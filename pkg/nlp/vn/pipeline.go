package vn

import (
	ctx "context"
	"database/sql"
	"encoding/json"
	"enews/pkg/configs"
	"enews/pkg/db"
	"log"
	"strings"
	"sync"
)

// List of constants
const (
	MissingContent string = "<missing>"
)

// - - - - - STAGES - - - - - //

func createArticleChannel(articles ...*db.Article) <-chan *db.Article {
	out := make(chan *db.Article)
	go func() {
		defer close(out)
		for _, a := range articles {
			out <- a
		}
	}()
	return out
}

func annotate(articles <-chan *db.Article) (<-chan *Annotation, <-chan *Annotation, <-chan *db.Article, <-chan error) {
	out1 := make(chan *Annotation)      // For Annotation
	out2 := make(chan *Annotation)      // For Annotation
	nilAnnots := make(chan *db.Article) // For nil annotation
	errc := make(chan error, 1)         // For errors

	// Create a client for NLP server
	conf := configs.LoadConfigs().VnNLP
	nlp := NewNLPClient(conf.Host, conf.Port, conf.Annotators)
	go func() {
		defer close(out1)
		defer close(out2)
		defer close(nilAnnots)
		defer close(errc)
		for a := range articles {
			content := a.Content.String
			if len(content) > 0 && content != MissingContent {
				parsed, err := nlp.DepParse(content)
				if err != nil {
					if e, ok := err.(*Error); ok && e.Type == ErrorTypeTextTooLong {
						nilAnnots <- a
					} else {
						errc <- err
					}
				} else {
					annot := Annotation{a.ID, parsed}
					out1 <- &annot
					out2 <- &annot
				}
			}
		}
	}()
	return out1, out2, nilAnnots, errc
}

func saveAnnotatedArticle(edb *db.Queries, annotation <-chan *Annotation) <-chan error {
	errc := make(chan error, 1)
	go func() {
		defer close(errc)
		for a := range annotation {
			annotJSON, err := json.Marshal(&a.ParsedContent)
			if err != nil {
				errc <- err
			} else {
				err = edb.CreateAnnotatedArticle(
					ctx.Background(),
					db.CreateAnnotatedArticleParams{
						ArticleID:  a.ArticleID,
						Annotation: string(annotJSON),
					},
				)
				if err != nil {
					if !strings.Contains(err.Error(), "violates unique constraint") {
						errc <- &Error{
							Type: ErrorTypeUniqueConstraintViolation,
							Err:  err,
							Msg:  strings.Join([]string{"Duplicate Article ID", string(a.ArticleID)}, " "),
						}
					} else {
						errc <- err
					}
				}
			}
		}
	}()
	return errc
}

func filterNamedEntities(annotations <-chan *Annotation) <-chan *db.CreateStageExtractedEntityParams {
	out := make(chan *db.CreateStageExtractedEntityParams)
	go func() {
		defer close(out)
		for a := range annotations {
			if a == nil {
				out <- nil
			}
			entities := getArticleEntities(a)

			// Create StageExtractedEntity for output channel
			id := a.ArticleID
			for ent, count := range entities {
				ner := db.CreateStageExtractedEntityParams{
					ArticleID:  id,
					Entity:     sql.NullString{ent.Name, true},
					EntityType: sql.NullString{ent.Type, true},
					Counts:     int16(count),
				}
				out <- &ner
			}
		}
	}()
	return out
}

func saveStageExtractedEntities(edb *db.Queries, entities <-chan *db.CreateStageExtractedEntityParams) <-chan error {
	errc := make(chan error, 1)
	go func() {
		defer close(errc)
		for e := range entities {
			err := edb.CreateStageExtractedEntity(ctx.Background(), *e)
			if err != nil {
				errc <- err
			}
		}
	}()
	return errc
}

func mergeErrors(echans ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error, len(echans))
	pickError := func(echan <-chan error) {
		for err := range echan {
			out <- err
		}
		wg.Done()
	}
	for _, c := range echans {
		wg.Add(1)
		go pickError(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// - - - - - PREPARE A BATCH - - - - - //

// runPipeline processes a channel of articles to return any that has nil annotations and any error
func runPipeline(edb *db.Queries, articles []*db.Article) ([]*db.Article, error) {
	var errchans []<-chan error

	// 1. Create a channel of articles
	articleChan := createArticleChannel(articles...)

	// 2. Annotate each article
	annots1, annots2, nilAnnots, serverErrs := annotate(articleChan)
	errchans = append(errchans, serverErrs)

	// 3. Save articles' annotation
	errc1 := saveAnnotatedArticle(edb, annots1)
	errchans = append(errchans, errc1)

	// 4. Filter Entities from annotation of each article
	entities := filterNamedEntities(annots2)

	// 5. Save articles' entities into a stage table
	errc2 := saveStageExtractedEntities(edb, entities)
	errchans = append(errchans, errc2)

	// Merge errors from every stage
	errs := mergeErrors(errchans...)

	numServerErrors := 0
	var unfinished []*db.Article

	select {
	case err := <-errs:
		log.Println("Errs: ", err)
		if err, ok := err.(*Error); ok {
			switch err.Type {
			case ErrorTypeRequestFailed, ErrorTypeServerNotResponding, ErrorTypeServerError:
				numServerErrors++
			}
		}
		// TODO: need to use context to cancel the pipeline
		if numServerErrors > 0 {
			return nil, &Error{Type: ErrorTypeResetServer}
		}
	case na := <-nilAnnots:
		log.Println("Nil: ", na.ID)
		unfinished = append(unfinished, na)
		log.Println("Unfinished: ", len(unfinished))
	}
	return unfinished, nil
}

func fetchArticles(edb *db.Queries, ids []int32) []*db.Article {
	articles, err := edb.GetArticle_ByListID(ctx.Background(), ids)
	if err != nil {
		return nil
	}
	out := []*db.Article{}
	for i := range articles {
		out = append(out, &articles[i])
	}
	return out
}

// processBatch processes a batch of article IDs
func processBatch(edb *db.Queries, ids []int32) error {
	// 1. Get the articles
	articles := fetchArticles(edb, ids)
	unfinished, err := runPipeline(edb, articles)
	if err != nil {
		return err
	}

	log.Println("Total unfinished: ", len(unfinished))
	return nil
}

// RunExtractPipeline is the main controller of the NLP pipeline
func RunExtractPipeline(dbc *db.DB) bool {
	// Get number of articles per batch
	numArticlePerBatch := configs.LoadConfigs().VnNLP.NumberArticlePerBatch

	// Get unprocessed article IDs
	edb := dbc.Connect()
	ids, err := edb.GetUnprocessedArticleID(ctx.Background())

	// TODO: resolve this
	if err != nil || len(ids) == 0 {
		return true
	}

	// Prepare batches of IDs
	batches := chunkIDs(numArticlePerBatch, ids)
	log.Println("Total batches: ", len(batches))

	// Run batch with pipeline
	for i, batch := range batches {
		log.Println("\n>>>> BATCH: ", i, batch)
		err := processBatch(edb, batch)

		if err != nil {
			log.Println(">> PIPELINE ERR: ", err)
			if err, ok := err.(*Error); ok {
				if err.Type == ErrorTypeResetServer {
					log.Println(">>> NEED TO RESET SERVER: ", err)
					return false
				}
			}
		}
	}
	log.Println(">>> FINISHED ALL BATCHES")
	return true
}
