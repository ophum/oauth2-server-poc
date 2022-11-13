package utils

import (
	"crypto"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func ParseEdPublicKeyFromPEMFile(path string) (crypto.PublicKey, error) {
	pem, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jwt.ParseEdPublicKeyFromPEM(pem)
}

func ParseEdPrivateKeyFromPEMFile(path string) (crypto.PrivateKey, error) {
	pem, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jwt.ParseEdPrivateKeyFromPEM(pem)
}
