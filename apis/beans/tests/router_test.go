package gobeansack_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"github.com/soumitsalman/cafecito-api-platform/apis/beans/router"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ROUTE_HEALTH     = "/health"
	ROUTE_DOCS       = "/docs/index.html"
	ROUTE_SEARCH     = "/articles/search"
	ROUTE_LATEST     = "/articles/latest"
	ROUTE_TRENDING   = "/articles/trending"
	ROUTE_HEADLINES  = "/news/top-headlines"
	ROUTE_ARTICLES   = "/articles"
	ROUTE_SOURCES    = "/sources"
	ROUTE_CATEGORIES = "/categories"
	ROUTE_ENTITIES   = "/entities"
	ROUTE_REGIONS    = "/regions"
	ROUTE_SENTIMENTS = "/sentiments"
	ROUTE_TAGS       = "/tags"
	ROUTE_STORIES    = "/stories"
)

func newTestHTTPServer(t *testing.T) *httptest.Server {
	t.Helper()
	db := setupTestDB()
	embedder := setupTestEmbedder()
	gin.SetMode(gin.TestMode)
	engine := router.NewRouter(db, embedder, nil)
	srv := httptest.NewServer(engine)
	t.Cleanup(func() {
		srv.Close()
		_ = embedder.Close()
		db.Close()
	})
	return srv
}

func routerURL(base, path string, params url.Values) string {
	raw := strings.TrimSuffix(base, "/") + path
	if enc := params.Encode(); enc != "" {
		raw += "?" + enc
	}
	return raw
}

func routerGET(t *testing.T, base, path string, params url.Values) (int, []byte) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, routerURL(base, path, params), nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, body
}

func requireStatus(t *testing.T, expected, actual int, body []byte) {
	t.Helper()
	require.Equal(t, expected, actual, "response body: %s", string(body))
}

type pageEnvelope[T any] struct {
	Data       []T            `json:"data"`
	Pagination map[string]any `json:"pagination"`
	Meta       map[string]any `json:"meta"`
}

func parseCollection(t *testing.T, body []byte) pageEnvelope[map[string]any] {
	t.Helper()
	var env pageEnvelope[map[string]any]
	require.NoError(t, json.Unmarshal(body, &env), "response body: %s", string(body))
	return env
}

func parseDetailObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var detail struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(body, &detail), "response body: %s", string(body))
	require.NotNil(t, detail.Data)
	return detail.Data
}

func printResponse(t *testing.T, label string, body []byte) {
	t.Helper()
	if os.Getenv("TEST_PRINT_RESPONSE") == "" && !testing.Verbose() {
		return
	}
	var v any
	if err := json.Unmarshal(bytes.TrimSpace(body), &v); err != nil {
		pp.Println(label, string(body))
		return
	}
	pp.Println(label, v)
}

func assertExpectedAPIError(t *testing.T, body []byte, expected_code string) {
	t.Helper()
	var response struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(body, &response), "response body: %s", string(body))
	assert.Equal(t, expected_code, response.Error.Code)
	assert.NotEmpty(t, response.Error.Message)
}

func assertExpectedPagination(t *testing.T, body []byte, expected_limit int) []map[string]any {
	t.Helper()
	env := parseCollection(t, body)
	require.NotNil(t, env.Data)
	require.Equal(t, float64(expected_limit), env.Pagination["limit"])
	require.Contains(t, env.Pagination, "num_results")
	assert.Equal(t, float64(len(env.Data)), env.Pagination["num_results"])
	require.Contains(t, env.Pagination, "next_cursor")
	assert.Len(t, env.Pagination, 3)
	assert.NotContains(t, env.Pagination, "page")
	assert.NotContains(t, env.Pagination, "offset")
	assert.NotContains(t, env.Pagination, "found")
	assert.NotContains(t, env.Pagination, "returned")
	assert.NotContains(t, env.Pagination, "cursor")
	var raw map[string]any
	require.NoError(t, json.Unmarshal(body, &raw))
	assert.NotContains(t, raw, "success")
	assert.NotContains(t, raw, "status")
	return env.Data
}

func nextCursorFromBody(t *testing.T, body []byte) string {
	t.Helper()
	env := parseCollection(t, body)
	if env.Pagination["next_cursor"] == nil {
		return ""
	}
	cursor, ok := env.Pagination["next_cursor"].(string)
	require.True(t, ok, "pagination.next_cursor must be a string when present")
	return cursor
}

func assertMetaAsOf(t *testing.T, body []byte) {
	t.Helper()
	env := parseCollection(t, body)
	require.Len(t, env.Meta, 1)
	as_of, ok := env.Meta["as_of"].(string)
	require.True(t, ok, "response is missing RFC3339 meta.as_of")
	_, err := time.Parse(time.RFC3339Nano, as_of)
	require.NoError(t, err, "invalid meta.as_of: %q", as_of)
}

func assertDetailMetaAsOf(t *testing.T, body []byte) {
	t.Helper()
	var response struct {
		Meta map[string]any `json:"meta"`
	}
	require.NoError(t, json.Unmarshal(body, &response), "response body: %s", string(body))
	require.Len(t, response.Meta, 1)
	as_of, ok := response.Meta["as_of"].(string)
	require.True(t, ok, "response is missing RFC3339 meta.as_of")
	_, err := time.Parse(time.RFC3339Nano, as_of)
	require.NoError(t, err, "invalid meta.as_of: %q", as_of)
}

func assertStoryArticleMeta(t *testing.T, body []byte, story_id string) {
	t.Helper()
	env := parseCollection(t, body)
	require.Len(t, env.Meta, 2)
	assert.Equal(t, story_id, env.Meta["story_id"])
	as_of, ok := env.Meta["as_of"].(string)
	require.True(t, ok, "response is missing RFC3339 meta.as_of")
	_, err := time.Parse(time.RFC3339Nano, as_of)
	require.NoError(t, err, "invalid meta.as_of: %q", as_of)
}

func assertNoDetailMeta(t *testing.T, body []byte) {
	t.Helper()
	var response map[string]any
	require.NoError(t, json.Unmarshal(body, &response), "response body: %s", string(body))
	assert.NotContains(t, response, "meta")
}

func assertStringArrayField(t *testing.T, item map[string]any, field string) {
	t.Helper()
	require.Contains(t, item, field)
	values, ok := item[field].([]any)
	require.True(t, ok, "%s must be an array, got %T", field, item[field])
	require.NotNil(t, values, "%s must be [] rather than null", field)
}

func assertOptionalStringField(t *testing.T, item map[string]any, field string) {
	t.Helper()
	value, exists := item[field]
	if !exists || value == nil {
		return
	}
	_, ok := value.(string)
	require.True(t, ok, "%s must be a string or null", field)
}

func assertOptionalUUIDField(t *testing.T, item map[string]any, field string) {
	t.Helper()
	value, exists := item[field]
	if !exists || value == nil {
		return
	}
	id, ok := value.(string)
	require.True(t, ok, "%s must be a UUID string or omitted", field)
	_, err := uuid.Parse(id)
	require.NoError(t, err, "invalid %s: %q", field, id)
}

func assertOptionalSourceObject(t *testing.T, parent map[string]any) {
	t.Helper()
	value, exists := parent["source"]
	if !exists || value == nil {
		return
	}
	source, ok := value.(map[string]any)
	require.True(t, ok, "source must be an object or null")
	assertExpectedSource(t, source)
}

func assertExpectedSourceSummary(t *testing.T, source map[string]any) {
	t.Helper()
	assertExpectedSource(t, source)
}

func assertExpectedSource(t *testing.T, source map[string]any) {
	t.Helper()
	require.Contains(t, source, "id")
	require.Contains(t, source, "domain")
	require.Contains(t, source, "url")
	if id, ok := source["id"].(string); ok && id != "" {
		_, err := uuid.Parse(id)
		require.NoError(t, err, "source has invalid id: %q", id)
	}
	assert.NotEmpty(t, source["domain"])
	assert.NotEmpty(t, source["url"])
	assert.NotContains(t, source, "base_url")
	for _, field := range []string{"name", "description", "favicon_url", "rss_feed_url"} {
		assertOptionalStringField(t, source, field)
	}
}

func assertExpectedArticle(t *testing.T, item map[string]any) {
	t.Helper()
	id, ok := item["id"].(string)
	require.True(t, ok, "article is missing string id")
	_, err := uuid.Parse(id)
	require.NoError(t, err, "article has invalid id: %q", id)

	assert.NotEmpty(t, item["url"])
	assert.Contains(t, item, "content_type")
	content_type, ok := item["content_type"].(string)
	require.True(t, ok, "article has invalid content_type")
	assert.Contains(t, []string{
		"blog", "contract", "earnings_report", "enforcement_action", "financial_report",
		"lawsuit", "news", "official_statement", "podcast", "post", "press_release",
		"research_paper", "site", "technical_documentation", "whitepaper",
	}, content_type)
	assert.NotEmpty(t, item["title"])
	assert.Contains(t, item, "summary")
	assert.Contains(t, item, "author")
	assert.Contains(t, item, "image_url")
	assertOptionalUUIDField(t, item, "story_id")
	assertOptionalSourceObject(t, item)

	published_at, ok := item["published_at"].(string)
	require.True(t, ok, "article is missing published_at")
	_, err = time.Parse(time.RFC3339Nano, published_at)
	require.NoError(t, err, "invalid published_at: %q", published_at)

	for _, field := range []string{"categories", "regions", "entities", "sentiments", "tags"} {
		assertStringArrayField(t, item, field)
		for _, raw := range item[field].([]any) {
			value, ok := raw.(string)
			require.True(t, ok, "%s values must be strings", field)
			assert.Equal(t, strings.ToLower(value), value)
			assert.NotContains(t, value, " ")
			assert.NotContains(t, value, "-")
		}
	}
}

func assertExpectedTrend(t *testing.T, item map[string]any) {
	t.Helper()
	trend, ok := item["trend"].(map[string]any)
	require.True(t, ok, "trending article is missing nested trend")
	for _, field := range []string{
		"likes", "comments", "mentions", "audiences", "related", "trend_score",
	} {
		require.Contains(t, trend, field)
		if trend[field] != nil {
			_, ok := trend[field].(float64)
			require.True(t, ok, "trend.%s must be numeric or null", field)
		}
	}
	assert.NotContains(t, trend, "score")
	assert.NotContains(t, trend, "audience")
}

func assertExpectedStory(t *testing.T, story map[string]any) {
	t.Helper()
	story_id, ok := story["id"].(string)
	require.True(t, ok, "story is missing string id")
	_, err := uuid.Parse(story_id)
	require.NoError(t, err, "story has invalid id: %q", story_id)
	assert.NotEmpty(t, story["title"])
	for _, field := range []string{"first_published_at", "last_published_at"} {
		value, ok := story[field].(string)
		require.True(t, ok, "story is missing %s", field)
		_, err := time.Parse(time.RFC3339Nano, value)
		require.NoError(t, err, "invalid story.%s: %q", field, value)
	}
	for _, field := range []string{"article_count", "source_count"} {
		value, ok := story[field].(float64)
		require.True(t, ok, "story is missing numeric %s", field)
		require.GreaterOrEqual(t, value, float64(1))
	}
	for _, field := range []string{"categories", "regions", "entities", "tags"} {
		assertStringArrayField(t, story, field)
	}
	previews, ok := story["top_articles"].([]any)
	require.True(t, ok, "story is missing top_articles")
	require.NotNil(t, previews)
	require.GreaterOrEqual(t, len(previews), 1)
	require.LessOrEqual(t, len(previews), 3)
	for _, raw_preview := range previews {
		preview, ok := raw_preview.(map[string]any)
		require.True(t, ok)
		assertExpectedStoryPreview(t, preview)
	}
	assert.NotContains(t, story, "links")
}

func assertExpectedStoryPreview(t *testing.T, preview map[string]any) {
	t.Helper()
	article_id, ok := preview["id"].(string)
	require.True(t, ok, "story preview is missing string id")
	_, err := uuid.Parse(article_id)
	require.NoError(t, err, "story preview has invalid id: %q", article_id)
	assert.NotEmpty(t, preview["url"])
	assert.NotEmpty(t, preview["title"])
	published_at, ok := preview["published_at"].(string)
	require.True(t, ok, "story preview is missing published_at")
	_, err = time.Parse(time.RFC3339Nano, published_at)
	require.NoError(t, err, "invalid story preview published_at: %q", published_at)
	assertOptionalSourceObject(t, preview)
	for _, field := range []string{"content_type", "summary", "content", "author", "image_url", "story_id", "categories", "regions", "entities", "sentiments", "tags", "trend", "links"} {
		assert.NotContains(t, preview, field)
	}
}

func assertExpectedTag(t *testing.T, item map[string]any) {
	t.Helper()
	value, ok := item["value"].(string)
	require.True(t, ok, "discovery item is missing string value")
	assert.NotEmpty(t, value)
	if item["type"] != nil {
		_, ok = item["type"].(string)
		require.True(t, ok)
	}
}

func skipIfEmbedderUnavailable(t *testing.T, status int, body []byte) {
	t.Helper()
	if status != http.StatusInternalServerError {
		return
	}
	var response struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &response) != nil {
		return
	}
	if response.Error.Code == shared.API_ERROR_EMBEDDING_ERROR {
		t.Skip("embedder unavailable:", response.Error.Message)
	}
}

func addArticleFilters(params url.Values) {
	params.Set("categories", test_categories[0])
	params.Set("limit", "5")
}

func firstArticleID(t *testing.T, base string) string {
	t.Helper()
	params := url.Values{}
	params.Set("limit", "20")
	params.Set("from", testSearchFrom().Format("2006-01-02"))
	status, body := routerGET(t, base, ROUTE_SEARCH, params)
	requireStatus(t, http.StatusOK, status, body)
	items := parseCollection(t, body).Data
	require.NotEmpty(t, items)
	for _, item := range items {
		article_id, _ := item["id"].(string)
		source, _ := item["source"].(map[string]any)
		source_id, _ := source["id"].(string)
		if article_id != "" && source_id != "" {
			return article_id
		}
	}
	id, ok := items[0]["id"].(string)
	require.True(t, ok)
	return id
}

func TestRouterHealth(t *testing.T) {
	srv := newTestHTTPServer(t)
	status, body := routerGET(t, srv.URL, ROUTE_HEALTH, nil)
	printResponse(t, "HEALTH", body)
	requireStatus(t, http.StatusOK, status, body)
	var response map[string]string
	require.NoError(t, json.Unmarshal(body, &response))
	assert.Equal(t, "alive", response["status"])
}

func TestRouterDocs(t *testing.T) {
	srv := newTestHTTPServer(t)
	status, body := routerGET(t, srv.URL, ROUTE_DOCS, nil)
	requireStatus(t, http.StatusOK, status, body)
}

func TestRouterDiscoveryRoutes(t *testing.T) {
	srv := newTestHTTPServer(t)
	for _, path := range []string{ROUTE_CATEGORIES, ROUTE_ENTITIES, ROUTE_REGIONS, ROUTE_SENTIMENTS} {
		t.Run(path, func(t *testing.T) {
			params := url.Values{}
			params.Set("limit", "5")
			status, body := routerGET(t, srv.URL, path, params)
			printResponse(t, path, body)
			requireStatus(t, http.StatusOK, status, body)
			assertMetaAsOf(t, body)
			items := assertExpectedPagination(t, body, 5)
			require.NotEmpty(t, items, path)
			assertExpectedTag(t, items[0])
		})
	}
}

func TestRouterDiscoveryQuery(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("q", "tech")
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_CATEGORIES, params)
	printResponse(t, "CATEGORIES_Q", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	for _, item := range items {
		assertExpectedTag(t, item)
	}
}

func TestRouterInvalidLimit(t *testing.T) {
	srv := newTestHTTPServer(t)
	for _, limit := range []string{"0", "101", "999"} {
		params := url.Values{}
		params.Set("limit", limit)
		status, body := routerGET(t, srv.URL, ROUTE_CATEGORIES, params)
		printResponse(t, "INVALID_LIMIT_"+limit, body)
		requireStatus(t, http.StatusBadRequest, status, body)
		assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
	}
}

func TestRouterGetSources(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("limit", "5")

	status, body := routerGET(t, srv.URL, ROUTE_SOURCES, params)
	printResponse(t, "SOURCES", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	sources := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, sources)
	for _, source := range sources {
		assertExpectedSourceSummary(t, source)
	}

	first_id, ok := sources[0]["id"].(string)
	require.True(t, ok)
	status, body = routerGET(t, srv.URL, ROUTE_SOURCES+"/"+first_id, nil)
	printResponse(t, "SOURCE_DETAIL", body)
	requireStatus(t, http.StatusOK, status, body)
	assertNoDetailMeta(t, body)
	detail := parseDetailObject(t, body)
	assert.Equal(t, first_id, detail["id"])
	assertExpectedSource(t, detail)
}

func TestRouterGetSourcesByQueryAndDomain(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("q", test_source_query)
	params.Set("domains", test_domains[0])
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_SOURCES, params)
	printResponse(t, "SOURCES_FILTERED", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	assertExpectedPagination(t, body, 5)
}

func TestRouterSourcesFilterByIDs(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_SOURCES, params)
	requireStatus(t, http.StatusOK, status, body)
	sources := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, sources)
	sourceID, ok := sources[0]["id"].(string)
	require.True(t, ok)

	params.Set("ids", sourceID)
	status, body = routerGET(t, srv.URL, ROUTE_SOURCES, params)
	printResponse(t, "SOURCES_BY_IDS", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	matched := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, matched)
	for _, source := range matched {
		assertExpectedSourceSummary(t, source)
		assert.Equal(t, sourceID, source["id"])
	}

	params.Set("ids", uuid.New().String())
	status, body = routerGET(t, srv.URL, ROUTE_SOURCES, params)
	printResponse(t, "SOURCES_BY_MISSING_IDS", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	assert.Empty(t, assertExpectedPagination(t, body, 5))
}

func TestRouterSourceNotFound(t *testing.T) {
	srv := newTestHTTPServer(t)
	status, body := routerGET(t, srv.URL, ROUTE_SOURCES+"/"+uuid.New().String(), nil)
	printResponse(t, "SOURCE_NOT_FOUND", body)
	requireStatus(t, http.StatusNotFound, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_NOT_FOUND)
}

func TestRouterSearchArticlesUnfiltered(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
	printResponse(t, "SEARCH_UNFILTERED", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, items)
	for _, item := range items {
		assertExpectedArticle(t, item)
		assert.NotContains(t, item, "content")
	}
}

func TestRouterSearchArticlesByCategories(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	addArticleFilters(params)
	params.Set("from", testSearchFrom().Format("2006-01-02"))
	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
	printResponse(t, "SEARCH_BY_CATEGORIES", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, items)
	for _, item := range items {
		assertExpectedArticle(t, item)
	}
}

func TestRouterSearchArticlesByExactIDsAndURLs(t *testing.T) {
	srv := newTestHTTPServer(t)
	article_id := firstArticleID(t, srv.URL)

	params := url.Values{}
	params.Set("ids", article_id)
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
	printResponse(t, "SEARCH_BY_IDS", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, items)
	assert.Equal(t, article_id, items[0]["id"])

	params = url.Values{}
	params.Set("urls", test_article_urls[0])
	params.Set("limit", "5")
	status, body = routerGET(t, srv.URL, ROUTE_SEARCH, params)
	printResponse(t, "SEARCH_BY_URLS", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items = assertExpectedPagination(t, body, 5)
	for _, item := range items {
		assert.Equal(t, test_article_urls[0], item["url"])
	}
}

func TestRouterSearchArticlesFilters(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("content_type", "news")
	params.Set("authors", test_authors[0])
	params.Set("regions", test_regions[2])
	params.Set("entities", test_entities[0])
	params.Set("sentiments", test_sentiments[0])
	params.Set("tags", test_tags[0])
	params.Set("from", testSearchFrom().Format("2006-01-02"))
	params.Set("to", testSearchTo().Format("2006-01-02"))
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
	printResponse(t, "SEARCH_FILTERS", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	for _, item := range items {
		assertExpectedArticle(t, item)
		assert.Equal(t, "news", item["content_type"])
	}
}

func TestRouterVectorSearchArticles(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("q", TEST_VECTOR_QUERY)
	params.Set("score_threshold", "0.6")
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
	printResponse(t, "VECTOR_SEARCH", body)
	skipIfEmbedderUnavailable(t, status, body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, items)
	for _, item := range items {
		assertExpectedArticle(t, item)
	}
}

func TestRouterScoreThresholdRequiresQ(t *testing.T) {
	srv := newTestHTTPServer(t)
	for _, path := range []string{ROUTE_SEARCH, ROUTE_LATEST, ROUTE_HEADLINES, ROUTE_TRENDING, ROUTE_STORIES} {
		t.Run(path, func(t *testing.T) {
			params := url.Values{}
			params.Set("score_threshold", "0.6")
			params.Set("limit", "5")
			status, body := routerGET(t, srv.URL, path, params)
			printResponse(t, path+"_SCORE_THRESHOLD_WITHOUT_Q", body)
			requireStatus(t, http.StatusBadRequest, status, body)
			assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
		})
	}
}

func TestRouterQueryMayOmitScoreThreshold(t *testing.T) {
	srv := newTestHTTPServer(t)
	for _, path := range []string{ROUTE_SEARCH, ROUTE_LATEST, ROUTE_HEADLINES, ROUTE_TRENDING, ROUTE_STORIES} {
		t.Run(path, func(t *testing.T) {
			params := url.Values{}
			params.Set("q", TEST_VECTOR_QUERY)
			params.Set("limit", "5")
			status, body := routerGET(t, srv.URL, path, params)
			skipIfEmbedderUnavailable(t, status, body)
			requireStatus(t, http.StatusOK, status, body)
		})
	}
}

func TestRouterRejectsRFC3339DateParams(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("from", "2026-02-01T00:00:00Z")
	params.Set("to", "2026-02-02T00:00:00Z")
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
	printResponse(t, "RFC3339_DATES", body)
	requireStatus(t, http.StatusBadRequest, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
}

func TestRouterGetLatestArticles(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	addArticleFilters(params)
	status, body := routerGET(t, srv.URL, ROUTE_LATEST, params)
	printResponse(t, "LATEST", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, items)
	for _, item := range items {
		assertExpectedArticle(t, item)
		assert.NotContains(t, item, "trend")
	}
}

func TestRouterVectorSearchLatest(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("q", TEST_VECTOR_QUERY)
	params.Set("score_threshold", "0.6")
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_LATEST, params)
	printResponse(t, "VECTOR_SEARCH_LATEST", body)
	skipIfEmbedderUnavailable(t, status, body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	for _, item := range items {
		assertExpectedArticle(t, item)
	}
}

func TestRouterGetTrendingArticles(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	addArticleFilters(params)
	status, body := routerGET(t, srv.URL, ROUTE_TRENDING, params)
	printResponse(t, "TRENDING", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, items)
	for _, item := range items {
		assertExpectedArticle(t, item)
		assertExpectedTrend(t, item)
	}
}

func TestRouterGetHeadlines(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	addArticleFilters(params)
	status, body := routerGET(t, srv.URL, ROUTE_HEADLINES, params)
	printResponse(t, "HEADLINES", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, items)
	for _, item := range items {
		assertExpectedArticle(t, item)
		assert.NotContains(t, item, "trend")
	}
}

func TestRouterLatestAndTrendingRejectExactIdentityFilters(t *testing.T) {
	srv := newTestHTTPServer(t)
	article_id := firstArticleID(t, srv.URL)
	for _, path := range []string{ROUTE_LATEST, ROUTE_TRENDING, ROUTE_HEADLINES} {
		t.Run(path, func(t *testing.T) {
			params := url.Values{}
			params.Set("ids", article_id)
			params.Set("urls", test_article_urls[0])
			params.Set("limit", "5")
			status, body := routerGET(t, srv.URL, path, params)
			printResponse(t, path+"_IDENTITY_FILTERS", body)
			requireStatus(t, http.StatusBadRequest, status, body)
			assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
		})
	}
}

func TestRouterFeedsRejectDisallowedDateBounds(t *testing.T) {
	srv := newTestHTTPServer(t)
	for _, path := range []string{ROUTE_LATEST, ROUTE_HEADLINES, ROUTE_TRENDING} {
		t.Run(path, func(t *testing.T) {
			params := url.Values{}
			params.Set("from", testSearchFrom().Format("2006-01-02"))
			params.Set("to", testSearchTo().Format("2006-01-02"))
			params.Set("limit", "5")
			status, body := routerGET(t, srv.URL, path, params)
			printResponse(t, path+"_DATE_BOUNDS", body)
			requireStatus(t, http.StatusBadRequest, status, body)
			assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
		})
	}
}

func TestRouterTopHeadlinesRejectsContentType(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("content_type", "news")
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_HEADLINES, params)
	printResponse(t, "HEADLINES_CONTENT_TYPE", body)
	requireStatus(t, http.StatusBadRequest, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
}

func TestRouterSearchRejectsUnsupportedParameters(t *testing.T) {
	srv := newTestHTTPServer(t)
	for key, value := range map[string]string{
		"content_type": "post",
		"offset":       "5",
		"page":         "2",
		"sort":         "published_at",
	} {
		t.Run(key, func(t *testing.T) {
			params := url.Values{}
			params.Set(key, value)
			params.Set("limit", "5")
			status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
			printResponse(t, "SEARCH_UNSUPPORTED_"+key, body)
			requireStatus(t, http.StatusBadRequest, status, body)
			assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
		})
	}
}

func TestRouterGetArticle(t *testing.T) {
	srv := newTestHTTPServer(t)
	article_id := firstArticleID(t, srv.URL)

	status, body := routerGET(t, srv.URL, ROUTE_ARTICLES+"/"+article_id, nil)
	printResponse(t, "ARTICLE_DETAIL", body)
	requireStatus(t, http.StatusOK, status, body)
	assertNoDetailMeta(t, body)
	detail := parseDetailObject(t, body)
	assert.Equal(t, article_id, detail["id"])
	assertExpectedArticle(t, detail)
	assert.NotContains(t, detail, "content")
	links, ok := detail["links"].(map[string]any)
	require.True(t, ok, "article detail is missing links")
	assert.Equal(t, "/articles/"+article_id+"/similar", links["similar"])
	assert.Equal(t, "/articles/"+article_id+"/mentions", links["mentions"])
}

func TestRouterGetArticleFullContent(t *testing.T) {
	srv := newTestHTTPServer(t)
	article_id := firstArticleID(t, srv.URL)
	params := url.Values{}
	params.Set("full_content", "true")
	status, body := routerGET(t, srv.URL, ROUTE_ARTICLES+"/"+article_id, params)
	printResponse(t, "ARTICLE_FULL_CONTENT", body)
	requireStatus(t, http.StatusOK, status, body)
	assertNoDetailMeta(t, body)
	detail := parseDetailObject(t, body)
	require.Contains(t, detail, "content")
}

func TestRouterArticleNotFound(t *testing.T) {
	srv := newTestHTTPServer(t)
	status, body := routerGET(t, srv.URL, ROUTE_ARTICLES+"/"+uuid.New().String(), nil)
	printResponse(t, "ARTICLE_NOT_FOUND", body)
	requireStatus(t, http.StatusNotFound, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_NOT_FOUND)
}

func TestRouterArticleInvalidID(t *testing.T) {
	srv := newTestHTTPServer(t)
	status, body := routerGET(t, srv.URL, ROUTE_ARTICLES+"/not-a-uuid", nil)
	printResponse(t, "ARTICLE_INVALID_ID", body)
	requireStatus(t, http.StatusBadRequest, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
}

func TestRouterSimilarArticles(t *testing.T) {
	srv := newTestHTTPServer(t)
	article_id := firstArticleID(t, srv.URL)
	params := url.Values{}
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_ARTICLES+"/"+article_id+"/similar", params)
	printResponse(t, "SIMILAR", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	for _, item := range items {
		assertExpectedArticle(t, item)
	}
}

func TestRouterArticleMentions(t *testing.T) {
	srv := newTestHTTPServer(t)
	article_id := firstArticleID(t, srv.URL)
	params := url.Values{}
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_ARTICLES+"/"+article_id+"/mentions", params)
	printResponse(t, "MENTIONS", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	for _, item := range items {
		require.Contains(t, item, "url")
		require.Contains(t, item, "platform")
		require.Contains(t, item, "forum")
		require.Contains(t, item, "observed_at")
		engagement, ok := item["engagement"].(map[string]any)
		require.True(t, ok)
		require.Contains(t, engagement, "likes")
		require.Contains(t, engagement, "comments")
		require.Contains(t, engagement, "audience")
	}
}

func TestRouterMissingArticleSubresources(t *testing.T) {
	srv := newTestHTTPServer(t)
	missing := uuid.New().String()
	for _, suffix := range []string{"/similar", "/mentions"} {
		status, body := routerGET(t, srv.URL, ROUTE_ARTICLES+"/"+missing+suffix, url.Values{"limit": {"5"}})
		printResponse(t, "MISSING"+suffix, body)
		requireStatus(t, http.StatusNotFound, status, body)
		assertExpectedAPIError(t, body, shared.API_ERROR_NOT_FOUND)
	}
}

func TestRouterStories(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("limit", "5")
	params.Set("from", testSearchFrom().Format("2006-01-02"))
	status, body := routerGET(t, srv.URL, ROUTE_STORIES, params)
	printResponse(t, "STORIES", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, items)
	story := items[0]
	assertExpectedStory(t, story)

	story_id, ok := story["id"].(string)
	require.True(t, ok)
	status, body = routerGET(t, srv.URL, ROUTE_STORIES+"/"+story_id, nil)
	printResponse(t, "STORY_DETAIL", body)
	requireStatus(t, http.StatusOK, status, body)
	assertNoDetailMeta(t, body)
	detail := parseDetailObject(t, body)
	assert.Equal(t, story_id, detail["id"])
	detail_story := make(map[string]any, len(detail)-1)
	for key, value := range detail {
		if key != "links" {
			detail_story[key] = value
		}
	}
	assertExpectedStory(t, detail_story)
	links, ok := detail["links"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "/stories/"+story_id+"/articles", links["articles"])

	status, body = routerGET(t, srv.URL, ROUTE_STORIES+"/"+story_id+"/articles", url.Values{"limit": {"5"}})
	printResponse(t, "STORY_ARTICLES", body)
	requireStatus(t, http.StatusOK, status, body)
	assertStoryArticleMeta(t, body, story_id)
	members := assertExpectedPagination(t, body, 5)
	for _, item := range members {
		assertExpectedArticle(t, item)
		assert.Equal(t, story_id, item["story_id"])
		assert.NotContains(t, item, "trend")
	}
}

func TestRouterStoryNotFound(t *testing.T) {
	srv := newTestHTTPServer(t)
	missing := uuid.New().String()
	status, body := routerGET(t, srv.URL, ROUTE_STORIES+"/"+missing, nil)
	printResponse(t, "STORY_NOT_FOUND", body)
	requireStatus(t, http.StatusNotFound, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_NOT_FOUND)
}

func TestRouterBackendAPIKey(t *testing.T) {
	db := setupTestDB()
	embedder := setupTestEmbedder()
	gin.SetMode(gin.TestMode)
	engine := router.NewRouter(db, embedder, map[string]string{"X-API-KEY": "test-token"})
	srv := httptest.NewServer(engine)
	t.Cleanup(func() {
		srv.Close()
		_ = embedder.Close()
		db.Close()
	})

	status, body := routerGET(t, srv.URL, ROUTE_HEALTH, nil)
	requireStatus(t, http.StatusOK, status, body)

	status, body = routerGET(t, srv.URL, ROUTE_SEARCH, url.Values{"limit": {"1"}})
	requireStatus(t, http.StatusUnauthorized, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_UNAUTHORIZED)

	req, err := http.NewRequest(http.MethodGet, routerURL(srv.URL, ROUTE_SEARCH, url.Values{"limit": {"1"}}), nil)
	require.NoError(t, err)
	req.Header.Set("X-API-KEY", "test-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	requireStatus(t, http.StatusOK, resp.StatusCode, body)
	assertMetaAsOf(t, body)
}

func TestRouterPostContentTypeIsResponseOnly(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("ids", fixturePostID.String())
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
	requireStatus(t, http.StatusOK, status, body)
	items := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, items)
	assert.Equal(t, "post", items[0]["content_type"])
}

func TestRouterCursorPagination(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("limit", "5")
	params.Set("from", time.Now().UTC().AddDate(0, 0, -30).Format("2006-01-02"))
	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
	requireStatus(t, http.StatusOK, status, body)
	first := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, first)

	cursor := nextCursorFromBody(t, body)
	if cursor == "" {
		t.Skip("not enough articles for a second page")
	}

	params.Set("cursor", cursor)
	status, body = routerGET(t, srv.URL, ROUTE_SEARCH, params)
	requireStatus(t, http.StatusOK, status, body)
	second := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, second)
	assert.NotEqual(t, first[0]["id"], second[0]["id"])
}

func TestRouterInvalidCursor(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("cursor", "not-a-valid-cursor")
	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
	printResponse(t, "INVALID_CURSOR", body)
	requireStatus(t, http.StatusBadRequest, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
}

func TestRouterDefaultPagination(t *testing.T) {
	srv := newTestHTTPServer(t)
	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, nil)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 20)
	require.NotEmpty(t, items)
	assert.LessOrEqual(t, len(items), 20)
}

func TestRouterEmptyCollection(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("from", "2000-01-01")
	params.Set("to", "2000-01-02")
	params.Set("limit", "5")
	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
	printResponse(t, "EMPTY_COLLECTION", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	items := assertExpectedPagination(t, body, 5)
	assert.Empty(t, items)
	assert.Equal(t, "", nextCursorFromBody(t, body))
}

func TestRouterFullContentProjection(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("limit", "5")
	params.Set("from", testSearchFrom().Format("2006-01-02"))

	status, body := routerGET(t, srv.URL, ROUTE_SEARCH, params)
	printResponse(t, "WITHOUT_FULL_CONTENT", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	without_content := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, without_content)
	assert.NotContains(t, without_content[0], "content")

	params.Set("full_content", "true")
	status, body = routerGET(t, srv.URL, ROUTE_SEARCH, params)
	printResponse(t, "FULL_CONTENT", body)
	requireStatus(t, http.StatusOK, status, body)
	assertMetaAsOf(t, body)
	with_content := assertExpectedPagination(t, body, 5)
	require.NotEmpty(t, with_content)
	assert.Equal(t, without_content[0]["id"], with_content[0]["id"])
	if content, present := with_content[0]["content"]; present && content != nil {
		_, ok := content.(string)
		require.True(t, ok, "content must be a string or null")
	}
}
