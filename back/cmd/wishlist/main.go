package main

import (
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/wishlist/back/internal/config"
	"github.com/iamvkosarev/wishlist/back/internal/http-server/router"
	"github.com/iamvkosarev/wishlist/back/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()
	log, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logger: %v\n", err)
		return
	}

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		return
	}

	log.Info("init server", slog.String("address", cfg.HTTPServer.Address))
	http.ListenAndServe(cfg.HTTPServer.Address, router.New(log, storage, cfg.SSOURL))
}
