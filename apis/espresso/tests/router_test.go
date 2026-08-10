package espressoapi_test

import (
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
	datautils "github.com/soumitsalman/data-utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const DEFAULT_STRESS_BASE_URL = "http://localhost:8080"

const (
	MIN_CONCURRENCY = 100
	MAX_CONCURRENCY = 10000
	HTTP_TIMEOUT    = 10 * time.Minute

	ROUTE_TAGS    = "/tags"
	ROUTE_EVENTS  = "/events"
	ROUTE_SIGNALS = "/signals"
	ROUTE_SOURCES = "/sources"
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

func TestRouterGetTags(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("q", "academic")
	params.Set("resource", "event")
	params.Set("limit", "5")

	status, body := routerGET(t, srv.URL, ROUTE_TAGS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	tags := parseStringArray(t, body)
	assert.NotEmpty(t, tags)
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
	events := parseDigestArray(t, body)
	assert.Greater(t, len(events), 0)
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
	events := parseDigestArray(t, body)
	assert.Greater(t, len(events), 0)
	pp.Println("EVENTS", events)
}

func TestRouterVectorSearchEvents(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("q", TEST_VECTOR_QUERY)
	params.Set("acc", "0.6")
	params.Set("limit", "5")
	params.Set("from", testSearchFrom().Format("2006-01-02"))

	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	events := parseDigestArray(t, body)
	assert.Greater(t, len(events), 0)
	pp.Println("EVENTS", events)
}

func TestRouterVectorSearchSignals(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("q", TEST_VECTOR_QUERY)
	params.Set("acc", "0.7")
	params.Set("limit", "5")
	params.Set("to", testSearchFrom().Format("2006-01-02"))

	status, body := routerGET(t, srv.URL, ROUTE_SIGNALS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	signals := parseDigestArray(t, body)
	assert.Greater(t, len(signals), 0)
	pp.Println("SIGNALS", signals)
}

func TestRouterGetSources(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("limit", "5")

	status, body := routerGET(t, srv.URL, ROUTE_SOURCES, params, "")
	requireStatus(t, http.StatusOK, status, body)
	sources := parseSourceArray(t, body)
	assert.NotEmpty(t, sources)
	pp.Println("SOURCES", sources)

	first_id, ok := sources[0]["id"].(string)
	require.True(t, ok)
	status, body = routerGET(t, srv.URL, ROUTE_SOURCES+"/"+first_id, nil, "")
	requireStatus(t, http.StatusOK, status, body)
	src := parseDetailObject(t, body)
	assert.Equal(t, first_id, src["id"])
}

func TestRouterSignalDetailsAndEvents(t *testing.T) {
	srv := newTestHTTPServer(t)

	params := url.Values{}
	params.Set("limit", "1")
	status, body := routerGET(t, srv.URL, ROUTE_SIGNALS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	signals := parseDigestArray(t, body)
	require.NotEmpty(t, signals)
	signal_id, ok := signals[0]["id"].(string)
	require.True(t, ok)

	status, body = routerGET(t, srv.URL, ROUTE_SIGNALS+"/"+signal_id, nil, "")
	requireStatus(t, http.StatusOK, status, body)
	detail := parseDetailObject(t, body)
	assert.Equal(t, signal_id, detail["id"])
	assert.Contains(t, detail, "links")
	assert.Contains(t, detail, "counts")
	pp.Println("DETAIL", detail)

	status, body = routerGET(t, srv.URL, ROUTE_SIGNALS+"/"+signal_id+"/events", nil, "")
	requireStatus(t, http.StatusOK, status, body)
	events := parseDigestArray(t, body)
	assert.NotEmpty(t, events)
	pp.Println("EVENTS", events)
}

func TestRouterEventDetailEvidenceAndSignals(t *testing.T) {
	srv := newTestHTTPServer(t)

	// get a root signal
	params := url.Values{}
	params.Set("limit", "1")
	status, body := routerGET(t, srv.URL, ROUTE_SIGNALS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	root_signals := parseDigestArray(t, body)
	require.NotEmpty(t, root_signals)
	signal_id, ok := root_signals[0]["id"].(string)
	require.True(t, ok)

	status, body = routerGET(t, srv.URL, ROUTE_SIGNALS+"/"+signal_id+"/events?limit=1", nil, "")
	requireStatus(t, http.StatusOK, status, body)
	events := parseDigestArray(t, body)
	assert.NotEmpty(t, events)
	event_id, ok := events[0]["id"].(string)
	require.True(t, ok)

	status, body = routerGET(t, srv.URL, ROUTE_EVENTS+"/"+event_id, nil, "")
	requireStatus(t, http.StatusOK, status, body)
	detail := parseDetailObject(t, body)
	assert.Equal(t, event_id, detail["id"])
	assert.Contains(t, detail, "links")
	assert.Contains(t, detail, "counts")
	pp.Println("DETAIL", detail)

	status, body = routerGET(t, srv.URL, ROUTE_EVENTS+"/"+event_id+"/evidence", nil, "")
	requireStatus(t, http.StatusOK, status, body)
	evidence := parseDigestArray(t, body)
	assert.NotEmpty(t, evidence)
	pp.Println("EVIDENCE", evidence)

	status, body = routerGET(t, srv.URL, ROUTE_EVENTS+"/"+event_id+"/signals", nil, "")
	requireStatus(t, http.StatusOK, status, body)
	signals := parseDigestArray(t, body)
	assert.NotEmpty(t, signals)
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
}

func TestRouterCursorPagination(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("limit", "2")
	params.Set("from", time.Now().UTC().AddDate(0, 0, -30).Format("2006-01-02"))

	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	first := parseDigestArray(t, body)
	require.NotEmpty(t, first)

	cursor := nextCursorFromBody(t, body)
	if cursor == "" {
		t.Skip("not enough events for a second page")
	}

	params.Set("cursor", cursor)
	status, body = routerGET(t, srv.URL, ROUTE_EVENTS, params, "")
	requireStatus(t, http.StatusOK, status, body)
	second := parseDigestArray(t, body)
	require.NotEmpty(t, second)
	assert.NotEqual(t, first[0]["id"], second[0]["id"])
}

func TestRouterInvalidCursor(t *testing.T) {
	srv := newTestHTTPServer(t)
	params := url.Values{}
	params.Set("cursor", "not-a-valid-cursor")

	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, params, "")
	requireStatus(t, http.StatusBadRequest, status, body)
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
