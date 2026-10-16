package shared

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

func init() { cache.RegisterResponseCacheGob() }

// ResponseCache holds the store and configuration shared by every cached route.
type ResponseCache struct {
	Store persistence.CacheStore
	Path  string
	TTL   time.Duration
}

// Wrap delegates response caching to the library using the initialized TTL.
// A nil cache disables caching, including when omitted by router callers.
func (responseCache *ResponseCache) Wrap(handler gin.HandlerFunc) gin.HandlerFunc {
	if responseCache == nil || responseCache.Store == nil {
		return handler
	}
	return cache.CachePage(responseCache.Store, responseCache.TTL, handler)
}

// NewResponseCache selects the library store. A nil cache disables route caching.
func NewResponseCache() (*ResponseCache, error) {
	if raw := strings.TrimSpace(os.Getenv("RESPONSE_CACHE_ENABLED")); raw != "" {
		enabled, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, errors.New("invalid RESPONSE_CACHE_ENABLED")
		}
		if !enabled {
			return nil, nil
		}
	}
	ttl := 10 * time.Minute
	if raw := strings.TrimSpace(os.Getenv("RESPONSE_CACHE_TTL")); raw != "" {
		var err error
		ttl, err = time.ParseDuration(raw)
		if err != nil || ttl <= 0 {
			return nil, errors.New("RESPONSE_CACHE_TTL must be a positive Go duration")
		}
	}
	path := strings.TrimSpace(os.Getenv("RESPONSE_CACHE_PATH"))
	if path == "" {
		return &ResponseCache{Store: persistence.NewInMemoryStore(ttl), Path: path, TTL: ttl}, nil
	}
	address, err := url.Parse(path)
	if err != nil || (address.Scheme != "redis" && address.Scheme != "rediss") || address.Hostname() == "" {
		return nil, errors.New("RESPONSE_CACHE_PATH must be a Redis connection string or blank")
	}
	return &ResponseCache{
		Store: persistence.NewRedisCacheWithURL(path, ttl),
		Path:  path,
		TTL:   ttl,
	}, nil
}
