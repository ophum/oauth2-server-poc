package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ophum/oauth2-server-poc/models"
	"github.com/ophum/oauth2-server-poc/utils"
)

func main() {
	clientPrivateKey, err := utils.ParseEdPrivateKeyFromPEMFile("./client-private.pem")
	if err != nil {
		log.Fatal(err)
	}

	clientCredentialsClaims := models.ClientCredentialsClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "client",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)), // 検証用に5分で作る。本来はもっと短い方がよさそう。
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodEdDSA, clientCredentialsClaims)
	token, err := t.SignedString(clientPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
}
