package gobeansack_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/soumitsalman/cafecito-api-platform/apis/beans/router"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared/embedding"
	"github.com/stretchr/testify/require"
)

type countingCacheEmbedder struct {
	embedding.Embedder
	calls atomic.Int32
}

func (embedder *countingCacheEmbedder) EmbedQuery(ctx context.Context, q string) []float32 {
	embedder.calls.Add(1)
	return embedder.Embedder.EmbedQuery(ctx, q)
}
func cacheGET(engine http.Handler, path, key string) *httptest.ResponseRecorder {
	request := httptest.NewRequest("GET", path, nil)
	request.Header.Set("X-API-KEY", key)
	writer := httptest.NewRecorder()
	engine.ServeHTTP(writer, request)
	return writer
}

func TestResponseCacheAllPublicRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database := setupTestDB()
	defer database.Close()
	embedder := &countingCacheEmbedder{Embedder: setupTestEmbedder()}
	t.Setenv("RESPONSE_CACHE_ENABLED", "true")
	t.Setenv("RESPONSE_CACHE_PATH", "")
	t.Setenv("RESPONSE_CACHE_TTL", "10m")
	responseCache, err := shared.NewResponseCache()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, responseCache.Store.Flush()) })
	engine := router.NewRouter(database, embedder, map[string]string{"X-API-KEY": "cache-test"}, responseCache)
	sources := cacheGET(engine, "/sources", "cache-test")
	require.Equal(t, 200, sources.Code, sources.Body.String())
	stories := cacheGET(engine, "/stories", "cache-test")
	require.Equal(t, 200, stories.Code, stories.Body.String())
	var story_page struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(stories.Body.Bytes(), &story_page))
	require.NotEmpty(t, story_page.Data)
	article := fixtureArticleIDs[0].String()
	story := story_page.Data[0].ID
	paths := []string{"/categories", "/entities", "/regions", "/sentiments", "/sources", "/sources/" + fixtureSourceTech.String(), "/articles/search", "/articles/latest", "/articles/trending", "/articles/" + article, "/articles/" + article + "/similar", "/articles/" + article + "/mentions", "/news/top-headlines", "/news/latest", "/news/trending", "/stories", "/stories/" + story, "/stories/" + story + "/articles"}
	warmed := make(map[string]string)
	for _, path := range paths {
		response := cacheGET(engine, path, "cache-test")
		require.Equal(t, 200, response.Code, path+": "+response.Body.String())
		warmed[path] = response.Body.String()
	}
	// Raw URI keys deliberately distinguish explicit defaults from omission.
	defaults := cacheGET(engine, "/articles/latest?limit=20&full_content=false", "cache-test")
	require.NotEqual(t, warmed["/articles/latest"], defaults.Body.String())
	first := cacheGET(engine, "/articles/search?q=technology&score_threshold=0.1&limit=20", "cache-test")
	require.Equal(t, 200, first.Code, first.Body.String())
	calls_after_miss := embedder.calls.Load()
	require.Positive(t, calls_after_miss)
	second := cacheGET(engine, "/articles/search?q=technology&score_threshold=0.1&limit=20", "cache-test")
	require.Equal(t, first.Body.String(), second.Body.String())
	require.Equal(t, calls_after_miss, embedder.calls.Load())
	// Distinct projection, identity, and unknown parameters must not reuse a success.
	projected := cacheGET(engine, "/articles/"+article+"?full_content=true", "cache-test")
	require.Equal(t, 200, projected.Code, projected.Body.String())
	require.NotEqual(t, warmed["/articles/"+article], projected.Body.String())
	ignored := cacheGET(engine, "/sources?ignored=value", "cache-test")
	require.Equal(t, 200, ignored.Code)
	require.NotEqual(t, warmed["/sources"], ignored.Body.String())
	require.NotEqual(t, 200, cacheGET(engine, "/sources/"+strings.ReplaceAll(fixtureSourceTech.String(), "a", "b"), "cache-test").Code)
	// Closing storage proves every warmed route actually skips its database work.
	database.Close()
	for _, path := range paths {
		response := cacheGET(engine, path, "cache-test")
		require.Equal(t, 200, response.Code, path)
		require.Equal(t, warmed[path], response.Body.String(), path)
	}
	require.Equal(t, 401, cacheGET(engine, "/sources", "wrong").Code)
	require.Equal(t, 400, cacheGET(engine, "/sources?cursor=invalid", "cache-test").Code)
	require.Equal(t, 400, cacheGET(engine, "/news/latest?content_type=blog", "cache-test").Code)
	require.Equal(t, 200, cacheGET(engine, "/health", "").Code)
}
