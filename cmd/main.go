package main

import "C"
import (
	"Avito_go/internal/config"
	"Avito_go/internal/getters/gocache"
	"Avito_go/internal/http-server/accessHTTP"
	"Avito_go/internal/http-server/handlers"
	"Avito_go/internal/lib/logger/postgres"
	storage2 "Avito_go/internal/storage"
	"Avito_go/internal/storage/postgresql"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting test")
	log.Debug("debug test")

	db, err := sql.Open("postgres", storage2.PsqlInfo)
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}
	defer db.Close()

	query := `
		SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_name IN ('banners', 'tags', 'features', 'banner_tag');
	`
	var count int
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		fmt.Println("Failed to execute query:", err)
		return
	}

	if count == 4 {
		log.Info("database exist")
	} else {
		storage, err := postgresql.New(cfg.Storage)
		if err != nil {
			log.Error("failed to init storage", postgres.Err(err))
			os.Exit(1)
		} else {
			log.Info("success init storage")
		}

		storage, err = postgresql.UpdateStorage(cfg.Storage)
		if err != nil {
			log.Error("failed to update storage", postgres.Err(err))
			os.Exit(1)
		} else {
			log.Info("success update storage")
		}
		_ = storage
	}

	if err := gocache.LoadDataIntoCache(); err != nil {
		log.Info("failed to load data into cache:", err)
		return
	}
	log.Info("success upload cache")

	go func() {
		ticker := time.Tick(5 * time.Minute)
		for {
			<-ticker
			if err := gocache.LoadDataIntoCache(); err != nil {
				log.Warn("failed to update cache")
			} else {
				log.Info("success update cache")
			}
		}
	}()

	router := http.NewServeMux()

	router.Handle("/user_banner", accessHTTP.AuthMiddlewareUserAdmin(http.HandlerFunc(handlers.GetUserBannerHandler)))
	router.Handle("/banner", accessHTTP.AuthMiddlewareAdmin(http.HandlerFunc(handlers.GetAllBannersHandler)))

	//http.HandleFunc("/user_banner", handlers.GetUserBannerHandler)
	err = http.ListenAndServe(":8085", router)
	if err != nil {
		fmt.Println("Error: ", err)
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
