package espressoapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/router"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared"
	datautils "github.com/soumitsalman/data-utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const DEFAULT_STRESS_BASE_URL = "http://localhost:8080"

const (
	MIN_CONCURRENCY = 100
	MAX_CONCURRENCY = 10000
	HTTP_TIMEOUT    = 10 * time.Minute

	ROUTE_DOCS        = "/docs/index.html"
	ROUTE_HEALTH      = "/health"
	ROUTE_TAGS        = "/tags"
	ROUTE_ENTITIES    = "/entities"
	ROUTE_REGIONS     = "/regions"
	ROUTE_EVENT_TYPES = "/event-types"
	ROUTE_EVENTS      = "/events"
	ROUTE_SIGNALS     = "/signals"
	ROUTE_SOURCES     = "/sources"
)

// stressEndpoint describes one API endpoint and its optional query params (router/routes.go).
type stressEndpoint struct {
	path         string
	accepts_q    bool
	accepts_tags bool
	accepts_from bool
}

var stress_endpoints = []stressEndpoint{
	{path: ROUTE_EVENTS, accepts_q: true, accepts_tags: true, accepts_from: true},
	{path: ROUTE_SIGNALS, accepts_q: true, accepts_tags: true, accepts_from: true},
	{path: ROUTE_TAGS},
	{path: ROUTE_SOURCES},
}

var sample_queries = []string{
	"artificial intelligence",
	"machine learning",
	"cloud computing",
	"cybersecurity breaches",
	"open source software",
	"startup funding",
	"climate change policy",
	"quantum computing",
	"electric vehicles",
	"blockchain technology",
}

var sample_tags = []string{
	"public_policy",
	"market_trends",
	"criminal_investigation",
	"OpenAI",
	"Google",
	"US",
	"Europe",
}

// --- in-process router integration tests (mirror tests/db_test.go) ---

func newTestHTTPServer(t *testing.T) *httptest.Server {
	t.Helper()
	db := setupTestDB()
	embedder := setupTestEmbedder()
	gin.SetMode(gin.TestMode)
	engine := router.NewRouter(db, embedder, nil)
	srv := httptest.NewServer(engine)
	t.Cleanup(func() {
		srv.Close()
		embedder.Close()
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

func routerGET(t *testing.T, base, path string, params url.Values, api_key string) (int, []byte) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, routerURL(base, path, params), nil)
	require.NoError(t, err)
	if api_key != "" {
		req.Header.Set("X-API-KEY", api_key)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, body
}

func routerPOST(t *testing.T, base, path string, payload any, api_key string) (int, []byte) {
	t.Helper()
	body, err := json.Marshal(payload)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, strings.TrimSuffix(base, "/")+path, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if api_key != "" {
		req.Header.Set("X-API-KEY", api_key)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	response_body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, response_body
}

func requireStatus(t *testing.T, expected int, actual int, body []byte) {
	t.Helper()
	require.Equal(t, expected, actual, "response body: %s", string(body))
}

// pageEnvelope mirrors router.PageResponse for parsing list responses.
type pageEnvelope[T any] struct {
	Data       []T            `json:"data"`
	Pagination map[string]any `json:"pagination"`
	Meta       map[string]any `json:"meta"`
}

func parseDigestArray(t *testing.T, body []byte) []map[string]any {
	t.Helper()
	var env pageEnvelope[map[string]any]
	require.NoError(t, json.Unmarshal(body, &env))
	return env.Data
}

func parseStringArray(t *testing.T, body []byte) []string {
	t.Helper()
	var env pageEnvelope[string]
	require.NoError(t, json.Unmarshal(body, &env))
	return env.Data
}

func parseSourceArray(t *testing.T, body []byte) []map[string]any {
	t.Helper()
	var env pageEnvelope[map[string]any]
	require.NoError(t, json.Unmarshal(body, &env))
	return env.Data
}

func parseDetailObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var detail struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(body, &detail))
	return detail.Data
}

func assertExpectedPagination(t *testing.T, body []byte, expected_limit int, expected_cursor string) []map[string]any {
	t.Helper()
	var env pageEnvelope[map[string]any]
	require.NoError(t, json.Unmarshal(body, &env))
	require.Equal(t, float64(expected_limit), env.Pagination["limit"])
	require.Contains(t, env.Pagination, "cursor")
	require.Contains(t, env.Pagination, "next_cursor")
	if expected_cursor == "" {
		assert.Nil(t, env.Pagination["cursor"])
	} else {
		assert.Equal(t, expected_cursor, env.Pagination["cursor"])
	}
	as_of, ok := env.Meta["as_of"].(string)
	require.True(t, ok, "pagination response is missing RFC3339 meta.as_of")
	_, err := time.Parse(time.RFC3339Nano, as_of)
	require.NoError(t, err, "invalid meta.as_of: %q", as_of)
	return env.Data
}

func assertExpectedSip(t *testing.T, item map[string]any, expected_kind string) {
	t.Helper()
	id, ok := item["id"].(string)
	require.True(t, ok, "response item is missing string id")
	_, err := uuid.Parse(id)
	require.NoError(t, err, "response item has invalid id: %q", id)
	assert.Equal(t, expected_kind, item["kind"])

	created_at, ok := item["created_at"].(string)
	require.True(t, ok, "response item is missing string created_at")
	_, err = time.Parse(time.RFC3339Nano, created_at)
	require.NoError(t, err, "response item has invalid created_at: %q", created_at)

	// require.Contains(t, item, "source_id") // TODO: enable later
	// require.Contains(t, item, "tags") // TODO: enable later
	if briefing, ok := item["briefing"].(string); ok && briefing != "" {
		assert.Equal(t, briefing, item["summary"])
	}
	assert.NotContains(t, item, "digest")
	assert.NotContains(t, item, "representation")
	assert.NotContains(t, item, "object")

	if source_id, ok := item["source_id"].(string); ok {
		_, err = uuid.Parse(source_id)
		require.NoError(t, err, "response item has invalid source_id: %q", source_id)
		if source, ok := item["source"].(map[string]any); ok {
			assert.Equal(t, source_id, source["id"])
		}
	}
}

func assertExpectedAPIError(t *testing.T, body []byte, expected_code string) {
	t.Helper()
	var response struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(body, &response))
	assert.Equal(t, expected_code, response.Error.Code)
	assert.NotEmpty(t, response.Error.Message)
}

func nextCursorFromBody(t *testing.T, body []byte) string {
	t.Helper()
	var env struct {
		Pagination struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"pagination"`
	}
	require.NoError(t, json.Unmarshal(body, &env))
	if env.Pagination.NextCursor == nil {
		return ""
	}
	return *env.Pagination.NextCursor
}

func cloneURLValues(params url.Values) url.Values {
	params_copy := make(url.Values, len(params))
	for key, values := range params {
		params_copy[key] = append([]string(nil), values...)
	}
	return params_copy
}

func requirePaginatedCollection(t *testing.T, base, path string, params url.Values) []map[string]any {
	t.Helper()
	params_copy := cloneURLValues(params)
	params_copy.Del("cursor")
	params_copy.Set("limit", "1")

	status, first_body := routerGET(t, base, path, params_copy, "")
	requireStatus(t, http.StatusOK, status, first_body)
	first := assertExpectedPagination(t, first_body, 1, "")
	require.NotEmpty(t, first, "first page for %s", path)

	cursor := nextCursorFromBody(t, first_body)
	require.NotEmpty(t, cursor, "expected next_cursor for %s", path)

	params_copy.Set("cursor", cursor)
	status, second_body := routerGET(t, base, path, params_copy, "")
	requireStatus(t, http.StatusOK, status, second_body)
	second := assertExpectedPagination(t, second_body, 1, cursor)
	require.NotEmpty(t, second, "second page for %s", path)
	assert.NotEqual(t, first[0], second[0], "cursor did not advance %s", path)

	return first
}

func collectionIDs(t *testing.T, base, path string) []string {
	t.Helper()
	status, body := routerGET(t, base, path, url.Values{"limit": {"100"}}, "")
	requireStatus(t, http.StatusOK, status, body)
	items := parseDigestArray(t, body)
	ids := make([]string, 0, len(items))
	for _, item := range items {
		id, ok := item["id"].(string)
		require.True(t, ok, "missing id in %s response item", path)
		ids = append(ids, id)
	}
	return ids
}

func findPaginatedRelationPath(t *testing.T, base, collection_path, suffix string) string {
	t.Helper()
	for _, id := range collectionIDs(t, base, collection_path) {
		path := collection_path + "/" + id + suffix
		status, body := routerGET(t, base, path, url.Values{"limit": {"1"}}, "")
		if status == http.StatusOK && nextCursorFromBody(t, body) != "" {
			return path
		}
	}
	require.Failf(t, "missing paginated relation fixture", "no %s%s collection has a next cursor", collection_path, suffix)
	return ""
}

func TestRouterHealth(t *testing.T) {
	srv := newTestHTTPServer(t)
	status, body := routerGET(t, srv.URL, ROUTE_HEALTH, nil, "")
	requireStatus(t, http.StatusOK, status, body)
	var response map[string]string
	require.NoError(t, json.Unmarshal(body, &response))
	assert.Equal(t, "alive", response["status"])
}

func TestRouterDocs(t *testing.T) {
	srv := newTestHTTPServer(t)
	status, body := routerGET(t, srv.URL, ROUTE_DOCS, nil, "")
	requireStatus(t, http.StatusOK, status, body)
}

func TestRouterPaginatedCollectionRoutes(t *testing.T) {
	srv := newTestHTTPServer(t)

	for _, route := range []string{
		ROUTE_TAGS,
		ROUTE_ENTITIES,
		ROUTE_REGIONS,
		ROUTE_EVENT_TYPES,
		ROUTE_EVENTS,
		ROUTE_SIGNALS,
		ROUTE_SOURCES,
	} {
		t.Run(route, func(t *testing.T) {
			requirePaginatedCollection(t, srv.URL, route, url.Values{})
		})
	}

	for _, relation := range []struct {
		collection_path string
		suffix          string
	}{
		{collection_path: ROUTE_EVENTS, suffix: "/evidence"},
		{collection_path: ROUTE_EVENTS, suffix: "/signals"},
		{collection_path: ROUTE_SIGNALS, suffix: "/events"},
	} {
		t.Run(relation.collection_path+"/:id"+relation.suffix, func(t *testing.T) {
			path := findPaginatedRelationPath(t, srv.URL, relation.collection_path, relation.suffix)
			requirePaginatedCollection(t, srv.URL, path, url.Values{})
		})
	}
}

func TestRouterGetTags(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("q", "academic")
	params.Set("resource", "event")
	params.Set("limit", "5")

	status, body := routerGET(t, srv.URL, ROUTE_TAGS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	tags := assertExpectedPagination(t, body, 5, "")
	assert.NotEmpty(t, tags)
	assert.NotEmpty(t, tags[0]["value"])
	pp.Println("TAGS", tags)
}

func TestRouterScalarSearchEvents(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	for _, tag := range test_scalar_tags {
		params.Add("tags", tag)
	}
	params.Set("from", testSearchFrom().Format("2006-01-02"))

	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	events := assertExpectedPagination(t, body, 20, "")
	assert.Greater(t, len(events), 0)
	for _, event := range events {
		assertExpectedSip(t, event, "event")
	}
	pp.Println("EVENTS", events)
}

func TestRouterSearchEventsByDigestTags(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	for _, tag := range []string{"us", "japan", "china"} {
		params.Add("regions", tag)
	}
	params.Set("from", testSearchFrom().Format("2006-01-02"))

	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	events := assertExpectedPagination(t, body, 20, "")
	assert.Greater(t, len(events), 0)
	for _, event := range events {
		assertExpectedSip(t, event, "event")
	}
	pp.Println("EVENTS", events)
}

func TestRouterVectorSearchEvents(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("q", TEST_VECTOR_QUERY)
	params.Set("score_threshold", "0.6")
	params.Set("limit", "5")
	params.Set("from", testSearchFrom().Format("2006-01-02"))

	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	events := assertExpectedPagination(t, body, 5, "")
	assert.Greater(t, len(events), 0)
	for _, event := range events {
		assertExpectedSip(t, event, "event")
	}
	pp.Println("EVENTS", events)
}

func TestRouterVectorSearchSignals(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("q", TEST_VECTOR_QUERY)
	params.Set("score_threshold", "0.7")
	params.Set("limit", "5")
	params.Set("to", testSearchFrom().Format("2006-01-02"))

	status, body := routerGET(t, srv.URL, ROUTE_SIGNALS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	signals := assertExpectedPagination(t, body, 5, "")
	assert.Greater(t, len(signals), 0)
	for _, signal := range signals {
		assertExpectedSip(t, signal, "signal")
	}
	pp.Println("SIGNALS", signals)
}

func TestRouterGetSources(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("limit", "5")

	status, body := routerGET(t, srv.URL, ROUTE_SOURCES, params, "")
	requireStatus(t, http.StatusOK, status, body)
	sources := assertExpectedPagination(t, body, 5, "")
	assert.NotEmpty(t, sources)
	for _, source := range sources {
		require.Contains(t, source, "id")
		assert.Contains(t, source, "url")
		assert.Contains(t, source, "domain")
		assert.Contains(t, source, "name")
		assert.Contains(t, source, "description")
		assert.Contains(t, source, "favicon_url")
		assert.Contains(t, source, "rss_feed_url")
		assert.NotContains(t, source, "base_url")
	}
	pp.Println("SOURCES", sources)

	first_id, ok := sources[0]["id"].(string)
	require.True(t, ok)
	status, body = routerGET(t, srv.URL, ROUTE_SOURCES+"/"+first_id, nil, "")
	requireStatus(t, http.StatusOK, status, body)
	src := parseDetailObject(t, body)
	assert.Equal(t, first_id, src["id"])
	require.Contains(t, src, "url")
	require.Contains(t, src, "domain")
	require.Contains(t, src, "name")
	assert.NotContains(t, src, "base_url")
}

func TestRouterSignalDetailsAndEvents(t *testing.T) {
	srv := newTestHTTPServer(t)

	params := url.Values{}
	params.Set("limit", "1")
	status, body := routerGET(t, srv.URL, ROUTE_SIGNALS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	signals := assertExpectedPagination(t, body, 1, "")
	require.NotEmpty(t, signals)
	for _, signal := range signals {
		assertExpectedSip(t, signal, "signal")
	}
	signal_id, ok := signals[0]["id"].(string)
	require.True(t, ok)

	status, body = routerGET(t, srv.URL, ROUTE_SIGNALS+"/"+signal_id, nil, "")
	requireStatus(t, http.StatusOK, status, body)
	detail := parseDetailObject(t, body)
	assert.Equal(t, signal_id, detail["id"])
	assertExpectedSip(t, detail, "signal")
	// require.Contains(t, detail, "source_id") // TODO: enable later
	assert.Contains(t, detail, "links")
	counts, ok := detail["counts"].(map[string]any)
	require.True(t, ok)
	require.Contains(t, counts, "events")
	pp.Println("DETAIL", detail)

	status, body = routerGET(t, srv.URL, ROUTE_SIGNALS+"/"+signal_id+"/events", nil, "")
	requireStatus(t, http.StatusOK, status, body)
	events := assertExpectedPagination(t, body, 20, "")
	assert.NotEmpty(t, events)
	for _, event := range events {
		assertExpectedSip(t, event, "event")
	}
	pp.Println("EVENTS", events)
}

func TestRouterEventDetailEvidenceAndSignals(t *testing.T) {
	srv := newTestHTTPServer(t)

	// get a root signal
	params := url.Values{}
	params.Set("limit", "1")
	status, body := routerGET(t, srv.URL, ROUTE_SIGNALS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	root_signals := assertExpectedPagination(t, body, 1, "")
	require.NotEmpty(t, root_signals)
	for _, signal := range root_signals {
		assertExpectedSip(t, signal, "signal")
	}
	signal_id, ok := root_signals[0]["id"].(string)
	require.True(t, ok)

	status, body = routerGET(t, srv.URL, ROUTE_SIGNALS+"/"+signal_id+"/events?limit=1", nil, "")
	requireStatus(t, http.StatusOK, status, body)
	events := assertExpectedPagination(t, body, 1, "")
	assert.NotEmpty(t, events)
	for _, event := range events {
		assertExpectedSip(t, event, "event")
	}
	event_id, ok := events[0]["id"].(string)
	require.True(t, ok)

	status, body = routerGET(t, srv.URL, ROUTE_EVENTS+"/"+event_id, nil, "")
	requireStatus(t, http.StatusOK, status, body)
	detail := parseDetailObject(t, body)
	assert.Equal(t, event_id, detail["id"])
	assertExpectedSip(t, detail, "event")
	// require.Contains(t, detail, "source_id") // TODO: enable later
	assert.Contains(t, detail, "links")
	counts, ok := detail["counts"].(map[string]any)
	require.True(t, ok)
	require.Contains(t, counts, "evidence")
	require.Contains(t, counts, "signals")
	pp.Println("DETAIL", detail)

	status, body = routerGET(t, srv.URL, ROUTE_EVENTS+"/"+event_id+"/evidence", nil, "")
	requireStatus(t, http.StatusOK, status, body)
	evidence := assertExpectedPagination(t, body, 20, "")
	pp.Println("EVIDENCE", evidence)
	assert.NotEmpty(t, evidence)
	for _, event := range evidence {
		assertExpectedSip(t, event, "event")
	}

	status, body = routerGET(t, srv.URL, ROUTE_EVENTS+"/"+event_id+"/signals", nil, "")
	requireStatus(t, http.StatusOK, status, body)
	signals := assertExpectedPagination(t, body, 20, "")
	assert.NotEmpty(t, signals)
	for _, signal := range signals {
		assertExpectedSip(t, signal, "signal")
	}
	signal_ids := datautils.Transform(signals, func(signal *map[string]any) string {
		return (*signal)["id"].(string)
	})
	assert.Contains(t, signal_ids, signal_id)
	pp.Println("SIGNALS", signals)
}

func TestRouterEventDetailNotFound(t *testing.T) {
	srv := newTestHTTPServer(t)
	bogus := uuid.New().String()
	status, body := routerGET(t, srv.URL, ROUTE_EVENTS+"/"+bogus, nil, "")
	requireStatus(t, http.StatusNotFound, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_NOT_FOUND)
}

func TestRouterCursorPagination(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("limit", "2")
	params.Set("from", time.Now().UTC().AddDate(0, 0, -30).Format("2006-01-02"))

	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	first := assertExpectedPagination(t, body, 2, "")
	require.NotEmpty(t, first)

	cursor := nextCursorFromBody(t, body)
	if cursor == "" {
		t.Skip("not enough events for a second page")
	}

	params.Set("cursor", cursor)
	status, body = routerGET(t, srv.URL, ROUTE_EVENTS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	second := assertExpectedPagination(t, body, 2, cursor)
	require.NotEmpty(t, second)
	assert.NotEqual(t, first[0]["id"], second[0]["id"])
}

func TestRouterInvalidCursor(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("cursor", "not-a-valid-cursor")

	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, params, "")
	requireStatus(t, http.StatusBadRequest, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
}

func TestRouterExpectedDefaultPagination(t *testing.T) {
	srv := newTestHTTPServer(t)
	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, nil, "")
	requireStatus(t, http.StatusOK, status, body)
	events := assertExpectedPagination(t, body, 20, "")
	require.NotEmpty(t, events)
	for _, event := range events {
		assertExpectedSip(t, event, "event")
	}
}

func TestRouterEventCategoriesAlias(t *testing.T) {
	srv := newTestHTTPServer(t)
	status, body := routerGET(t, srv.URL, ROUTE_EVENT_TYPES, url.Values{"limit": {"1"}}, "")
	requireStatus(t, http.StatusOK, status, body)
	event_types := assertExpectedPagination(t, body, 1, "")
	require.NotEmpty(t, event_types)
	event_type, ok := event_types[0]["value"].(string)
	require.True(t, ok)

	by_event_type_params := url.Values{"event_types": {event_type}, "limit": {"1"}}
	status, body = routerGET(t, srv.URL, ROUTE_EVENTS, by_event_type_params, "")
	requireStatus(t, http.StatusOK, status, body)
	by_event_type := assertExpectedPagination(t, body, 1, "")
	require.Len(t, by_event_type, 1)
	assertExpectedSip(t, by_event_type[0], "event")
	assert.Equal(t, event_type, by_event_type[0]["event_type"])

	by_category_params := url.Values{"categories": {event_type}, "limit": {"1"}}
	status, body = routerGET(t, srv.URL, ROUTE_EVENTS, by_category_params, "")
	requireStatus(t, http.StatusOK, status, body)
	by_category := assertExpectedPagination(t, body, 1, "")
	require.Len(t, by_category, 1)
	assertExpectedSip(t, by_category[0], "event")
	assert.Equal(t, by_event_type[0]["id"], by_category[0]["id"])
}

func TestRouterAcceptsRFC3339TimeBounds(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{
		"from":  {"2000-01-01T00:00:00Z"},
		"to":    {"2100-01-01T00:00:00Z"},
		"limit": {"1"},
	}
	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	events := assertExpectedPagination(t, body, 1, "")
	require.NotEmpty(t, events)
	assertExpectedSip(t, events[0], "event")
}

func TestRouterExpectedAggregateRoutes(t *testing.T) {
	srv := newTestHTTPServer(t)
	count_params := url.Values{
		"from": {"2026-08-01"},
		"to":   {"2026-08-10"},
	}
	status, body := routerGET(t, srv.URL, ROUTE_EVENTS+"/count", count_params, "")
	requireStatus(t, http.StatusOK, status, body)
	var count_response map[string]any
	require.NoError(t, json.Unmarshal(body, &count_response))
	count_data, ok := count_response["data"].(map[string]any)
	require.True(t, ok)
	count, ok := count_data["count"].(float64)
	require.True(t, ok)
	assert.GreaterOrEqual(t, count, float64(0))
	assert.Contains(t, count_data, "event_types")
	assert.Contains(t, count_data, "impact_levels")
	count_meta, ok := count_response["meta"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "created_at", count_meta["time_field"])
	count_as_of, ok := count_meta["as_of"].(string)
	require.True(t, ok)
	_, err := time.Parse(time.RFC3339Nano, count_as_of)
	require.NoError(t, err)

	summary_params := url.Values{
		"from":     {"2026-08-01"},
		"to":       {"2026-08-10"},
		"group_by": {"event_type"},
	}
	status, body = routerGET(t, srv.URL, ROUTE_EVENTS+"/summary", summary_params, "")
	requireStatus(t, http.StatusOK, status, body)
	var summary_response map[string]any
	require.NoError(t, json.Unmarshal(body, &summary_response))
	assert.Equal(t, "event_type", summary_response["group_by"])
	_, ok = summary_response["data"].([]any)
	require.True(t, ok)
	summary_meta, ok := summary_response["meta"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "event", summary_meta["counted_resource"])
	assert.Equal(t, "created_at", summary_meta["time_field"])
	summary_as_of, ok := summary_meta["as_of"].(string)
	require.True(t, ok)
	_, err = time.Parse(time.RFC3339Nano, summary_as_of)
	require.NoError(t, err)
}

func TestRouterExpectedEventSearchPost(t *testing.T) {
	srv := newTestHTTPServer(t)
	status, body := routerPOST(t, srv.URL, ROUTE_EVENTS+"/search", map[string]any{
		"q":     TEST_VECTOR_QUERY,
		"limit": 1,
	}, "")
	requireStatus(t, http.StatusOK, status, body)
	events := assertExpectedPagination(t, body, 1, "")
	require.NotEmpty(t, events)
	assertExpectedSip(t, events[0], "event")
}

// --- stress tests against a live server ---

type stressResult struct {
	endpoint    string
	status_code int
	latency     time.Duration
	err         error
	item_count  int
}

func buildStressURL(base_url string, ep stressEndpoint, rng *rand.Rand) string {
	params := url.Values{}
	path := ep.path

	if ep.accepts_q && rng.Intn(2) == 0 {
		params.Set("q", sample_queries[rng.Intn(len(sample_queries))])
	}

	if ep.accepts_tags && rng.Intn(2) == 0 {
		n := 1 + rng.Intn(min(3, len(sample_tags)))
		perm := rng.Perm(len(sample_tags))
		for i := 0; i < n; i++ {
			params.Add("tags", sample_tags[perm[i]])
		}
	}

	if ep.accepts_from && rng.Intn(2) == 0 {
		days_ago := 1 + rng.Intn(30)
		params.Set("from", time.Now().UTC().AddDate(0, 0, -days_ago).Format("2006-01-02"))
	}

	params.Set("limit", strconv.Itoa(1+rng.Intn(50)))

	raw := strings.TrimSuffix(base_url, "/") + path
	if enc := params.Encode(); enc != "" {
		raw += "?" + enc
	}
	return raw
}

// countEnvelopeItems counts items in a collection envelope's data array, falling back to a bare
// JSON array for resilience. Returns 0 when the body is neither.
func countEnvelopeItems(body []byte) int {
	var env struct {
		Data []any `json:"data"`
	}
	if json.Unmarshal(body, &env) == nil {
		return len(env.Data)
	}
	var arr []any
	if json.Unmarshal(body, &arr) == nil {
		return len(arr)
	}
	return 0
}

func runStressTest(base_url string, concurrency int, api_key string) []stressResult {
	results := make([]stressResult, concurrency)
	var wg sync.WaitGroup
	wg.Add(concurrency)

	client := &http.Client{Timeout: HTTP_TIMEOUT}
	master_rng := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	seeds := make([]int64, concurrency)
	for i := range seeds {
		seeds[i] = master_rng.Int63()
	}

	for i := 0; i < concurrency; i++ {
		go func(idx int) {
			defer wg.Done()

			rng := rand.New(rand.NewSource(seeds[idx])) //nolint:gosec
			ep := stress_endpoints[rng.Intn(len(stress_endpoints))]
			raw_url := buildStressURL(base_url, ep, rng)

			parsed, err := url.Parse(raw_url)
			if err != nil {
				results[idx] = stressResult{endpoint: ep.path, err: err}
				return
			}
			endpoint := parsed.Path

			req, err := http.NewRequest(http.MethodGet, raw_url, nil)
			if err != nil {
				results[idx] = stressResult{endpoint: endpoint, err: err}
				return
			}
			if api_key != "" {
				req.Header.Set("X-API-KEY", api_key)
			}

			start := time.Now()
			resp, err := client.Do(req)
			latency := time.Since(start)

			if err != nil {
				results[idx] = stressResult{endpoint: endpoint, latency: latency, err: err}
				return
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			count := countEnvelopeItems(body)

			results[idx] = stressResult{
				endpoint:    endpoint,
				status_code: resp.StatusCode,
				latency:     latency,
				item_count:  count,
			}
		}(i)
	}

	wg.Wait()
	return results
}

func printStressSummary(t *testing.T, results []stressResult) {
	t.Helper()

	type epStats struct {
		total, success, failures int
		total_ms, total_items    int64
	}

	stats := map[string]*epStats{}
	for _, ep := range stress_endpoints {
		stats[ep.path] = &epStats{}
	}

	total_success, total_failure := 0, 0
	for _, r := range results {
		s, ok := stats[r.endpoint]
		if !ok {
			s = &epStats{}
			stats[r.endpoint] = s
		}
		s.total++
		s.total_ms += r.latency.Milliseconds()
		s.total_items += int64(r.item_count)

		if r.err != nil || r.status_code >= 500 {
			s.failures++
			total_failure++
		} else {
			s.success++
			total_success++
		}
	}

	t.Log("=== Stress Test Summary ===")
	t.Logf("Total requests: %d | Success: %d | Failure: %d",
		len(results), total_success, total_failure)

	var total_items int64
	for _, r := range results {
		total_items += int64(r.item_count)
	}
	t.Logf("Total items received: %d", total_items)
	t.Log("--- Per-endpoint breakdown ---")

	seen := make([]string, 0, len(stats))
	for _, ep := range stress_endpoints {
		seen = append(seen, ep.path)
	}

	for _, path := range seen {
		s := stats[path]
		if s == nil || s.total == 0 {
			continue
		}
		avg_ms, avg_items := int64(0), int64(0)
		if s.total > 0 {
			avg_ms = s.total_ms / int64(s.total)
			avg_items = s.total_items / int64(s.total)
		}
		t.Logf("  %-32s  total=%-5d  ok=%-5d  err=%-5d  avg_latency=%dms  avg_items=%d",
			path, s.total, s.success, s.failures, avg_ms, avg_items)
	}

	if len(results) > 0 {
		fail_rate := float64(total_failure) / float64(len(results))
		if fail_rate > 0.10 {
			t.Errorf("stress test failure rate %.1f%% exceeds 10%% threshold", fail_rate*100)
		}
	}
}

func concurrencyFromEnv() int {
	raw := os.Getenv("STRESS_CONCURRENCY")
	if raw == "" {
		return 200
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < MIN_CONCURRENCY {
		return MIN_CONCURRENCY
	}
	if n > MAX_CONCURRENCY {
		return MAX_CONCURRENCY
	}
	return n
}

func stressBaseURL() string {
	if base_url := os.Getenv("STRESS_BASE_URL"); base_url != "" {
		return base_url
	}
	return DEFAULT_STRESS_BASE_URL
}

func skipIfStressServerUnreachable(t *testing.T, base_url string) {
	t.Helper()
	client := &http.Client{Timeout: HTTP_TIMEOUT}
	if _, err := client.Get(base_url + "/health"); err != nil {
		t.Skipf("API server not reachable at %s (%v) — skipping stress test", base_url, err)
	}
}

func stressEndpointFailures(results []stressResult) map[string]int {
	failures := map[string]int{}
	for _, r := range results {
		if r.err != nil || r.status_code >= 500 {
			failures[r.endpoint]++
		}
	}
	return failures
}

func TestStressAPI(t *testing.T) {
	base_url := stressBaseURL()
	api_key := os.Getenv("STRESS_API_KEY")
	concurrency := concurrencyFromEnv()

	t.Logf("Stress testing %s with %d concurrent requests", base_url, concurrency)
	skipIfStressServerUnreachable(t, base_url)

	results := runStressTest(base_url, concurrency, api_key)
	printStressSummary(t, results)
}

func TestStressAPIEndpoints(t *testing.T) {
	base_url := stressBaseURL()
	api_key := os.Getenv("STRESS_API_KEY")
	skipIfStressServerUnreachable(t, base_url)

	const REQUESTS_PER_ENDPOINT = 10
	concurrency := len(stress_endpoints) * REQUESTS_PER_ENDPOINT

	t.Logf("Endpoint smoke stress: %d endpoints × %d requests = %d total",
		len(stress_endpoints), REQUESTS_PER_ENDPOINT, concurrency)

	results := runStressTest(base_url, concurrency, api_key)
	printStressSummary(t, results)

	failures := stressEndpointFailures(results)
	for _, ep := range stress_endpoints {
		path := ep.path
		t.Run(fmt.Sprintf("endpoint=%s", path), func(t *testing.T) {
			if f := failures[path]; f > 0 {
				t.Errorf("%s had %d failure(s)", path, f)
			}
		})
	}
}

func TestStressVectorSearch(t *testing.T) {
	base_url := stressBaseURL()
	api_key := os.Getenv("STRESS_API_KEY")
	concurrency := concurrencyFromEnv()
	skipIfStressServerUnreachable(t, base_url)

	t.Logf("Vector search stress testing with %d concurrent requests", concurrency)

	results := make([]stressResult, concurrency)
	var wg sync.WaitGroup
	wg.Add(concurrency)

	client := &http.Client{Timeout: HTTP_TIMEOUT}
	master_rng := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	seeds := make([]int64, concurrency)
	for i := range seeds {
		seeds[i] = master_rng.Int63()
	}

	vector_endpoints := []string{ROUTE_EVENTS, ROUTE_SIGNALS}

	for i := 0; i < concurrency; i++ {
		go func(idx int) {
			defer wg.Done()

			rng := rand.New(rand.NewSource(seeds[idx])) //nolint:gosec
			endpoint := vector_endpoints[rng.Intn(len(vector_endpoints))]

			params := url.Values{}
			params.Set("q", sample_queries[rng.Intn(len(sample_queries))])
			params.Set("limit", strconv.Itoa(1+rng.Intn(50)))

			raw_url := strings.TrimSuffix(base_url, "/") + endpoint + "?" + params.Encode()
			req, err := http.NewRequest(http.MethodGet, raw_url, nil)
			if err != nil {
				results[idx] = stressResult{endpoint: endpoint, err: err}
				return
			}
			if api_key != "" {
				req.Header.Set("X-API-KEY", api_key)
			}

			start := time.Now()
			resp, err := client.Do(req)
			latency := time.Since(start)

			if err != nil {
				results[idx] = stressResult{endpoint: endpoint, latency: latency, err: err}
				return
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			count := countEnvelopeItems(body)

			results[idx] = stressResult{
				endpoint:    endpoint,
				status_code: resp.StatusCode,
				latency:     latency,
				item_count:  count,
			}
		}(i)
	}

	wg.Wait()
	printStressSummary(t, results)
}

func TestRouterDiscoveryRoutes(t *testing.T) {
	srv := newTestHTTPServer(t)
	for _, path := range []string{"/entities", "/regions", "/event-types"} {
		status, body := routerGET(t, srv.URL, path, url.Values{"limit": {"5"}}, "")
		requireStatus(t, http.StatusOK, status, body)
		values := assertExpectedPagination(t, body, 5, "")
		require.NotEmpty(t, values, path)
		value, ok := values[0]["value"].(string)
		require.True(t, ok, "missing discovery value for %s", path)
		assert.NotEmpty(t, value, path)
		type_value, ok := values[0]["type"].(string)
		require.True(t, ok, "missing discovery type for %s", path)
		switch path {
		case "/entities":
			assert.Contains(t, []string{"company", "people"}, type_value)
		case "/regions":
			assert.Equal(t, "region", type_value)
		case "/event-types":
			assert.Equal(t, "event_type", type_value)
		}
	}
}
