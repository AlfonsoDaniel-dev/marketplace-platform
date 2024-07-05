package config

import (
	"database/sql"
	"fmt"
	"os"
	"shopperia/src/db"
)

func makeConnStr(usr, pwd, host, port, DBName string) string {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", usr, pwd, host, port, DBName)

	return connStr
}

func ConnectPostgresDB() *sql.DB {
	var (
		usr    = os.Getenv("DB_USER")
		pwd    = os.Getenv("DB_PASSWORD")
		host   = os.Getenv("DB_HOST")
		port   = os.Getenv("DB_PORT")
		DBName = os.Getenv("DB_NAME")
	)

	connStr := makeConnStr(usr, pwd, host, port, DBName)

	db.ConnectToDatabase(connStr)

	database := db.PooldDB()

	return database
}
