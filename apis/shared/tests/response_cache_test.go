package internal_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared"
	"github.com/stretchr/testify/require"
)

func newMemoryResponseStore(t *testing.T) *shared.ResponseCache {
	t.Helper()
	t.Setenv("RESPONSE_CACHE_ENABLED", "true")
	t.Setenv("RESPONSE_CACHE_PATH", "")
	t.Setenv("RESPONSE_CACHE_TTL", "")
	store, err := shared.NewResponseCache()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, store.Store.Flush()) })
	return store
}
func cacheRequest(engine http.Handler, path string) *httptest.ResponseRecorder {
	writer := httptest.NewRecorder()
	engine.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, path, nil))
	return writer
}
func TestResponseCacheSelection(t *testing.T) {
	store := newMemoryResponseStore(t)
	require.IsType(t, &persistence.InMemoryStore{}, store.Store)
	require.Empty(t, store.Path)
	require.Equal(t, 10*time.Minute, store.TTL)
	// Removed disk/size settings must not influence the memory-backed selection.
	t.Setenv("RESPONSE_CACHE_DIR", "/does/not/exist")
	t.Setenv("RESPONSE_CACHE_MAX_ENTRY_BYTES", "invalid")
	_, err := shared.NewResponseCache()
	require.NoError(t, err)
	for _, path := range []string{"/tmp/cache", "https://example.com", "redis://", "redis://localhost:invalid"} {
		t.Setenv("RESPONSE_CACHE_PATH", path)
		_, err := shared.NewResponseCache()
		require.Error(t, err)
		require.NotContains(t, err.Error(), path)
	}
	t.Setenv("RESPONSE_CACHE_ENABLED", "false")
	disabled, err := shared.NewResponseCache()
	require.NoError(t, err)
	require.Nil(t, disabled)
	t.Setenv("RESPONSE_CACHE_ENABLED", "invalid")
	_, err = shared.NewResponseCache()
	require.Error(t, err)
}
func TestResponseCacheTTLConfiguration(t *testing.T) {
	t.Setenv("RESPONSE_CACHE_ENABLED", "true")
	t.Setenv("RESPONSE_CACHE_PATH", "")
	for _, test := range []struct {
		raw string
		ttl time.Duration
	}{
		{"", 10 * time.Minute}, {"  ", 10 * time.Minute},
		{"30s", 30 * time.Second}, {"1h", time.Hour}, {" 2m ", 2 * time.Minute},
	} {
		t.Run(test.raw, func(t *testing.T) {
			t.Setenv("RESPONSE_CACHE_TTL", test.raw)
			responseCache, err := shared.NewResponseCache()
			require.NoError(t, err)
			require.Equal(t, test.ttl, responseCache.TTL)
			require.NoError(t, responseCache.Store.Flush())
		})
	}
	for _, raw := range []string{"invalid", "10", "0", "0s", "-1m", "999999999999999999h"} {
		t.Run(raw, func(t *testing.T) {
			t.Setenv("RESPONSE_CACHE_TTL", raw)
			responseCache, err := shared.NewResponseCache()
			require.ErrorContains(t, err, "RESPONSE_CACHE_TTL")
			require.Nil(t, responseCache)
		})
	}
}

func TestResponseCacheDisabledWrap(t *testing.T) {
	t.Setenv("RESPONSE_CACHE_ENABLED", "false")
	t.Setenv("RESPONSE_CACHE_PATH", "invalid")
	t.Setenv("RESPONSE_CACHE_TTL", "invalid")
	disabled, err := shared.NewResponseCache()
	require.NoError(t, err)
	require.Nil(t, disabled)
	for _, responseCache := range []*shared.ResponseCache{disabled, {}} {
		var calls int
		engine := gin.New()
		engine.GET("/data", responseCache.Wrap(func(c *gin.Context) {
			calls++
			c.String(200, "ok")
		}))
		for i := 0; i < 2; i++ {
			require.Equal(t, "ok", cacheRequest(engine, "/data").Body.String())
		}
		require.Equal(t, 2, calls)
	}
}

func TestResponseCacheMemoryHitExpiryAndRawQueryKeys(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("RESPONSE_CACHE_ENABLED", "true")
	t.Setenv("RESPONSE_CACHE_PATH", "")
	t.Setenv("RESPONSE_CACHE_TTL", "100ms")
	store, err := shared.NewResponseCache()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, store.Store.Flush()) })
	// Initialization captures configuration; later environment changes do not alter it.
	t.Setenv("RESPONSE_CACHE_TTL", "1h")
	t.Setenv("RESPONSE_CACHE_PATH", "invalid")
	var calls atomic.Int32
	engine := gin.New()
	engine.GET("/data/:id", store.Wrap(func(c *gin.Context) {
		c.Header("X-Example", "retained")
		c.JSON(200, gin.H{"call": calls.Add(1), "id": c.Param("id")})
	}))
	first := cacheRequest(engine, "/data/a?q=x&limit=20")
	repeated := cacheRequest(engine, "/data/a?q=x&limit=20")
	require.Equal(t, first.Body.String(), repeated.Body.String())
	require.EqualValues(t, 1, calls.Load())
	require.Equal(t, "retained", repeated.Header().Get("X-Example"))
	// RequestURI semantics intentionally distinguish query ordering and defaults.
	require.NotEqual(t, first.Body.String(), cacheRequest(engine, "/data/a?limit=20&q=x").Body.String())
	require.NotEqual(t, first.Body.String(), cacheRequest(engine, "/data/a?q=x").Body.String())
	require.NotEqual(t, first.Body.String(), cacheRequest(engine, "/data/b?q=x&limit=20").Body.String())
	before := calls.Load()
	require.Eventually(t, func() bool { cacheRequest(engine, "/data/a?q=x&limit=20"); return calls.Load() > before }, time.Second, 20*time.Millisecond)
}
func TestResponseCacheLibraryStatusAndAbortBehavior(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := newMemoryResponseStore(t)
	var created_calls, aborted_calls, error_calls atomic.Int32
	engine := gin.New()
	engine.GET("/created", store.Wrap(func(c *gin.Context) {
		created_calls.Add(1)
		c.SetCookie("example", "value", 60, "/", "", false, true)
		c.String(201, "created")
	}))
	engine.GET("/aborted", store.Wrap(func(c *gin.Context) { aborted_calls.Add(1); c.String(200, "aborted"); c.Abort() }))
	engine.GET("/error", store.Wrap(func(c *gin.Context) { error_calls.Add(1); c.String(500, "error") }))
	for i := 0; i < 2; i++ {
		response := cacheRequest(engine, "/created")
		require.Equal(t, 201, response.Code)
		require.NotEmpty(t, response.Header().Get("Set-Cookie"))
		require.Equal(t, 200, cacheRequest(engine, "/aborted").Code)
		require.Equal(t, 500, cacheRequest(engine, "/error").Code)
	}
	require.EqualValues(t, 1, created_calls.Load())
	require.EqualValues(t, 2, aborted_calls.Load())
	require.EqualValues(t, 2, error_calls.Load())
}
func TestResponseCacheRedisHitExpiryAndFlush(t *testing.T) {
	path := os.Getenv("RESPONSE_CACHE_TEST_REDIS")
	if path == "" {
		t.Skip("set RESPONSE_CACHE_TEST_REDIS to an isolated Redis for integration verification")
	}
	t.Setenv("RESPONSE_CACHE_PATH", path)
	t.Setenv("RESPONSE_CACHE_ENABLED", "true")
	t.Setenv("RESPONSE_CACHE_TTL", "1s")
	store, err := shared.NewResponseCache()
	require.NoError(t, err)
	require.Equal(t, path, store.Path)
	require.Equal(t, time.Second, store.TTL)
	require.IsType(t, &persistence.RedisStore{}, store.Store)
	require.NoError(t, store.Store.Flush())
	t.Cleanup(func() { _ = store.Store.Flush() })
	var calls atomic.Int32
	engine := gin.New()
	engine.GET("/sources", store.Wrap(func(c *gin.Context) {
		c.String(200, "response-%d", calls.Add(1))
	}))
	for i := 0; i < 2; i++ {
		require.Equal(t, "response-1", cacheRequest(engine, "/sources").Body.String())
	}
	key := "expiry"
	require.NoError(t, store.Store.Set(key, "expires", store.TTL))
	var value string
	require.NoError(t, store.Store.Get(key, &value))
	require.Equal(t, "expires", value)
	require.Eventually(t, func() bool { return store.Store.Get(key, &value) == persistence.ErrCacheMiss }, 3*time.Second, 50*time.Millisecond)
	require.NoError(t, store.Store.Flush())
	require.Equal(t, "response-2", cacheRequest(engine, "/sources").Body.String())
}
func TestResponseCacheRedisFailurePreservesComputedResponse(t *testing.T) {
	t.Setenv("RESPONSE_CACHE_PATH", "redis://127.0.0.1:1/0")
	t.Setenv("RESPONSE_CACHE_ENABLED", "true")
	t.Setenv("RESPONSE_CACHE_TTL", "")
	store, err := shared.NewResponseCache()
	require.NoError(t, err)
	engine := gin.New()
	engine.GET("/data", store.Wrap(func(c *gin.Context) { c.String(200, strings.Repeat("ok", 2)) }))
	response := cacheRequest(engine, "/data")
	require.Equal(t, 200, response.Code)
	require.Equal(t, "okok", response.Body.String())
}
