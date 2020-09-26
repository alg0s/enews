package vn

import (
	ctx "context"
	"enews/pkg/configs"
	"enews/pkg/db"
)

// Extractor to extract NER from articles
type Extractor struct {
	db *db.Queries
	batches
}

func (e *Extractor) fetchArticles(ids []int32) []db.Article {
	articles, err := e.db.GetArticle_ByListID(ctx.Background(), ids)
	if err != nil {
		return nil
	}
	return articles
}

func (e *Extractor) run() error {
	numPerBatch := configs.LoadConfigs().VnNLP.NumberArticlePerBatch

	// Get list of unprocessed article IDs
	ids, err := e.db.GetUnprocessedArticleID(ctx.Background())
	if err != nil {
		return err
	}

	//
}
