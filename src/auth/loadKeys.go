package auth

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"io/ioutil"
	"sync"
)

var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
	once      sync.Once
)

func LoadFiles(publicKeyPath string, privateKeyPath string) error {
	var err error
	once.Do(func() {
		err = loadFiles(publicKeyPath, privateKeyPath)
	})

	return err
}

func loadFiles(publicPath string, privatePath string) error {
	publicKey, err := ioutil.ReadFile(publicPath)
	if err != nil {
		return err
	}

	privateKey, err := ioutil.ReadFile(privatePath)
	if err != nil {
		return err
	}

	return ParseRSA(publicKey, privateKey)
}

func ParseRSA(publicBytes, privateBytes []byte) error {
	var err error

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return err
	}

	return nil
}
