package models

import "github.com/golang-jwt/jwt/v4"

type GrantType string

const (
	GrantTypeClientCredentials GrantType = "client_credentials"
)

type ClientCredentialsClaims struct {
	*jwt.RegisteredClaims
}

type AccessTokenClaims struct {
	*jwt.RegisteredClaims
	Scope string `json:"scope"`
}
