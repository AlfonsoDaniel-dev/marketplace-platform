package models

import "github.com/golang-jwt/jwt"

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	IsAdmin  bool   `json:"is_admin"`
	Version  string `json:"version"`
	jwt.StandardClaims
}

type SendTSVLoginEmail struct {
	UserName  string
	Text      string
	Link      string
	FinalText string
}

type LoginStatus struct {
	Token                   string
	Error                   error
	HasTwoStepsVerification bool
	TsvConfirmationLink     string
}

type LoginConfirmationChan struct {
	Token     string
	Confirmed bool
	Error     error
}
