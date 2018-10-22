package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"log"
	"sync"
)

var db *sql.DB
var once sync.Once

func GetInstance() *sql.DB {
	once.Do(func() {
		connStr := "postgres://docker:docker@localhost/docker?sslmode=disable"
		var err error
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		if err = db.Ping(); err != nil {
			log.Fatal(err)
		}
	})
	return db
}
