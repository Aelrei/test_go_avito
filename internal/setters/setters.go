package setters

import (
	"Avito_go/internal/getters"
	"Avito_go/internal/gocache"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func InsertBanner(data map[string]interface{}, db *sql.DB) (int, error) {
	contentData, ok := data["content"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("content field is missing or not a valid JSON object")
	}

	requiredFields := []string{"title", "text", "url"}
	for _, field := range requiredFields {
		if _, ok := contentData[field]; !ok {
			return 0, fmt.Errorf("%s field is missing in content", field)
		}
	}

	jsonData, err := json.Marshal(contentData)
	if err != nil {
		return 0, fmt.Errorf("error marshalling content: %v", err)
	}

	featureID, ok := data["feature_id"].(float64)
	if !ok || featureID <= 0 {
		return 0, fmt.Errorf("feature_id field is missing, not an integer or below zero")
	}
	intValue := int(featureID)

	active, ok := data["is_active"].(bool)
	if !ok {
		return 0, fmt.Errorf("is_active field is missing or not a boolean")
	}

	var maxID int
	maxID, err = getters.GetMaxBannerIdFromDB(db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		maxID = 0
		return 0, nil
	}
	maxID++

	tagIDs, ok := data["tag_ids"].([]interface{})
	if !ok {
		return 0, fmt.Errorf("tag_ids field is missing or not an array")
	}

	var tagIDArray []int
	for _, tagID := range tagIDs {
		tagIDFloat, ok := tagID.(float64)
		if !ok || tagIDFloat <= 0 {
			return 0, fmt.Errorf("tag_id is not an integer or below zero")
		}
		tagIDArray = append(tagIDArray, int(tagIDFloat))
	}

	now := time.Now()

	var id int
	err = db.QueryRow(`
        INSERT INTO banners (id, content, feature_id, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `, maxID, jsonData, intValue, active, now, now).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting banner: %v", err)
	}

	nameFeature := "feature" + strconv.Itoa(intValue)
	_, err = db.Exec(`
	INSERT INTO features (id, name)
	VALUES ($1, $2)
	ON CONFLICT (id) DO NOTHING
`, intValue, nameFeature)
	if err != nil {
		return 0, fmt.Errorf("error inserting into features table: %v", err)
	}

	for _, tagID := range tagIDArray {
		_, err := db.Exec(`
        INSERT INTO banner_tag (banner_id, tag_id)
        VALUES ($1, $2)
    `, id, tagID)
		if err != nil {
			return 0, fmt.Errorf("error inserting into banner_tag table: %v", err)
		}
	}

	for _, tagID := range tagIDArray {
		nameTag := "tag" + strconv.Itoa(tagID)
		_, err := db.Exec(`
        INSERT INTO tags (id, name)
        VALUES ($1, $2)
        ON CONFLICT (id) DO NOTHING
    `, tagID, nameTag)
		if err != nil {
			return 0, fmt.Errorf("error inserting into banner_tag table: %v", err)
		}
	}

	return id, nil
}

func ChangeInfoBanner(data map[string]interface{}, db *sql.DB, idNum int) error {
	maxID, err := getters.GetMaxBannerIdFromDB(db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		maxID = 0
		return nil
	}
	if idNum > maxID {
		return fmt.Errorf("too big id")
	} else if idNum <= 0 {
		return fmt.Errorf("id less than 1")
	}

	contentData, ok := data["content"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("content field is missing or not a valid JSON object")
	}

	requiredFields := []string{"title", "text", "url"}
	for _, field := range requiredFields {
		if _, ok := contentData[field]; !ok {

			return fmt.Errorf("%s field is missing in content", field)
		}
	}

	jsonData, err := json.Marshal(contentData)
	if err != nil {
		return fmt.Errorf("error marshalling content: %v", err)
	}

	maxFeatureId, err := getters.GetMaxBannerFeatureIdFromDB(db)
	if err != nil {
		return fmt.Errorf("can't count max(feature_id)")
	}
	featureID, ok := data["feature_id"].(float64)
	featureIDInt := int(featureID)
	if !ok || featureIDInt <= 0 || featureIDInt > maxFeatureId {
		return fmt.Errorf("feature_id field is missing, not an integer, below zero or more than max(feature_id) ")
	}

	active, ok := data["is_active"].(bool)
	if !ok {
		return fmt.Errorf("is_active field is missing or not a boolean")
	}

	tagIDs, ok := data["tag_ids"].([]interface{})
	if !ok {
		return fmt.Errorf("tag_ids field is missing or not an array")
	}

	maxTagID, err := getters.GetMaxBannerTagIdFromDB(db)
	if err != nil {
		return fmt.Errorf("can't count max(feature_id)")
	}
	var tagIDArray []int
	for _, tagID := range tagIDs {
		tagIDFloat, ok := tagID.(float64)
		tagIdInt := int(tagIDFloat)
		if !ok || tagIdInt <= 0 || tagIdInt > maxTagID {
			return fmt.Errorf("tag_id is not an integer or below zero")
		}
		check, err := getters.CheckBannerByTagFeature(tagIdInt, featureIDInt, idNum, db)
		if check || err != nil {
			return fmt.Errorf("may cause duplicates")
		}
		tagIDArray = append(tagIDArray, int(tagIDFloat))
	}

	updateTime := time.Now()

	_, err = db.Exec(`
        UPDATE banners
        SET content = $1, feature_id = $2, active = $3, updated_at = $4
        WHERE id = $5
    `, jsonData, featureIDInt, active, updateTime, idNum)
	if err != nil {
		return fmt.Errorf("error updating banner: %v", err)
	}

	_, err = db.Exec("DELETE FROM banner_tag WHERE banner_id = $1", idNum)
	if err != nil {
		return fmt.Errorf("error deleting old banner_tag associations: %v", err)
	}

	for _, tagID := range tagIDArray {
		_, err := db.Exec("INSERT INTO banner_tag (banner_id, tag_id) VALUES ($1, $2)", idNum, tagID)
		if err != nil {
			return fmt.Errorf("error inserting new banner_tag association: %v", err)
		}
	}

	err = gocache.LoadDataIntoCache()
	if err != nil {
		return err
	}

	return nil
}

func DeleteBanner(db *sql.DB, idNum int) error {
	maxID, err := getters.GetMaxBannerIdFromDB(db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		maxID = 0
		return nil
	}
	if idNum > maxID {
		return fmt.Errorf("too big id")
	} else if idNum <= 0 {
		return fmt.Errorf("id less than 1")
	}

	check, err := getters.CheckBannerById(idNum, db)
	if !check || err != nil {
		return fmt.Errorf("no such banner")
	}

	_, err = db.Exec("DELETE FROM banners WHERE id = $1", idNum)
	if err != nil {
		return fmt.Errorf("failed to delete banner: %v", err)
	}

	_, err = db.Exec("DELETE FROM banner_tag WHERE banner_id = $1", idNum)
	if err != nil {
		return fmt.Errorf("failed to delete banner: %v", err)
	}

	return nil
}
