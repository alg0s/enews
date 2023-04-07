package vn

import (
	ctx "context"
	"database/sql"
	"encoding/json"
	"enews/pkg/configs"
	"enews/pkg/db"
	"log"
	"strings"
)

// List of constants
const (
	MissingContent string = "<missing>"
)

// Extractor to extract NER from articles
type Extractor struct {
	db         *db.Queries
	nlp        *NLPClient
	batches    [][]int32
	unfinished []*db.Article
}

func (e *Extractor) fetchArticles(ids []int32) []db.Article {
	articles, err := e.db.GetArticle_ByListID(ctx.Background(), ids)
	if len(articles) == 0 {
		log.Println(">>> Articles Empty: ", articles)
	}
	if err != nil {
		return nil
	}
	return articles
}

func (e *Extractor) getBatches() error {
	// Get list of unprocessed article IDs
	ids, err := e.db.GetUnprocessedArticleID(ctx.Background())
	if err != nil {
		return err
	}
	// Group IDs into small batches
	numPerBatch := configs.LoadConfigs().VnNLP.NumberArticlePerBatch
	batches := chunkIDs(numPerBatch, ids)
	e.batches = batches
	return nil
}

func (e *Extractor) annotate(s string) (*[]ParsedSentence, error) {
	parsed, err := e.nlp.DepParse(s)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func (e *Extractor) saveAnnotatedArticle(id int32, annotation *[]ParsedSentence) error {
	annotJSON, err := json.Marshal(annotation)
	if err != nil {
		return err
	}
	err = e.db.CreateAnnotatedArticle(
		ctx.Background(),
		db.CreateAnnotatedArticleParams{ArticleID: id, Annotation: string(annotJSON)},
	)
	if err != nil {
		if !strings.Contains(err.Error(), "violates unique constraint") {
			return &Error{
				Type: ErrorTypeUniqueConstraintViolation,
				Err:  err,
				Msg:  strings.Join([]string{"Duplicate Article ID", string(id)}, " "),
			}
		}
		return err
	}
	return nil
}

func (e *Extractor) saveToStageTable(id int32, annotation *[]ParsedSentence) error {
	// 1. Filter NERs from parsed sentences
	ners := getArticleNERs(annotation)

	// 2. Save to staging table
	for entity, count := range ners {
		err := e.db.CreateStageExtractedEntity(
			ctx.Background(),
			db.CreateStageExtractedEntityParams{
				ArticleID:  id,
				Entity:     sql.NullString{entity.Name, true},
				EntityType: sql.NullString{entity.Type, true},
				Counts:     int16(count),
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Extractor) processArticle(a *db.Article) error {
	id := a.ID
	content := a.Content.String
	if len(content) == 0 || content == MissingContent {
		log.Println("Article is nill: ", id)
		return nil
	}

	// 1. Annotate
	annotation, err := e.annotate(content)
	if err != nil {
		return err
	}
	if annotation == nil {
		return &Error{Type: ErrorTypeNilAnnotation}
	}

	// 2. Save article's annotation
	err = e.saveAnnotatedArticle(id, annotation)
	if err != nil {
		return err
	}

	// 3. Save to staging table
	err = e.saveToStageTable(id, annotation)
	if err != nil {
		return err
	}

	return nil
}

func (e *Extractor) processBatch(ids []int32) {
	articles := e.fetchArticles(ids)

	for _, a := range articles {
		err := e.processArticle(&a)
		if err != nil {
			log.Println("ERR:  ", err)
			e.unfinished = append(e.unfinished, &a)
		}
	}
	err := e.processStagedEntities()
	if err != nil {
		log.Println("Unable to process staged entities: ", err)
	}
}

func (e *Extractor) processStagedEntities() error {
	// 1. Save new entities
	err := e.db.InsertNewEntitiesFromStagedEntities(ctx.Background())
	if err != nil {
		return err
	}

	// 2. Save new article entities
	err = e.db.InsertNewArticleEntitiesFromStagedEntities(ctx.Background())
	if err != nil {
		return err
	}

	// 3. Truncate staged entities table
	return nil
}

func (e *Extractor) run() error {
	// var wg sync.WaitGroup

	err := e.getBatches()
	if err != nil {
		log.Panic("Unable to get article batches: ", err)
		return err
	}

	log.Println(">> Total batches: ", len(e.batches))
	for i, b := range e.batches {
		log.Printf("** >>>: Processing batch #%v: %v", i, b)
		// wg.Add(1)
		e.processBatch(b)
	}
	// wg.Wait()
	log.Println(">>> All Batches Finished. Unfinished: ", len(e.unfinished))

	// Process StagedEntities
	return nil
}

// RunExtractor run an extractor to retrieve and extract entities from unprocessed articles
func RunExtractor(db *db.DB) error {
	conf := configs.LoadConfigs().VnNLP
	dbQuery := db.GetConn()
	nlp := NewNLPClient(conf.Host, conf.Port, conf.Annotators)
	e := Extractor{db: dbQuery, nlp: nlp}
	err := e.run()
	log.Println(">>> Execution Err: ", err)
	return err
}
