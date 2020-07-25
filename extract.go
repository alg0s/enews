package main

/*
   REF:
   http://jmoiron.github.io/sqlx/#mapping
*/

import (
	"log"

	_ "github.com/lib/pq"
)

const (
	dbHost string = `localhost`
	dbPort string = `5432`
	dbName string = `postgres`
	dbUser string = `steve`
	dbPw   string = ``
)

// Entity represents the table in enet
type Entity struct {
	ID   int
	Name string
	Type string
}

const (
	createEntities = `
		CREATE TABLE IF NOT EXISTS enews.entity (
			"id" 	SERIAL,
			"name" 	VARCHAR(500) NOT NULL,
			"type"	VARCHAR(20) NOT NULL 
		)
	`
)

func test() {
	var db = Db{Host: dbHost, Port: dbPort, DbName: dbName, Username: dbUser}
	var isConnected = db.Connect()
	log.Println("Is Connected: ", isConnected)

	// test select
	rows, err := db.Conn.Queryx(`Select * FROM enews.entity`)

	if err != nil {
		log.Panic(err)
	}

	var entities []Entity

	for rows.Next() {
		var e Entity
		err = rows.StructScan(&e)
		entities = append(entities, e)
		log.Println(e)
	}

	for i, e := range entities {
		log.Println("\n%s - %s", i, e)
	}
}
