package auth

import (
	"net/http"

	"github.com/urfave/negroni"
)

func BasicAuth(username, password string) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		user, pass, _ := r.BasicAuth()

		if user != username || pass != password {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Incorrect username/password"))
			return
		}

		next(w, r)
	})
}
