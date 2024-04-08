package handlers

import (
	"Avito_go/internal/http-server/gocache"
	"Avito_go/internal/storage"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func GetUserBanner(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tagID := r.URL.Query().Get("tag_id")
		featureID := r.URL.Query().Get("feature_id")
		useLastVersion := r.URL.Query().Get("use_last_version")

		if useLastVersion == "" {
			useLastVersion = "false"
		}

		if useLastVersion == "true" {
			banner, err := GetBannerByTagAndFeature(tagID, featureID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
			jsonBytes, _ := json.MarshalIndent(banner, "", " ")
			jsonBytes = append(jsonBytes, '\n')
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonBytes)
		} else {
			banner, err := gocache.GetCache(tagID, featureID)
			banner := fmt.Sprintf(string(banner.([]byte)))
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
			jsonBytes, _ := json.MarshalIndent(banner, "", " ")
			jsonBytes = append(jsonBytes, '\n')
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonBytes)
		}
	}
}

func GetBannerByTagAndFeature(tagID, featureID string) (*storage.Banner, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		storage.Host, storage.Port, storage.User, storage.Password, storage.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	query := `
		SELECT b.id, b.content, b.feature_id, t.id, b.active, b.created_at, b.updated_at 
		FROM features as f, banners as b, banner_tag as bt, tags as t
		WHERE t.id = bt.tag_id 
		  AND bt.banner_id = b.id 
		  AND f.id = b.feature_id  
		  AND t.id = $1
		  AND b.feature_id = $2;
	`

	rows := db.QueryRow(query, tagID, featureID)

	var banner storage.Banner
	err = rows.Scan(&banner.Id, &banner.Content, &banner.Feature_id, &banner.Tag_id, &banner.Active, &banner.Created_at, &banner.Updated_at)
	if err != nil {
		return nil, fmt.Errorf("error during scanning result set: %w", err)
	}

	return &banner, nil
}
