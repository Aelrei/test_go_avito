package access

import (
	"net/http"
	"strings"
)

func AuthMiddlewareUserAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("token"), " ")
		if len(token) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if token[0] == "user_token" {
			next.ServeHTTP(w, r)
		} else if token[0] == "admin_token" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

	})
}

func AuthMiddlewareAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("token"), " ")
		if len(token) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if token[0] != "admin_token" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
