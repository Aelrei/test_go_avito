package storage

import (
	"errors"
	"os"
)

var (
	Host     = os.Getenv("DB_HOST")
	User     = os.Getenv("DB_USER")
	Password = os.Getenv("DB_PASSWORD")
	Dbname   = os.Getenv("DB_NAME")
	Port     = os.Getenv("DB_PORT")
)

type Banner struct {
	Id         int    `json:"banner_id"`
	Content    string `json:"content"`
	Feature_id int    `json:"feature_id"`
	Tag_id     int    `json:"tag_id"`
	Active     bool   `json:"is_active"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

var (
	ErrTagNotFound     = errors.New("tag not found")
	ErrFeatureNotFound = errors.New("feature not found")
)
