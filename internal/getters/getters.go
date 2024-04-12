package getters

import (
	"Avito_go/internal/getters/gocache"
	"Avito_go/internal/storage"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

func GetBannerByTagAndFeature(tagID, featureID string) (string, error) {
	db, err := sql.Open("postgres", storage.PsqlInfo)
	if err != nil {
		return "", fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	query := `
		SELECT b.content
        FROM  banners as b
            LEFT OUTER JOIN banner_tag as bt
        ON bt.banner_id = b.id
		  WHERE bt.tag_id = $1 AND b.feature_id = $2 AND b.active = true;
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

func GetAllBanners(tagID, featureID, limit, offset string) ([]*storage.AllBanner, error) {

	db, err := sql.Open("postgres", storage.PsqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	query := `
        SELECT b.id, b.content, b.feature_id, bt.tag_id, b.active, b.created_at, b.updated_at
        FROM  banners as b
        LEFT OUTER JOIN features as f
            ON b.feature_id = f.id
        INNER JOIN banner_tag as bt
            ON b.id = bt.banner_id
    `
	if featureID != "" && tagID == "" {
		query = query + ` AND b.feature_id = $1 `
	} else if tagID != "" && featureID == "" {
		query = query + ` AND bt.tag_id = $1 `
	} else if featureID != "" && tagID != "" {
		query = query + ` AND b.feature_id = $1 AND bt.tag_id = $2 `
	}

	if featureID != "" && tagID != "" {
		query = query + ` LIMIT $3 OFFSET $4;`
	} else {
		query = query + ` LIMIT $2 OFFSET $3;`
	}

	var rows *sql.Rows
	if featureID != "" && tagID != "" {
		rows, err = db.Query(query, featureID, tagID, limit, offset)
	} else if featureID != "" && tagID == "" {
		rows, err = db.Query(query, featureID, limit, offset)
	} else if tagID != "" && featureID == "" {
		rows, err = db.Query(query, tagID, limit, offset)
	} else {
		return nil, fmt.Errorf("at least one of featureID or tagID must be provided")
	}
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	bannerMap := make(map[int]*storage.AllBanner)

	for rows.Next() {
		var banner storage.Banner
		err := rows.Scan(&banner.Id, &banner.Content, &banner.Feature_id, &banner.Tag_id, &banner.Active, &banner.Created_at, &banner.Updated_at)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		if b, ok := bannerMap[banner.Id]; ok {
			b.Tag_ids = append(b.Tag_ids, banner.Tag_id)
		} else {
			newBanner := &storage.AllBanner{
				Id:         banner.Id,
				Content:    banner.Content,
				Feature_id: banner.Feature_id,
				Tag_ids:    []int{banner.Tag_id},
				Active:     banner.Active,
				Created_at: banner.Created_at,
				Updated_at: banner.Updated_at,
			}
			bannerMap[banner.Id] = newBanner
		}
	}

	var banners []*storage.AllBanner
	for _, b := range bannerMap {
		banners = append(banners, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return banners, nil
}

func GetCache(TagId, FeatureId string) (interface{}, error) {
	str := fmt.Sprintf("%s %s", TagId, FeatureId)
	cachedValue, found := gocache.Cah.Get(str)
	if !found {
		return nil, errors.New("value not found in cache")
	}
	return cachedValue, nil
}

func GetMaxBannerIdFromDB(db *sql.DB) (int, error) {
	var maxId int
	query := "SELECT MAX(id) FROM banners"
	err := db.QueryRow(query).Scan(&maxId)
	if err != nil {
		return 0, fmt.Errorf("error getting max Id from database: %v", err)
	}
	return maxId, nil
}

func GetMaxBannerFeatureIdFromDB(db *sql.DB) (int, error) {
	var maxFeatureId int
	query := "SELECT MAX(id) FROM features"
	err := db.QueryRow(query).Scan(&maxFeatureId)
	if err != nil {
		return 0, fmt.Errorf("error getting max Id from database: %v", err)
	}
	return maxFeatureId, nil
}

func GetMaxBannerTagIdFromDB(db *sql.DB) (int, error) {
	var maxTagId int
	query := "SELECT MAX(id) FROM tags"
	err := db.QueryRow(query).Scan(&maxTagId)
	if err != nil {
		return 0, fmt.Errorf("error getting max Id from database: %v", err)
	}
	return maxTagId, nil
}

func CheckBannerByTagFeature(tagID int, featureID int, db *sql.DB) (bool, error) {
	var exists bool
	err := db.QueryRow(`
        SELECT EXISTS (
            SELECT 1
            FROM banners AS b
            INNER JOIN banner_tag AS bt ON b.id = bt.banner_id
            WHERE b.feature_id = $1 AND bt.tag_id = $2
        )
    `, featureID, tagID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking tag_id validity: %v", err)
	}
	return exists, nil
}

func CheckBannerById(bannerId int, db *sql.DB) (bool, error) {
	var exists bool
	err := db.QueryRow(`
        SELECT EXISTS (
            SELECT 1
            FROM banners AS b
            WHERE b.id = $1
        )
    `, bannerId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking tag_id validity: %v", err)
	}
	return exists, nil
}