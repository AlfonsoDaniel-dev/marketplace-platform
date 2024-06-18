package helpers

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func EncryptPassword(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return "", err
	}

	return string(encryptedPassword), nil
}

func ComparePasswords(OldPassword []byte, password string) bool {
	ok := bcrypt.CompareHashAndPassword(OldPassword, []byte(password))
	if ok != nil {
		log.Println("error while comparing password. ERR: ", ok)
		return false
	}

	return true
}
