package vn

import (
	ctx "context"
	"database/sql"
	"encoding/json"
	"enews/pkg/configs"
	"enews/pkg/db"
	"log"
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

func annotate(articles <-chan *db.Article) (<-chan *Annotation, <-chan *Annotation, <-chan *db.Article) {
	out1 := make(chan *Annotation)
	out2 := make(chan *Annotation)
	failed := make(chan *db.Article)

	// Create a client for NLP server
	st := configs.LoadConfigs().VnNLP
	nlp := NewNLPClient(st.Host, st.Port, st.Annotators)

	go func() {
		defer close(out1)
		defer close(out2)
		defer close(failed)
		for a := range articles {
			// TODO: handle empty content, probably in the SQL
			content := a.Content.String
			// content := a.Title.String
			parsed, err := nlp.DepParse(content)
			if err != nil || parsed == nil {
				failed <- a
			} else {
				annot := Annotation{a.ID, parsed}
				out1 <- &annot
				out2 <- &annot
			}
		}
	}()
	return out1, out2, failed
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
					errc <- err
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
				log.Println("Err: ", err)
				errc <- err
			}
		}
	}()
	return errc
}

func handleFailed(articles <-chan *db.Article) {
	go func() {
		for a := range articles {
			log.Println("\n FAILED article: ", a.Title.String)
		}
	}()
}

func mergeErrors(echans ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error, len(echans))
	pickError := func(echan <-chan error) {
		for e := range echan {
			log.Println("1 err: ", e)
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

// RunExtractPipeline is the main controller of the NLP pipeline
func RunExtractPipeline() {
	numArticles := configs.LoadConfigs().VnNLP.VnNLPConfigs.NumberArticlePerBatch

	batch := fetchArticles(numArticles)

	achan := createArticleChannel(batch...)

	annots1, annots2, failed := annotate(achan)

	var errchans []<-chan error

	errc1 := saveAnnotatedArticle(annots1)
	errchans = append(errchans, errc1)

	entities := filterNamedEntities(annots2)

	errc2 := saveStageExtractedEntities(entities)
	errchans = append(errchans, errc2)

	handleFailed(failed)

	errs := mergeErrors(errchans...)
	for e := range errs {
		log.Println(e)
	}
}
