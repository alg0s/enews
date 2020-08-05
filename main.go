package main

import (
	"context"
	"enews/pkg/db"
	"log"
)

func connWithSQL() {
	var conn = db.Conn()
	var edb = db.New(conn)

	articles, err := edb.ListArticles(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	for i, a := range articles {
		log.Println(i, a)
	}

	log.Println("Done")
}

func main() {
	connWithSQL()
}
