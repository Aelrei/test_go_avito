package handlers

import (
	"Avito_go/internal/http-server/gocache"
	"Avito_go/internal/storage"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetUserBannerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tagID := r.URL.Query().Get("tag_id")
		featureID := r.URL.Query().Get("feature_id")
		useLastRevision := r.URL.Query().Get("use_last_revision")

		if useLastRevision == "" {
			useLastRevision = "false"
		}

		_, err := validateID(tagID)
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = validateID(featureID)
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if useLastRevision == "true" {
			banner, err := GetBannerByTagAndFeature(tagID, featureID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					w.WriteHeader(http.StatusNotFound)
				} else {
					sendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
				}
				return
			}
			jsonBytes := append([]byte(banner), '\n')

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(jsonBytes)
			if err != nil {
				return
			}
		} else if useLastRevision == "false" {
			banner, found := gocache.GetCache(tagID, featureID)
			if found != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			var data map[string]interface{}
			var ans storage.Banner
			err = json.Unmarshal(banner.([]byte), &ans)
			str := []byte(ans.Content)
			err = json.Unmarshal([]byte(str), &data)
			if err != nil {
				sendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
				return
			}
			formattedJSON, err := json.MarshalIndent(data, "", "  ")
			jsonBytes := append(formattedJSON, '\n')
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(jsonBytes)
			if err != nil {
				return
			}
		}
	}
}

//func GetAllBannersHandler()

func GetBannerByTagAndFeature(tagID, featureID string) (string, error) {
	db, err := sql.Open("postgres", storage.PsqlInfo)
	if err != nil {
		return "", fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	query := `
		SELECT b.content
		FROM features as f, banners as b, banner_tag as bt, tags as t
		WHERE t.id = bt.tag_id 
		  AND bt.banner_id = b.id 
		  AND f.id = b.feature_id  
		  AND t.id = $1
		  AND b.feature_id = $2;
	`

	rows := db.QueryRow(query, tagID, featureID)

	var content string
	err = rows.Scan(&content)
	if err != nil {
		return "", fmt.Errorf("error during scanning result set: %w", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling content: %w", err)
	}

	formattedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error formatting JSON: %w", err)
	}

	return string(formattedJSON), nil
}

func validateID(id string) (int, error) {
	parsedID, err := strconv.Atoi(id)
	if err != nil || parsedID <= 0 {
		return 0, errors.New("not correct one of parameters")
	}
	if parsedID <= 0 {
		return 0, errors.New("not correct one of parameters")
	}
	return parsedID, nil
}

func sendErrorResponse(w http.ResponseWriter, status int, message string) {
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
