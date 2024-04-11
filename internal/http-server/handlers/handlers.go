package handlers

import (
	"Avito_go/internal/getters"
	"Avito_go/internal/http-server/accessHTTP"
	"Avito_go/internal/setters"
	"Avito_go/internal/storage"
	"Avito_go/internal/storage/accessDB"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

		_, err := accessDB.ValidateID(tagID)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = accessDB.ValidateID(featureID)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if useLastRevision == "true" {
			banner, err := getters.GetBannerByTagAndFeature(tagID, featureID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					w.WriteHeader(http.StatusNotFound)
				} else {
					accessHTTP.SendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
				}
				return
			}
			jsonBytes := append([]byte(banner), '\n')

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonBytes)
			if err != nil {
				return
			}
		} else if useLastRevision == "false" {
			banner, found := getters.GetCache(tagID, featureID)
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
				accessHTTP.SendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
				return
			}
			formattedJSON, err := json.MarshalIndent(data, "", "  ")
			jsonBytes := append(formattedJSON, '\n')
			if err != nil {
				accessHTTP.SendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonBytes)
			if err != nil {
				return
			}
		}
	}
}

func GetAllBannersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tagID := r.URL.Query().Get("tag_id")
		featureID := r.URL.Query().Get("feature_id")
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		if tagID == "" && featureID != "" {
			_, err := accessDB.ValidateID(featureID)
			if err != nil {
				accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else if tagID != "" && featureID == "" {
			_, err := accessDB.ValidateID(tagID)
			if err != nil {
				accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else {
			_, err := accessDB.ValidateID(tagID)
			if err != nil {
				accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			_, err = accessDB.ValidateID(featureID)
			if err != nil {
				accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		_, err := accessDB.ValidateLimitOffset(limit)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = accessDB.ValidateLimitOffset(offset)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		banner, err := getters.GetAllBanners(tagID, featureID, limit, offset)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				accessHTTP.SendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
			}
			return
		}
		jsonBytes, _ := json.MarshalIndent(banner, "", " ")
		jsonBytes = append(jsonBytes, '\n')
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonBytes)
		if err != nil {
			return
		}
	case "POST":
		db, err := sql.Open("postgres", storage.PsqlInfo)
		if err != nil {
			return
		}
		defer db.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Вставка нового баннера в базу данных
		newBannerID, err := setters.InsertBanner(data, db)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Формирование ответа с ID нового баннера
		response := map[string]int{"banner_id": newBannerID}
		responseBody, err := json.Marshal(response)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
			return
		}

		// Отправка ответа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseBody)
	}
}
