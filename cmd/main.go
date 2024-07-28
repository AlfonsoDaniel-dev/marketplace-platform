package main

import (
	"log"
	"os"
	"shopperia/cmd/config"
	"shopperia/src/External/db/migrations"
	"shopperia/src/auth"
)

func main() {
	config.LoadEnvVars()

	if err := config.ValidateEnvVars(); err != nil {
		log.Fatalf("Error validating environment variables: %v", err)
	}

	database := config.ConnectPostgresDB()

	if err := auth.LoadFiles("./cmd/certificates/app.rsa.pub", "./cmd/certificates/app.rsa"); err != nil {
		log.Fatalf("Error loading certificates: %v", err)
	}

	migrator := migrations.NewMigrator(database)

	err := migrator.Migrate()
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	server := config.NewHttp(database)

	PORT := os.Getenv("APP_PORT")
	HOST := os.Getenv("APP_HOST")

	address := HOST + ":" + PORT

	if err := server.Start(address); err != nil {
		log.Fatalf("Error starting server. ERR: %v", err)
	}
	/*
		severCert := "./cmd/ssl/server.crt"
		severKey := "./cmd/ssl/server.key"

		if err := server.StartTLS(":443", severCert, severKey); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	*/
}
