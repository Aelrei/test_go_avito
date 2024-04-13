package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gitlab.com/Aelrei/test_go_avito/internal/getters"
	"gitlab.com/Aelrei/test_go_avito/internal/http-server/accessHTTP"
	"gitlab.com/Aelrei/test_go_avito/internal/setters"
	"gitlab.com/Aelrei/test_go_avito/internal/storage"
	"gitlab.com/Aelrei/test_go_avito/internal/storage/accessDB"
	"io"
	"net/http"
	"strconv"
)

type S struct {
	db *sql.DB
}

func New(db *sql.DB) S {
	return S{db: db}
}

func (ss *S) GetUserBannerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tagID := r.URL.Query().Get("tag_id")
		featureID := r.URL.Query().Get("feature_id")
		useLastRevision := r.URL.Query().Get("use_last_revision")
		token := r.Header.Get("token")

		if useLastRevision == "" {
			useLastRevision = "false"
		} else {
			err := accessDB.ValidateLastRevision(useLastRevision)
			if err != nil {
				accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		err := accessDB.ValidateID(tagID)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		err = accessDB.ValidateID(featureID)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if useLastRevision == "true" {
			trueLastRevision(w, tagID, featureID, token, ss.db)
		} else if useLastRevision == "false" {
			falseLastRevision(w, tagID, featureID, token)
		}
	default:
		accessHTTP.SendErrorResponse(w, http.StatusBadRequest, "Not allowed Method")
	}
}

func (ss *S) GetPostBannersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tagID := r.URL.Query().Get("tag_id")
		featureID := r.URL.Query().Get("feature_id")
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		if tagID == "" && featureID != "" {
			err := accessDB.ValidateID(featureID)
			if err != nil {
				accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else if tagID != "" && featureID == "" {
			err := accessDB.ValidateID(tagID)
			if err != nil {
				accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else {
			err := accessDB.ValidateID(tagID)
			if err != nil {
				accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			err = accessDB.ValidateID(featureID)
			if err != nil {
				accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		err := accessDB.ValidateLimitOffset(limit)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		err = accessDB.ValidateLimitOffset(offset)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		banner, err := getters.GetAllBanners(tagID, featureID, limit, offset, ss.db)
		if banner == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				accessHTTP.SendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
			}
			return
		}
		jsonBytes, _ := json.MarshalIndent(banner, "", " ")
		jsonBytes = append(jsonBytes, '\n')
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonBytes)
		if err != nil {
			return
		}
	case "POST":
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

		newBannerID, err := setters.InsertBanner(data, ss.db)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		response := map[string]int{"banner_id": newBannerID}
		responseBody, err := json.Marshal(response)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		jsonBytes := append(responseBody, '\n')
		w.Write(jsonBytes)
	default:
		accessHTTP.SendErrorResponse(w, http.StatusBadRequest, "Not allowed Method")
	}
}

func (ss *S) PatchDeleteBannerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PATCH":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
			return
		}
		strID := r.PathValue("id")
		id, err := strconv.Atoi(strID)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = setters.ChangeInfoBanner(data, ss.db, id)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	case "DELETE":
		strID := r.PathValue("id")
		id, err := strconv.Atoi(strID)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		err = setters.DeleteBanner(ss.db, id)
		if err != nil {
			accessHTTP.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		accessHTTP.SendErrorResponse(w, http.StatusBadRequest, "Not allowed Method")
	}
}

func trueLastRevision(w http.ResponseWriter, tagID, featureID, token string, db *sql.DB) {
	banner, err := getters.GetBannerByTagAndFeature(tagID, featureID, db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			accessHTTP.SendErrorResponse(w, http.StatusInternalServerError, "StatusInternalServerError")
		}
		return
	}
	active, err := getters.GetActive(tagID, featureID, db)
	if active == false && token == "user_token" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonBytes := append([]byte(banner), '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		return
	}
}

func falseLastRevision(w http.ResponseWriter, tagID, featureID, token string) {
	banner, found := getters.GetCache(tagID, featureID)
	if found != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var data map[string]interface{}
	var ans storage.Banner
	err := json.Unmarshal(banner.([]byte), &ans)
	if ans.Active == false && token == "user_token" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
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
