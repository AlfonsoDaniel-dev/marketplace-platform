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
	jwt.StandardClaims
}

type SendTSVLoginEmail struct {
	UserName  string
	Text      string
	Link      string
	FinalText string
}
