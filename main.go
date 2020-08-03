package main

import (
	pkg "enews/pkg/db"
	schm "enews/schema"
	"log"
)

func main() {

	var db = pkg.ConnectLocalDB()
	db.Conn.Ping()
	log.Println("Here")

	var articles = []schm.RawArticle{}

	db.Conn.Select(&articles, "SELECT * FROM enews.raw_articles LIMIT 3")

	log.Println("Total articles: ", len(articles))

	for i, a := range articles {
		log.Println("\n---", i)
		log.Println(a.AddedID, a.Category.String, a.Summary.String, a.PublishDate)
	}

	log.Println("End")
}
