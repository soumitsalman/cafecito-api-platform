// @title 			Beans News API & MCP
// @version 		0.8
// @description 	MCP-ready news and blog intelligence over RSS-sourced articles, semantic enrichment, and propagation tracking.
// @description 	A **bean** is one article or post keyed by canonical URL. Records include publisher metadata, summary/full text, publish timestamp, inferred categories, regions, entities, sentiments, and optional social trend metrics.
// @description 	Agent workflow: (1) listCategories, listEntities, listRegions to discover exact filter values; (2) searchArticles for full-corpus retrieval; (3) getLatestArticles, getTrendingArticles, or getTopHeadlines for feed-style monitoring; (4) getPublishers to resolve source IDs; (5) getArticlePropagation or postArticlePropagation to check story spread.
// @description 	Conventions: Auth is optional at the backend but API-key protected through the gateway. Pagination uses `limit` default 16 max 128 and `offset` default 0. Empty result sets return HTTP 204, not an error. Use fuzzy `tags` when spelling is uncertain, exact `categories`/`regions`/`entities` after discovery, and `q` + `acc` for semantic vector search.
// @schemes 		https
// @license.name 	MIT
// @contact.name 	Project Cafecito
// @contact.url  	http://cafecito.tech
// @contact.email 	soumitsrah@cafecito.tech
package main

import (
	"context"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/soumitsalman/cafecito-api-platform/apis/beans/db"
	_ "github.com/soumitsalman/cafecito-api-platform/apis/beans/docs"
	r "github.com/soumitsalman/cafecito-api-platform/apis/beans/router"
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
	beansack := db.NewPGSack(ctx, connStr)
	defer beansack.Close()

	// determine concurrency limit from environment
	maxStr := config.GetEnv("MAX_CONCURRENT_REQUESTS", "", false)
	max_requests, err := strconv.Atoi(maxStr)
	if err != nil && max_requests < 0 {
		max_requests = 0
	}

	api := r.NewRouter(
		beansack,
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
