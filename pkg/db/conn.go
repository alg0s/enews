package db

import (
	"database/sql"
	"enews/pkg/configs"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// TODO: - add db engine verification

// Conn returns a sql.DB connection that connects to the database in user settings `configs.json`
func conn() *sql.DB {
	var conf = configs.LoadConfigs().Database
	var dsn = fmt.Sprintf(`host=%s port=%s user=%s dbname=%s sslmode=disable`,
		conf.Host,
		conf.Port,
		conf.Username,
		conf.Dbname,
	)
	conn, err := sql.Open(conf.Engine, dsn)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// Connect makes a connection to the db and provides a db object
func Connect() *Queries {
	var c = conn()
	return New(c)
}
