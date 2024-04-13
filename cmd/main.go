package main

import (
	"database/sql"
	"fmt"
	"gitlab.com/Aelrei/test_go_avito/internal/config"
	"gitlab.com/Aelrei/test_go_avito/internal/gocache"
	"gitlab.com/Aelrei/test_go_avito/internal/http-server/accessHTTP"
	"gitlab.com/Aelrei/test_go_avito/internal/http-server/handlers"
	"gitlab.com/Aelrei/test_go_avito/internal/lib/logger"
	"gitlab.com/Aelrei/test_go_avito/internal/storage"
	"gitlab.com/Aelrei/test_go_avito/internal/storage/postgresql"
	"net/http"
	"time"
)

func main() {

	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting...")

	db, err := sql.Open("postgres", storage.PsqlInfo)
	if err != nil {
		log.Warn("Failed to connect to the database:", err)
		return
	}
	defer db.Close()

	err = postgresql.CheckPostgresDB(db, cfg, log)
	if err != nil {
		log.Warn("Failed to connect to the database:", err)
		return
	}

	if err := gocache.LoadDataIntoCache(db); err != nil {
		log.Warn("failed to load data into cache:", err)
		return
	}
	log.Info("success upload cache")

	go func() {
		ticker := time.Tick(5 * time.Minute)
		for range ticker {
			gocache.Cache.Flush()
			if err := gocache.LoadDataIntoCache(db); err != nil {
				log.Warn("failed to update cache")
			} else {
				log.Info("success update cache")
			}
		}
	}()

	S := handlers.New(db)
	router := http.NewServeMux()

	router.Handle("/user_banner", accessHTTP.AuthMiddlewareUserAdmin(http.HandlerFunc(S.GetUserBannerHandler)))
	router.Handle("/banner", accessHTTP.AuthMiddlewareAdmin(http.HandlerFunc(S.GetPostBannersHandler)))
	router.Handle("/banner/{id}", accessHTTP.AuthMiddlewareAdmin(http.HandlerFunc(S.PatchDeleteBannerHandler)))

	address := cfg.HTTPServer.Address
	err = http.ListenAndServe(address, router)
	if err != nil {
		fmt.Println("Error: ", err)
	}

}
