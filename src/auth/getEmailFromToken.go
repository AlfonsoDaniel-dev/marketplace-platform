package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func GetEmailFromToken(TokenString string) (string, error) {
	parts := strings.Split(TokenString, ":")
	if len(parts) != 3 {
		return "", errors.New("invalid token")
	}

	tokenPayLoad, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.New("Error while getting payload from token")
	}

	var Claims map[string]interface{}
	if err := json.Unmarshal(tokenPayLoad, &Claims); err != nil {
		return "", errors.New("Error while unmarshalling payload")
	}

	email, ok := Claims["email"].(string)
	if !ok {
		return "", errors.New("Error while getting email from token")
	}

	return email, nil
}

func GetFromToken(tokenString, whatThing string) (interface{}, error) {
	parts := strings.Split(tokenString, ":")
	if len(parts) != 3 {
		return "", errors.New("invalid token")
	}

	tokenPayLoad, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.New("Error while getting payload from token")
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(tokenPayLoad, &claims); err != nil {
		return "", errors.New("error while unmarshaling token")
	}

	var value interface{}
	var ok bool

	switch whatThing {
	case "email":
		value, ok = claims["email"].(string)
		if !ok {
			errStr := fmt.Sprintf("no %v field in token", whatThing)
			return "", errors.New(errStr)
		}
	case "user_name":
		value, ok = claims["user_name"].(string)
		if !ok {
			errStr := fmt.Sprintf("no %v field in token", whatThing)
			return "", errors.New(errStr)
		}
	case "is_admin":
		value, ok = claims["is_admin"].(bool)
		if !ok {
			errStr := fmt.Sprintf("no %v field in token", whatThing)
			return "", errors.New(errStr)
		}
	}

	return value, nil
}
