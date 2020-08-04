package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	dbHostLocal string = `localhost`
	dbPortLocal string = `5432`
	dbNameLocal string = `postgres`
	dbUserLocal string = `steve`
	dbPwLocal   string = ``
)

// Dbx wraps an instance of DB connection using sqlx
type Dbx struct {
	Conn     *sqlx.DB
	Host     string
	Port     string
	DbName   string
	Username string
	Password string
}

// Connect estsablishes a db connection and return true if successful, otherwise false
func (db *Dbx) Connect() bool {
	var dsn = fmt.Sprintf(`host=%s port=%s user=%s dbname=%s sslmode=disable`,
		db.Host, db.Port, db.Username, db.DbName)
	db.Conn = sqlx.MustConnect("postgres", dsn)
	return true
}

// ConnectLocalDBx Create a connection with the default local Postgresql database
func ConnectLocalDBx() *Dbx {
	var db = Dbx{Host: dbHostLocal, Port: dbPortLocal, DbName: dbNameLocal, Username: dbUserLocal}
	db.Connect()
	return &db
}

// Db wraps an instance of DB connection using database/sql
type Db struct {
	Conn     *sql.DB
	Host     string
	Port     string
	DbName   string
	Username string
	Password string
}

// Connect estsablishes a db connection and return true if successful, otherwise false
func (db *Db) Connect() bool {
	var dsn = fmt.Sprintf(`host=%s port=%s user=%s dbname=%s sslmode=disable`,
		db.Host, db.Port, db.Username, db.DbName)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	// defer conn.Close()
	db.Conn = conn
	return true
}

// ConnectLocalDB Create a connection with the default local Postgresql database
func ConnectLocalDB() *Db {
	var db = Db{Host: dbHostLocal, Port: dbPortLocal, DbName: dbNameLocal, Username: dbUserLocal}
	db.Connect()
	return &db
}
