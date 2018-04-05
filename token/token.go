package token

import (
	"net/http"

	"github.com/urfave/negroni"

	"golang.org/x/oauth2"
)

//go:generate counterfeiter . TokenRetriever
type TokenRetriever interface {
	GetToken() (*oauth2.Token, error)
}

func TokenHandler(tr TokenRetriever) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		token, err := tr.GetToken()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Error retrieving OAuth token"))
			return
		}

		r.Header.Set("Authorization", "Bearer "+token.AccessToken)

		next(w, r)
	})
}
