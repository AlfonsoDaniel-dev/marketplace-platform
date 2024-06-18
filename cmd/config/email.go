package config

import (
	"errors"
	"os"
	"strings"
)

type emailConfig struct {
	Email      string
	Name       string
	Password   string
	Host       string
	ServerName string
}

func getVars() (string, string, string, string, string, error) {
	Email := strings.TrimSpace(os.Getenv("EMAIL_PROVIDER_ACCOUNT"))
	Name := os.Getenv("EMAIL_PROVIDER_NAME")
	Password := strings.TrimSpace(os.Getenv("EMAIL_PROVIDER_PASSWORD"))
	Host := os.Getenv("EMAIL_PROVIDER_HOST")
	serverName := os.Getenv("EMAIL_PROVIDER_SERVER_NAME")

	enVar := []string{Email, Name, Password, Host, serverName}
	for _, v := range enVar {
		if v == "" {
			return Email, Name, Password, Host, serverName, errors.New("Error al cargar var")
		}
	}

	return Email, Name, Password, Host, serverName, nil
}

func NewEmailConfig() *emailConfig {
	// email, name, password, host, serverName, err := getVars()
	/*if err != nil {
		log.Fatal(err)
	} */
	return &emailConfig{
		Email:      "alfonso.dancerod@gmail.com",
		Name:       "Alfonso Daniel Cervantes Rodriguez",
		Password:   "mqgj nzzd mqxe ndce",
		Host:       "smtp.gmail.com",
		ServerName: "smtp.gmail.com:465",
	}
}

func (e *emailConfig) GetFields() (string, string, string, string, string) {
	return e.Email, e.Name, e.Password, e.Host, e.ServerName
}
