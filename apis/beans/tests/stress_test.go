package gobeansack_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const DEFAULT_STRESS_BASE_URL = "http://localhost:8080"

const (
	MIN_CONCURRENCY = 100
	MAX_CONCURRENCY = 10000
	HTTP_TIMEOUT    = 10 * time.Minute
)

type stressEndpoint struct {
	path         string
	accepts_q    bool
	accepts_tags bool
	accepts_from bool
}

var stress_endpoints = []stressEndpoint{
	{path: ROUTE_SEARCH, accepts_q: true, accepts_tags: true, accepts_from: true},
	{path: ROUTE_LATEST, accepts_q: true, accepts_tags: true},
	{path: ROUTE_TRENDING, accepts_q: true, accepts_tags: true},
	{path: ROUTE_HEADLINES, accepts_q: true, accepts_tags: true},
	{path: ROUTE_SOURCES, accepts_q: true},
	{path: ROUTE_CATEGORIES, accepts_q: true},
	{path: ROUTE_ENTITIES, accepts_q: true},
	{path: ROUTE_REGIONS, accepts_q: true},
	{path: ROUTE_SENTIMENTS, accepts_q: true},
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
	"data privacy and GDPR",
	"generative AI models",
	"zero trust security",
	"edge computing",
	"5G networks",
	"sustainable technology",
	"DevOps practices",
	"API security",
	"supply chain resilience",
	"digital transformation",
}

var sample_tags = []string{
	"openai",
	"google",
	"microsoft",
	"united_states",
	"europe",
	"artificial_intelligence",
	"cybersecurity",
	"public_policy_and_administration",
}

type stressResult struct {
	endpoint    string
	status_code int
	latency     time.Duration
	err         error
	item_count  int
}

func buildStressURL(base_url string, ep stressEndpoint, rng *rand.Rand) string {
	params := url.Values{}

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

	raw := strings.TrimSuffix(base_url, "/") + ep.path
	if enc := params.Encode(); enc != "" {
		raw += "?" + enc
	}
	return raw
}

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

			req, err := http.NewRequest(http.MethodGet, raw_url, nil)
			if err != nil {
				results[idx] = stressResult{endpoint: ep.path, err: err}
				return
			}
			if api_key != "" {
				req.Header.Set("X-API-KEY", api_key)
			}

			start := time.Now()
			resp, err := client.Do(req)
			latency := time.Since(start)
			if err != nil {
				results[idx] = stressResult{endpoint: ep.path, latency: latency, err: err}
				return
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			results[idx] = stressResult{
				endpoint:    ep.path,
				status_code: resp.StatusCode,
				latency:     latency,
				item_count:  countEnvelopeItems(body),
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
	for _, ep := range stress_endpoints {
		s := stats[ep.path]
		avg_ms, avg_items := int64(0), int64(0)
		if s.total > 0 {
			avg_ms = s.total_ms / int64(s.total)
			avg_items = s.total_items / int64(s.total)
		}
		t.Logf("  %-28s  total=%-5d  ok=%-5d  err=%-5d  avg_latency=%dms  avg_items=%d",
			ep.path, s.total, s.success, s.failures, avg_ms, avg_items)
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

	failures := map[string]int{}
	for _, r := range results {
		if r.err != nil || r.status_code >= 500 {
			failures[r.endpoint]++
		}
	}
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

	vector_endpoints := []string{ROUTE_SEARCH, ROUTE_LATEST, ROUTE_TRENDING}

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
			results[idx] = stressResult{
				endpoint:    endpoint,
				status_code: resp.StatusCode,
				latency:     latency,
				item_count:  countEnvelopeItems(body),
			}
		}(i)
	}

	wg.Wait()
	printStressSummary(t, results)
}
