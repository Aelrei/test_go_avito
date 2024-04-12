package storage

import (
	"errors"
	"fmt"
	"os"
)

var NumberBanners = 1000

var PsqlInfo = fmt.Sprintf("host=%s port=%s user=%s "+
	"password=%s dbname=%s sslmode=disable",
	Host, Port, User, Password, Dbname)

var (
	Host     = os.Getenv("DB_HOST")
	User     = os.Getenv("DB_USER")
	Password = os.Getenv("DB_PASSWORD")
	Dbname   = os.Getenv("DB_NAME")
	Port     = os.Getenv("DB_PORT")
)

type BannerContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

type Banner struct {
	Id         int    `json:"banner_id,omitempty"`
	Content    string `json:"content,omitempty"`
	Feature_id int    `json:"feature_id,omitempty"`
	Tag_id     int    `json:"tag_id,omitempty"`
	Active     bool   `json:"is_active,omitempty"`
	Created_at string `json:"created_at,omitempty"`
	Updated_at string `json:"updated_at,omitempty"`
}

type AllBanner struct {
	Id         int    `json:"banner_id"`
	Content    string `json:"content"`
	Feature_id int    `json:"feature_id"`
	Tag_ids    []int  `json:"tag_ids"`
	Active     bool   `json:"is_active"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

var (
	ErrTagNotFound     = errors.New("tag not found")
	ErrFeatureNotFound = errors.New("feature not found")
)
