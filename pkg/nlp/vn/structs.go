package vn

// Token details of an extract entity
type Token struct {
	Index    int    `json:"index"`
	Form     string `json:"form"`
	PosTag   string `json:"posTag"`
	NerLabel string `json:"nerLabel"`
	Head     int    `json:"head"`
	DepLabel string `json:"depLabel"`
}

// ParsedSentence is an array of Token's
type ParsedSentence []Token

// ServerResponse is the response from NLPServer
type ServerResponse struct {
	Sentences []ParsedSentence `json:"sentences"`
	Status    bool             `json:"status"`
	Error     string           `json:"error"`
	Language  string           `json:"language"`
}

// ServerInfo holds information of a running NLPServer
type ServerInfo struct {
	Host       string
	Port       string
	Address    string
	Annotators string
}

// ParsedArticle represents an extracted article with identified tokens
type ParsedArticle struct {
	ID        int32
	Sentences *[]ParsedSentence
}

// Entity represents a named entity recognized by the NLP service
type Entity struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Annotation represents an article that is already annotated
type Annotation struct {
	ArticleID     int32
	ParsedContent *[]ParsedSentence
}
