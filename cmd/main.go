package main

import (
	"log"
	"os"
	"shopperia/cmd/config"
	"shopperia/src/auth"
	"shopperia/src/database/migrations"
)

func main() {
	config.LoadEnvVars()

	if err := config.ValidateEnvVars(); err != nil {
		log.Fatalf("Error validating environment variables: %v", err)
	}

	db := config.ConnectPostgresDB()

	if err := auth.LoadFiles("./cmd/certificates/app.rsa.pub", "./cmd/certificates/app.rsa"); err != nil {
		log.Fatalf("Error loading certificates: %v", err)
	}

	migrator := migrations.NewMigrator(db)

	if err := migrator.Migrate(); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	server := config.NewHttp(db)

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
