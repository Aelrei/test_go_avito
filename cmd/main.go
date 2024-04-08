package main

import "C"
import (
	"Avito_go/internal/config"
	"Avito_go/internal/http-server/gocache"
	"Avito_go/internal/http-server/handlers"
	"Avito_go/internal/lib/logger/postgres"
	"Avito_go/internal/storage/postgresql"

	"fmt"
	"log/slog"
	"net/http"
	"os"
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

	//gocache.Cah = cache.New(0, 0)

	if err := gocache.LoadDataIntoCache(); err != nil {
		fmt.Println("Failed to load data into cache:", err)
		return
	}

	cachedData, found := gocache.Cah.Get("4 438")
	if found {
		fmt.Println("Data found in cache:")
		fmt.Println(string(cachedData.([]byte)))
	} else {
		fmt.Println("Data not found in cache")
	}

	http.HandleFunc("/user_banner", handlers.GetUserBanner)
	err = http.ListenAndServe(":8085", nil)
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
