package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"os"
	"shopperia/src/common/models"
	"time"
)

func GenerateToken(loginModel models.Login, username string, isAdmin bool) (string, error) {
	TokenDuration := time.Now().Add(time.Hour * 48).Unix()
	Issuer := os.Getenv("JWT_ISSUER")
	tokenCreationTime := time.Now().Unix()
	tokenId := uuid.New().String()

	claim := models.Claims{
		Email:    loginModel.Email,
		UserName: username,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			Id:        tokenId,
			ExpiresAt: TokenDuration,
			IssuedAt:  tokenCreationTime,
			Issuer:    Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
