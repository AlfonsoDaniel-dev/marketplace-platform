package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"sync"
)

var (
	database *sql.DB
	once     sync.Once
)

func ConnectToDatabase(connStr string) {
	once.Do(func() {
		var err error
		database, err = sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}

		if err := database.Ping(); err != nil {
			panic(err)
		}
	})
}

func PooldDB() *sql.DB {
	return database
}
