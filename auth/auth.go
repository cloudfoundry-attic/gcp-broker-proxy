package auth

import (
	"net/http"
)

func BasicAuth(handler http.HandlerFunc, username, password string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()

		if user != username || pass != password {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Incorrect username/password"))
			return
		}

		handler(w, r)
	})
}
