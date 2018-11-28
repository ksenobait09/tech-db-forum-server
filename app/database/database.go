package database

import (
	"github.com/jackc/pgx"
	"log"
	"sync"
)

var db *pgx.ConnPool

var once sync.Once

func GetInstance() *pgx.ConnPool {
	once.Do(func() {
		var err error
		db, err = pgx.NewConnPool(pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "docker",
				Password: "docker",
				Database: "docker",
			},
			MaxConnections: 50,
		})
		if err != nil {
			log.Fatal(err)
		}
	})
	return db
}
