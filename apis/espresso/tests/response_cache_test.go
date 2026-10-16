package espressoapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/router"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared"
	"github.com/stretchr/testify/require"
)

type countingCacheEmbedder struct{ calls atomic.Int32 }

func (embedder *countingCacheEmbedder) EmbedQuery(ctx context.Context, q string) []float32 {
	embedder.calls.Add(1)
	return append([]float32(nil), test_query_embedding...)
}
func (embedder *countingCacheEmbedder) EmbedDocuments(ctx context.Context, docs []string) [][]float32 {
	return nil
}
func (embedder *countingCacheEmbedder) Close() error { return nil }
func cacheGET(engine http.Handler, path, key string) *httptest.ResponseRecorder {
	request := httptest.NewRequest("GET", path, nil)
	request.Header.Set("X-API-KEY", key)
	writer := httptest.NewRecorder()
	engine.ServeHTTP(writer, request)
	return writer
}
func cacheFirstID(t *testing.T, engine http.Handler, path string) string {
	t.Helper()
	response := cacheGET(engine, path, "cache-test")
	require.Equal(t, 200, response.Code, response.Body.String())
	var page struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &page))
	require.NotEmpty(t, page.Data)
	return page.Data[0].ID
}

func TestResponseCacheAllPublicRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database := setupTestDB(t)
	defer database.Close()
	embedder := &countingCacheEmbedder{}
	t.Setenv("RESPONSE_CACHE_ENABLED", "true")
	t.Setenv("RESPONSE_CACHE_PATH", "")
	t.Setenv("RESPONSE_CACHE_TTL", "10m")
	responseCache, err := shared.NewResponseCache()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, responseCache.Store.Flush()) })
	engine := router.NewRouter(database, embedder, map[string]string{"X-API-KEY": "cache-test"}, responseCache)
	event := cacheFirstID(t, engine, "/events")
	signal := cacheFirstID(t, engine, "/signals")
	source := cacheFirstID(t, engine, "/sources")
	paths := []string{"/tags", "/entities", "/regions", "/event-types", "/sources", "/sources/" + source, "/events", "/events/" + event, "/events/" + event + "/signals", "/events/" + event + "/evidence", "/signals", "/signals/" + signal, "/signals/" + signal + "/events", "/events?response_type=yaml", "/events?response_type=toon"}
	warmed := make(map[string]string)
	content_types := make(map[string]string)
	for _, path := range paths {
		response := cacheGET(engine, path, "cache-test")
		require.Equal(t, 200, response.Code, path+": "+response.Body.String())
		warmed[path] = response.Body.String()
		content_types[path] = response.Header().Get("Content-Type")
	}
	defaults := cacheGET(engine, "/events?limit=20&score_threshold=0.6&response_type=json", "cache-test")
	require.NotEqual(t, warmed["/events"], defaults.Body.String())
	first := cacheGET(engine, "/events?q=technology&score_threshold=0.1&limit=20", "cache-test")
	require.Equal(t, 200, first.Code, first.Body.String())
	require.EqualValues(t, 1, embedder.calls.Load())
	second := cacheGET(engine, "/events?q=technology&score_threshold=0.1&limit=20", "cache-test")
	require.Equal(t, first.Body.String(), second.Body.String())
	require.EqualValues(t, 1, embedder.calls.Load())
	database.Close()
	for _, path := range paths {
		response := cacheGET(engine, path, "cache-test")
		require.Equal(t, 200, response.Code, path)
		require.Equal(t, warmed[path], response.Body.String(), path)
		require.Equal(t, content_types[path], response.Header().Get("Content-Type"))
	}
	require.Equal(t, 401, cacheGET(engine, "/events", "wrong").Code)
	require.Equal(t, 400, cacheGET(engine, "/events?cursor=invalid", "cache-test").Code)
	require.Equal(t, 400, cacheGET(engine, "/events?response_type=xml", "cache-test").Code)
	require.NotEqual(t, 200, cacheGET(engine, "/events?ignored=value", "cache-test").Code)
	require.Equal(t, 200, cacheGET(engine, "/health", "").Code)
}
