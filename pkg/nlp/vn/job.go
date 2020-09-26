package vn

// ExtractJob represents an article in extractor
type ExtractJob struct {
	articleID         int32
	content           string
	annotation        *[]Sentence
	isTooLong         bool
	savedAnnotArticle bool
	savedToStage      bool
}
