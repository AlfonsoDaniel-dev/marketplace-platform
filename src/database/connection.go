package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

func ConnectToDatabase(connStr string) {
	once.Do(func() {
		var err error
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}

		if err := db.Ping(); err != nil {
			panic(err)
		}
	})
}

func PoodDB() *sql.DB {
	return db
}
