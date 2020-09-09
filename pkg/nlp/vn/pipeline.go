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

func fetchArticles(quantity int) []db.Article {
	edb := db.Connect()
	articles, err := edb.GetArticles_Limit(ctx.Background(), int32(quantity))
	log.Println("Total articles in batch: ", len(articles))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return articles
}

func createArticleChannel(articles ...db.Article) <-chan *db.Article {
	out := make(chan *db.Article)
	go func() {
		defer close(out)
		for i, a := range articles {
			log.Println("# ", i)
			out <- &a
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
	st := configs.LoadConfigs().VnNLP
	nlp := NewNLPClient(st.Host, st.Port, st.Annotators)

	go func() {
		defer close(out1)
		defer close(out2)
		defer close(nilAnnots)
		defer close(errc)
		for a := range articles {
			// TODO: handle empty content, probably in the SQL
			content := a.Content.String
			parsed, err := nlp.DepParse(content)

			if err != nil {
				if e, ok := err.(*Error); ok && e.Type == ErrorTypeNilAnnotation {
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
	}()
	return out1, out2, nilAnnots, errc
}

func saveAnnotatedArticle(annotation <-chan *Annotation) <-chan error {
	errc := make(chan error, 1)
	go func() {
		defer close(errc)
		edb := db.Connect()
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
						errc <- &Error{Type: ErrorTypeUniqueConstraintViolation, Err: err}
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

func saveStageExtractedEntities(entities <-chan *db.CreateStageExtractedEntityParams) <-chan error {
	errc := make(chan error, 1)
	go func() {
		defer close(errc)
		edb := db.Connect()
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
		for e := range echan {
			out <- e
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

// TODO: implement this to handle large articles, not for now
func handleNilAnnotation(articles <-chan *db.Article) {
	go func() {
		for a := range articles {
			log.Println("\n Nil Annotation: ", a.ID, len(a.Content.String))
		}
	}()
}

// RunExtractPipeline is the main controller of the NLP pipeline
func RunExtractPipeline() bool {
	var errchans []<-chan error

	numArticles := configs.LoadConfigs().VnNLP.VnNLPConfigs.NumberArticlePerBatch

	batch := fetchArticles(numArticles)

	achan := createArticleChannel(batch...)

	annots1, annots2, nilAnnots, serverErrs := annotate(achan)
	errchans = append(errchans, serverErrs)

	errc1 := saveAnnotatedArticle(annots1)
	errchans = append(errchans, errc1)

	entities := filterNamedEntities(annots2)

	errc2 := saveStageExtractedEntities(entities)
	errchans = append(errchans, errc2)

	handleNilAnnotation(nilAnnots)

	errs := mergeErrors(errchans...)

	numServerErrors := 0

	for err := range errs {
		log.Println(">>> ERR: ", err)
		if err, ok := err.(*Error); ok {
			switch err.Type {
			case ErrorTypeRequestFailed:
				numServerErrors++
			case ErrorTypeServerNotResponding:
				numServerErrors++
			case ErrorTypeServerError:
				numServerErrors++
			}
		}
	}

	if numServerErrors > 0 {
		return false
	}
	return true
}

/** Thoughts **

It is actually not bad to fail a pipeline once a stage fails. There are two main types of behaviors:
First is acceptable failure which will be recorded but will not fail the pipeline.
Second is unacceptable failure, for instance server not responding. It will fail the pipeline
and signal the Main controller to do something with the server.
*/
