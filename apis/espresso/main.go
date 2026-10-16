package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/db"
	_ "github.com/soumitsalman/cafecito-api-platform/apis/espresso/docs"
	r "github.com/soumitsalman/cafecito-api-platform/apis/espresso/router"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared/config"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared/embedding"
)

const (
	DEFAULT_PORT              = "8080"
	DEFAULT_EMBEDDER_BASE_URL = "http://localhost:10000"
)

func main() {
	if err := run(); err != nil {
		log.Error().Err(err).Msg("API stopped")
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	responseCache, cache_err := shared.NewResponseCache()
	if cache_err != nil {
		return cache_err
	}

	conn_str := config.GetEnv("PG_CONNECTION_STRING", "", true)
	cupboard := db.NewCupboard(ctx, conn_str)
	defer cupboard.Close()

	api := r.NewRouter(
		cupboard,
		embedding.NewHTTPEmbedder(
			config.GetEnv("EMBEDDER_BASE_URL", DEFAULT_EMBEDDER_BASE_URL, true),
			config.GetEnv("EMBEDDER_API_KEY", "", false),
			config.GetEnv("EMBEDDER_MODEL", "", false),
		),
		config.ParseAPIKeys(os.Getenv("API_KEY")),
		responseCache,
	)

	port := config.GetEnv("PORT", DEFAULT_PORT, false)
	addr := "0.0.0.0:" + port
	log.Info().Str("module", "MAIN").Str("addr", addr).Msg("Routes Initialized. Server starting")

	return shared.Serve(api, addr)
}
