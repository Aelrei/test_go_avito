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

	Cah.Flush()

	db, err := sql.Open("postgres", storage.PsqlInfo)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	query := `SELECT
    b.id AS banner_id,
    b.content AS banner_content,
    b.feature_id AS banner_feature_id,
    bt.tag_id AS tag_id,
    b.active AS banner_active,
    b.created_at AS banner_created_at,
    b.updated_at AS banner_updated_at
FROM
    banners AS b
JOIN
    banner_tag AS bt ON bt.banner_id = b.id;`

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
		Cah.Set(fmt.Sprintf("%d %d", banner.Tag_id, banner.Feature_id), jsonData, cache.NoExpiration)
	}

	return nil
}