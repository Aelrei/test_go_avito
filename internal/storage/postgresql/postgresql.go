package postgresql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"gitlab.com/Aelrei/test_go_avito/internal/config"
	"gitlab.com/Aelrei/test_go_avito/internal/lib/logger/postgres"
	"gitlab.com/Aelrei/test_go_avito/internal/storage"
	"log/slog"
	"math/rand/v2"
	"os"
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
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            CONSTRAINT unique_banner_feature UNIQUE (id, feature_id)
                                           
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

	sqlStatementBanners := `
		INSERT INTO banners (content, feature_id, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING
		RETURNING id;
	`

	for id := 1; id < storage.NumberBanners; id++ {
		content := storage.BannerContent{
			Title: "some_title" + strconv.Itoa(id),
			Text:  "some_text" + strconv.Itoa(id),
			URL:   "some_url" + strconv.Itoa(id),
		}
		contentJSON, _ := json.Marshal(content)

		featureID := (id % 7) + 1

		createdAt := time.Now()
		updatedAt := time.Now()

		_, err := db.Exec(sqlStatementBanners, contentJSON, featureID, true, createdAt, updatedAt)
		if err != nil {
			fmt.Println("Error during request:", err)
			return nil, err
		}

		for i := 0; i < rand.IntN(3)+1; i++ {
			tagID := i + id

			sqlBannerTag := `
				INSERT INTO banner_tag (banner_id, tag_id)
				VALUES ($1, $2)
				ON CONFLICT (banner_id, tag_id) DO NOTHING;
			`
			_, err := db.Exec(sqlBannerTag, id, tagID)
			if err != nil {
				fmt.Println("Error during request:", err)
				return nil, err
			}
		}
	}

	sqlStatementFeatures := `
		INSERT INTO features (name)
		VALUES ($1)
		ON CONFLICT (id) DO NOTHING
		RETURNING id;
	`

	for id := 1; id <= (storage.NumberBanners / 2); id++ {
		value := "feature" + strconv.Itoa(id)
		_, err := db.Exec(sqlStatementFeatures, value)
		if err != nil {
			fmt.Println("Error during request:", err)
			return nil, err
		}
	}

	sqlStatementTags := `
		INSERT INTO tags (name)
		VALUES ($1)
		ON CONFLICT (id) DO NOTHING
		RETURNING id;
	`

	for id := 1; id < storage.NumberBanners; id++ {
		value := "tag" + strconv.Itoa(id)
		_, err := db.Exec(sqlStatementTags, value)
		if err != nil {
			fmt.Println("Error during request:", err)
			return nil, err
		}
	}

	return &Storage{db: db}, nil
}

func CheckPostgresDB(db *sql.DB, cfg *config.Config, log *slog.Logger) error {
	query := `
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_name IN ('banners', 'tags', 'features', 'banner_tag');
	`
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		fmt.Println("Failed to execute query:", err)
		return err
	}

	if count == 4 {
		log.Info("database exists")
	} else {
		storage, err := New(cfg.Storage)
		if err != nil {
			log.Error("failed to init storage", postgres.Err(err))
			os.Exit(1)
		} else {
			log.Info("success init storage")
		}

		storage, err = UpdateStorage(cfg.Storage)
		if err != nil {
			log.Error("failed to update storage", postgres.Err(err))
			os.Exit(1)
		} else {
			log.Info("success update storage")
		}
		_ = storage
	}
	return nil
}
