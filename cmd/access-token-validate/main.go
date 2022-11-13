package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ophum/oauth2-server-poc/models"
	"github.com/ophum/oauth2-server-poc/utils"
)

var (
	accessToken = ""
)

func main() {
	flag.StringVar(&accessToken, "token", "", "access token")
	flag.Parse()

	serverPublicKey, err := utils.ParseEdPublicKeyFromPEMFile("./server-public.pem")
	if err != nil {
		log.Fatal(err)
	}
	var accessTokenClaims models.AccessTokenClaims
	t, err := jwt.ParseWithClaims(accessToken, &accessTokenClaims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errors.New("invalid alg")
		}
		return serverPublicKey, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if !t.Valid {
		log.Fatal("invalid token")
	}

	j, err := json.MarshalIndent(accessTokenClaims, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(j))
}
