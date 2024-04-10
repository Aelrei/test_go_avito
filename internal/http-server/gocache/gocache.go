package gocache

import (
	"Avito_go/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
)

var (
	Cah = cache.New(0, 0)
)

func LoadDataIntoCache() error {

	db, err := sql.Open("postgres", storage.PsqlInfo)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	query := `SELECT b.id, b.content, b.feature_id, t.id as tag_id, b.active, b.created_at, b.updated_at
		FROM features as f, banners as b, banner_tag as bt, tags as t
		WHERE t.id = bt.tag_id
		  AND bt.banner_id = b.id
		  AND f.id = b.feature_id;`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var banner storage.Banner
		err := rows.Scan(&banner.Id, &banner.Content, &banner.Feature_id, &banner.Tag_id, &banner.Active, &banner.Created_at, &banner.Updated_at)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		jsonData, err := json.Marshal(banner)
		if err != nil {
			return fmt.Errorf("failed to marshal data to JSON: %w", err)
		}
		Cah.Set(fmt.Sprintf("%d %d", banner.Tag_id, banner.Feature_id), jsonData, cache.NoExpiration) // NoExpiration - без истечения срока действия
	}

	return nil
}
