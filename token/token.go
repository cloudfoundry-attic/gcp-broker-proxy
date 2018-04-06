package token

import (
	"fmt"
	"log"
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

			msg := fmt.Sprintf("Error retrieving OAuth token: %s", err.Error())
			log.Println(msg)
			w.Write([]byte(msg))
			return
		}

		r.Header.Set("Authorization", "Bearer "+token.AccessToken)

		next(w, r)
	})
}
