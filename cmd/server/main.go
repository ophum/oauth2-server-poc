package main

import (
	"crypto"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ophum/oauth2-server-poc/models"
	"github.com/ophum/oauth2-server-poc/requests"
	"github.com/ophum/oauth2-server-poc/utils"
)

var (
	// private key of oauth2 server
	serverPrivateKey crypto.PrivateKey
	// public key of client
	clientPublicKey crypto.PublicKey
)

// available scope
const (
	ScopeRead  = "read"
	ScopeWrite = "write"
)

func main() {
	var err error

	if clientPublicKey, err = utils.ParseEdPublicKeyFromPEMFile("./client-public.pem"); err != nil {
		log.Fatal(err)
	}

	if serverPrivateKey, err = utils.ParseEdPrivateKeyFromPEMFile("./server-private.pem"); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.POST("/token", token)
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func token(ctx *gin.Context) {
	var req requests.TokenRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println(err)
		ctx.Abort()
		return
	}

	token := ""
	switch req.GrantType {
	case models.GrantTypeClientCredentials:
		if t, err := clientCredentialsHandler(ctx, &req); err != nil {
			log.Println(err)
			ctx.Abort()
			return
		} else {
			token = t
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

func clientCredentialsHandler(ctx *gin.Context, req *requests.TokenRequest) (string, error) {
	clientJWT, err := getBearerToken(ctx)
	if err != nil {
		return "", err
	}

	if err := validateToken(clientJWT); err != nil {
		return "", err
	}

	scope := createAccessTokenScope(req.Scope)
	accessTokenClaims := models.AccessTokenClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "ophum",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Scope: scope,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(serverPrivateKey)
	if err != nil {
		return "", err
	}
	return accessTokenString, nil
}

func getBearerToken(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return "", errors.New("invalid header")
	}
	return strings.TrimPrefix(header, "Bearer "), nil
}

func validateToken(token string) error {
	var claims models.ClientCredentialsClaims
	t, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errors.New("invalid alg")
		}
		return clientPublicKey, nil
	})
	if err != nil {
		return err
	}

	if !t.Valid {
		return errors.New("token is invalid")
	}

	return nil
}

func createAccessTokenScope(requestScope string) string {
	requestScopes := strings.Split(requestScope, " ")
	accessTokenScopes := []string{}
	for _, v := range requestScopes {
		switch v {
		case ScopeRead:
			accessTokenScopes = append(accessTokenScopes, ScopeRead)
		case ScopeWrite:
			accessTokenScopes = append(accessTokenScopes, ScopeWrite)
		default:
			log.Println("invalid scope", v)
		}
	}

	return strings.Join(accessTokenScopes, " ")
}
