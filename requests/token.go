package requests

import "github.com/ophum/oauth2-server-poc/models"

type TokenRequest struct {
	GrantType models.GrantType `form:"grant_type"`
	Scope     string           `form:"scope"`
}
