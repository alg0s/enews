package db

import (
	"database/sql"
	"enews/pkg/configs"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// DB provides a connection to the database
type DB struct {
	db *sql.DB
}

// Init starts a connection with the database defined in `configs.json`
func (d *DB) init() error {
	var conf = configs.LoadConfigs().Database
	var dsn = fmt.Sprintf(`host=%s port=%s user=%s dbname=%s sslmode=disable`,
		conf.Host,
		conf.Port,
		conf.Username,
		conf.Dbname,
	)
	conn, err := sql.Open(conf.Engine, dsn)

	// Setup Connection Pool
	conn.SetMaxOpenConns(conf.MaxOpenConnection)
	conn.SetMaxIdleConns(conf.MaxIdleConnection)
	conn.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Minute)

	if err != nil {
		return err
	}

	if err = conn.Ping(); err != nil {
		return err
	}
	d.db = conn
	return nil
}

// GetConn returns a Queries pointer that can be used to perform db transactions
func (d *DB) GetConn() *Queries {
	return New(d.db)
}

// ConnectDB creates a new instance of DB
func ConnectDB() (*DB, error) {
	var db DB
	err := db.init()
	if err != nil {
		return nil, err
	}
	return &db, nil
}
