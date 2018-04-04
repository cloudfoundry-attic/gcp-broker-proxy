package token

import (
	"net/http"

	"golang.org/x/oauth2"
)

//go:generate counterfeiter . TokenRetriever
type TokenRetriever interface {
	GetToken() (*oauth2.Token, error)
}

func TokenHandler(handler http.HandlerFunc, tr TokenRetriever) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := tr.GetToken()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Error retrieving OAuth token"))
			return
		}

		r.Header.Set("Authorization", "Bearer "+token.AccessToken)

		handler(w, r)
	})
}
