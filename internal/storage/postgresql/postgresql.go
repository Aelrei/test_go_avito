package postgresql

import (
	"Avito_go/internal/storage"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

type Storage struct {
	db *sql.DB
}

var (
	host     = os.Getenv("DB_HOST")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
	port     = os.Getenv("DB_PORT")
)

func New(storagePath string) (*Storage, error) {
	const fn = "storage.postgresql.New"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
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
--             FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE,
--             FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
        );
    `)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", gen, err)
	}

	return &Storage{db: db}, nil

}

func UpdateStorage(storagePath string) (*Storage, error) {
	const fn = "storage.postgresql.Update"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
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
		content, _ := json.Marshal("content" + strconv.Itoa(id))
		featureID := rand.IntN(1000) + 1
		active := "true"
		createdAt := time.Now().Format("2006-01-02 15:04:05")
		updatedAt := time.Now().Format("2006-01-02 15:04:05")

		_, err := db.Exec(sqlStatementBanners, id, content, featureID, active, createdAt, updatedAt)
		if err != nil {
			fmt.Println("Error during request:", err)
			panic(err)
		}
	}

	sqlBannerTag := `
	   INSERT INTO banner_tag (banner_id, tag_id)
	   VALUES ($1, $2)
	   ON CONFLICT (banner_id, tag_id) DO NOTHING;
	`
	for id := 1; id < 1000; id++ {
		for i := 0; i <= rand.IntN(3)+1; i++ {
			randId := rand.IntN(1000) + 1
			_, err := db.Exec(sqlBannerTag, id, randId)
			if err != nil {
				fmt.Println("Error during request:", err)
				panic(err)
			}
		}
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrTagNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}
