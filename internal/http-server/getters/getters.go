package getters

import (
	"Avito_go/internal/http-server/gocache"
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
	//fmt.Println(string(formattedJSON))
	return string(formattedJSON), nil
}

func GetAllBanners(tagID, featureID string) ([]*storage.AllBanner, error) {
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
          AND b.feature_id = $1;
    `

	rows, err := db.Query(query, featureID)
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

		// Проверяем, существует ли баннер с данным ID в карте
		if b, ok := bannerMap[banner.Id]; ok {
			// Если баннер уже существует, добавляем Tag_id в его массив Tag_ids
			b.Tag_ids = append(b.Tag_ids, banner.Tag_id)
		} else {
			// Иначе создаем новый баннер и добавляем его в карту
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

	// Преобразуем карту в список баннеров
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
