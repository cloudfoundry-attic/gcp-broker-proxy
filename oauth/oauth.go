package oauth

import (
	"context"
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

var (
	scopes = "https://www.googleapis.com/auth/cloud-platform"
)

type GCPOAuth struct {
	jwt *jwt.Config
}

func NewGCPOAuth(serviceAccountJSON string) (*GCPOAuth, error) {
	rawJSON := []byte(serviceAccountJSON)

	jwt, err := google.JWTConfigFromJSON(rawJSON, scopes)
	if err != nil {
		return nil, err
	}

	oauth := GCPOAuth{jwt}

	return &oauth, nil
}

func (o *GCPOAuth) GetToken() (*oauth2.Token, error) {
	var token *oauth2.Token
	token, err := o.jwt.TokenSource(context.Background()).Token()

	if err != nil {
		return token, err
	}

	if token.AccessToken == "" {
		return nil, errors.New("Missing access_token in oauth response")
	}

	return token, err
}
