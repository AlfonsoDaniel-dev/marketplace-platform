package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func LoadEnvVars() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var envVars = []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "APP_PORT", "ALLOWED_ORIGINS"}

func ValidateEnvVars() error {
	for _, envVar := range envVars {
		if strings.TrimSpace(os.Getenv(envVar)) == "" {
			log.Fatalf("Environment variable %s is not set", envVar)
		}
	}

	return nil
}
