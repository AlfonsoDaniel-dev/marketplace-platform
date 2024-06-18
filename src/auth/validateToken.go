package auth

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"shopperia/src/common/models"
)

func ValidateToken(tokenString string) (models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, verifyFunc)
	if err != nil {
		return models.Claims{}, err
	}

	if !token.Valid {
		return models.Claims{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return models.Claims{}, errors.New("No se pudieron obtener los tokens")
	}

	return *claims, nil

}

func verifyFunc(t *jwt.Token) (interface{}, error) {
	return verifyKey, nil
}
