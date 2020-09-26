package vn

import (
	"enews/pkg/db"
)

// ExtractorPipeline defines the required functions for an NLP extractor
type ExtractorPipeline interface {
	createArticleChannel(as ...*db.Article) <-chan *db.Article
	annotate(as <-chan *db.Article) (<-chan *Annotation, <-chan *Annotation, <-chan error)
	saveAnnotatedArticle(edb *db.Queries, annots <-chan *Annotation) <-chan error
	filterNamedEntities(annots <-chan *Annotation) <-chan *db.CreateStageExtractedEntityParams
	saveStageExtractedEntities(edb *db.Queries, es <-chan *db.CreateStageExtractedEntityParams) <-chan error
	mergeErrors(errcs ...<-chan error) <-chan error
}
