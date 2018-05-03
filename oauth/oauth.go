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
	jwt   *jwt.Config
	token *oauth2.Token
}

func NewGCPOAuth(serviceAccountJSON string) (*GCPOAuth, error) {
	rawJSON := []byte(serviceAccountJSON)

	jwt, err := google.JWTConfigFromJSON(rawJSON, scopes)
	if err != nil {
		return nil, err
	}

	oauth := GCPOAuth{jwt, nil}

	return &oauth, nil
}

func (o *GCPOAuth) GetToken() (*oauth2.Token, error) {
	tokenSource := oauth2.ReuseTokenSource(o.token, o.jwt.TokenSource(context.Background()))

	var err error
	o.token, err = tokenSource.Token()

	if err != nil {
		return o.token, err
	}

	if o.token.AccessToken == "" {
		return nil, errors.New("Missing access_token in oauth response")
	}

	return o.token, err
}
