package postgresql

import (
	"Avito_go/internal/storage"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"math/rand/v2"
	"strconv"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.postgresql.New"

	db, err := sql.Open("postgres", storage.PsqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	gen, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS banners (
            id SERIAL PRIMARY KEY,
            content JSONB NOT NULL,
            feature_id INT NOT NULL,
            active BOOLEAN NOT NULL DEFAULT TRUE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                                           
        );

        CREATE TABLE IF NOT EXISTS tags (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL
            
        );

        CREATE TABLE IF NOT EXISTS features (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL
        );

        CREATE TABLE IF NOT EXISTS banner_tag (
            banner_id INT,
            tag_id INT,
            PRIMARY KEY (banner_id, tag_id)
        );
    `)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", gen, err)
	}

	return &Storage{db: db}, nil

}

func UpdateStorage(storagePath string) (*Storage, error) {
	const fn = "storage.postgresql.Update"

	db, err := sql.Open("postgres", storage.PsqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	sqlStatementTags := `
		INSERT INTO tags (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING
		RETURNING id;`

	for id := 1; id < 1000; id++ {
		value := "tag" + strconv.Itoa(id)
		err = db.QueryRow(sqlStatementTags, id, value).Scan(&id)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			} else {
				fmt.Println("Error during request:", err)
				panic(err)
			}
		}

	}

	sqlStatementFeatures := `
		INSERT INTO features (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING
		RETURNING id;
		`

	for id := 1; id < 1000; id++ {
		value := "feature" + strconv.Itoa(id)
		err = db.QueryRow(sqlStatementFeatures, id, value).Scan(&id)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			} else {
				fmt.Println("Error during request:", err)
				panic(err)
			}
		}
	}

	sqlStatementBanners := `
    INSERT INTO banners (id, content, feature_id, active, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (id) DO NOTHING
    RETURNING id;
`

	for id := 1; id < 1000; id++ {
		content := storage.BannerContent{
			Title: "some_title" + strconv.Itoa(id),
			Text:  "some_text" + strconv.Itoa(id),
			URL:   "some_url" + strconv.Itoa(id),
		}
		contentJSON, _ := json.Marshal(content)
		for i := 0; i <= rand.IntN(3)+1; i++ {
			featureID := rand.IntN(7) + id - 2
			active := "true"
			createdAt := time.Now().Format("2006-01-02 15:04:05")
			updatedAt := time.Now().Format("2006-01-02 15:04:05")

			_, err := db.Exec(sqlStatementBanners, id, contentJSON, featureID, active, createdAt, updatedAt)
			if err != nil {
				fmt.Println("Error during request:", err)
				panic(err)
			}
		}
	}

	sqlBannerTag := `
	   INSERT INTO banner_tag (banner_id, tag_id)
	   VALUES ($1, $2)
	   ON CONFLICT (banner_id, tag_id) DO NOTHING;
	`
	for id := 1; id < 1000; id++ {
		for i := 0; i <= rand.IntN(3)+1; i++ {
			randId := rand.IntN(4) + id - 3
			_, err := db.Exec(sqlBannerTag, id, randId)
			if err != nil {
				fmt.Println("Error during request:", err)
				panic(err)
			}
		}
	}

	return &Storage{db: db}, nil
}
