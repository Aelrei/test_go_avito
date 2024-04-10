package access

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

func ValidateID(id string) (int, error) {
	parsedID, err := strconv.Atoi(id)
	if err != nil || parsedID <= 0 {
		return 0, errors.New("not correct one of parameters")
	}
	return parsedID, nil
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
