package main

import (
	"context"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/cupboard"
	_ "github.com/soumitsalman/cafecito-api-platform/apis/espresso/docs"
	r "github.com/soumitsalman/cafecito-api-platform/apis/espresso/router"
	"github.com/soumitsalman/cafecito-api-platform/apis/internal/config"
	"github.com/soumitsalman/cafecito-api-platform/apis/internal/embedding"
)

const (
	DEFAULT_PORT              = "8080"
	DEFAULT_EMBEDDER_BASE_URL = "http://localhost:10000"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connStr := config.GetEnv("PG_CONNECTION_STRING", "", true)
	db := cupboard.NewCupboard(ctx, connStr)
	defer db.Close()

	// determine concurrency limit from environment
	maxStr := config.GetEnv("MAX_CONCURRENT_REQUESTS", "", false)
	max_requests, err := strconv.Atoi(maxStr)
	if err != nil && max_requests < 0 {
		max_requests = 0
	}

	api := r.NewRouter(
		db,
		embedding.NewHTTPEmbedder(
			config.GetEnv("EMBEDDER_BASE_URL", DEFAULT_EMBEDDER_BASE_URL, true),
			config.GetEnv("EMBEDDER_API_KEY", "", false),
			config.GetEnv("EMBEDDER_MODEL", "", false),
		),
		config.ParseAPIKeys(os.Getenv("API_KEY")),
		max_requests,
	)

	port := config.GetEnv("PORT", DEFAULT_PORT, false)
	addr := "0.0.0.0:" + port
	log.Info().Str("module", "MAIN").Str("addr", addr).Msg("Routes Initialized. Server starting")

	if err := api.Run(addr); err != nil {
		log.Fatal().Str("module", "MAIN").Err(err).Msg("server error")
	}
}
