package accessHTTP

import (
	"encoding/json"
	"net/http"
	"strings"
)

func AuthMiddlewareUserAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("token"), " ")
		if len(token) != 1 || token[0] == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if token[0] == "user_token" {
			next.ServeHTTP(w, r)
		} else if token[0] == "admin_token" {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusForbidden)
			return
		}

	})
}

func AuthMiddlewareAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("token"), " ")
		if len(token) != 1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if token[0] == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if token[0] != "admin_token" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SendErrorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	errorMessage := map[string]string{"error": message}
	jsonBytes, err := json.Marshal(errorMessage)
	jsonBytes, err = json.MarshalIndent(errorMessage, "", " ")
	jsonBytes = append(jsonBytes, '\n')
	if err != nil {
		http.Error(w, "Error during request ", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}
