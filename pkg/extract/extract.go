package extract

/*
   Entity Extraction steps:
   - Initiatize queue
   - Initialize NLP servers
   - Truncate table `stage_extracted_article_entities`
   - Check table `articles`, terminates if empty
   - Loader starts loading articles by batches, which is a parameter of Extractor,
   and drop them into the queue
   - Extractor workers start picking up articles from the queue and process them,
   inserting outputs into the table `stage_extracted_article_entities`
   - Once the queue is empty, meaning workers have processed all articles in the batch,
   the Controller will run the SQL query to process the extracted articles:
        - insert article ID into table `extracted_articles`
        - insert new article and entities into table `article_entities`
        - insert new entity_type into table `entity_type`
          and new entities into table `unique_entities`
   - Repeat for a new batch
*/

/*
   Ref:
      https://www.opsdash.com/blog/job-queues-in-go.html
*/

import (
	ctx "context"
	"enews/pkg/db"
	nlp "enews/pkg/nlp/vn"
	"enews/pkg/queue"
	"log"
)

// Extract extracts entities from article
func Extract() {
	// Get articles
	edb := db.Connect()

	articles, err := edb.GetArticle_Limit(ctx.Background(), 3)

	if err != nil {
		log.Fatal(err)
	}

	// Initialize an ArticleQueue
	var q = queue.NewQueue()

	for _, a := range articles {
		q.Enqueue(a)
	}

	// Initialize NLP service
	var s = nlp.NewVnNLPServer()

	// Extract entities from articles
	for q.IsEmpty() == false {
		a := q.Dequeue()
		if a == nil {
			break
		}
		content := a.(db.Article).Content.String
		parsed := s.Ner(content)

		log.Println("NER: ", parsed)
	}
}
