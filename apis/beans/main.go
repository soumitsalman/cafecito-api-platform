// @title 			Beans News API & MCP
// @version 		1.0
// @description Beans finds and verifies what publishers published. It returns citable Articles, Source metadata, attention-ranked feeds, similar publisher reading, external Article mentions, and normalized filter discovery.
// @schemes 		https
// @license.name 	MIT
// @contact.name 	Project Cafecito
// @contact.url  	https://cafecito.tech
// @contact.email 	soumitsrah@cafecito.tech
package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/soumitsalman/cafecito-api-platform/apis/beans/db"
	_ "github.com/soumitsalman/cafecito-api-platform/apis/beans/docs"
	r "github.com/soumitsalman/cafecito-api-platform/apis/beans/router"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared/config"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared/embedding"
)

const (
	DEFAULT_PORT              = "8080"
	DEFAULT_EMBEDDER_BASE_URL = "http://localhost:10000"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn_str := config.GetEnv("PG_CONNECTION_STRING", "", true)
	beansack := db.NewPGSack(ctx, conn_str)
	defer beansack.Close()

	api := r.NewRouter(
		beansack,
		embedding.NewHTTPEmbedder(
			config.GetEnv("EMBEDDER_BASE_URL", DEFAULT_EMBEDDER_BASE_URL, true),
			config.GetEnv("EMBEDDER_API_KEY", "", false),
			config.GetEnv("EMBEDDER_MODEL", "", false),
		),
		config.ParseAPIKeys(os.Getenv("API_KEY")),
	)

	port := config.GetEnv("PORT", DEFAULT_PORT, false)
	addr := "0.0.0.0:" + port
	log.Info().Str("module", "MAIN").Str("addr", addr).Msg("Routes Initialized. Server starting")

	if err := api.Run(addr); err != nil {
		log.Fatal().Str("module", "MAIN").Err(err).Msg("server error")
	}
}
