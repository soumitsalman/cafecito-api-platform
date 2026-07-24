package gobeansack_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/soumitsalman/cafecito-api-platform/apis/beans/db"
	"github.com/soumitsalman/cafecito-api-platform/apis/internal/config"
	"github.com/stretchr/testify/assert"

	"github.com/k0kubun/pp"
)

var test_ctx = context.Background()

func TestBuildVectorSQLUsesHNSWCandidateQuery(t *testing.T) {
	distance := 0.25
	conditions := db.Condition{
		Categories: []string{"technology"},
		Created:    time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
		Embedding:  []float32{0.1, 0.2, 0.3},
		Distance:   &distance,
	}

	query, params := (&db.PGSack{}).BuildVectorSQL(
		db.BEANS,
		conditions,
		[]string{db.ORDER_BY_LATEST},
		db.Pagination{Limit: 16, Offset: 4},
		"*",
	)

	expected_parts := []string{
		"WITH nearest_results AS MATERIALIZED",
		"ORDER BY embedding <=> @embedding::vector ASC",
		"LIMIT @candidate_limit",
		"WHERE distance <= @distance",
		"ORDER BY created DESC, distance ASC",
		"LIMIT @limit OFFSET @offset",
	}
	for _, part := range expected_parts {
		if !strings.Contains(query, part) {
			t.Errorf("query does not contain %q:\n%s", part, query)
		}
	}
	expected_candidate_limit := (16 + 4) * config.VECTOR_QUERY_CANDIDATE_LIMIT_MULTIPLIER
	if got := params["candidate_limit"]; got != expected_candidate_limit {
		t.Errorf("candidate_limit = %v, want %d", got, expected_candidate_limit)
	}
}

func TestCountRowsRejectsEmbeddingWithoutDistance(t *testing.T) {
	pg_sack := &db.PGSack{}
	count, err := pg_sack.CountRows(context.Background(), db.BEANS, db.Condition{
		Embedding: []float32{0.1, 0.2, 0.3},
	})

	assert.ErrorIs(t, err, db.ErrVectorDistanceRequired)
	assert.Zero(t, count)
}

func TestQueryDistinctValues(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()
	page := db.Pagination{Limit: 5}
	regions, err := pg_sack.DistinctRegions(context.Background(), page)
	assert.NoError(t, err)
	pp.Println("REGIONS", regions)

	entities, err := pg_sack.DistinctEntities(context.Background(), page)
	assert.NoError(t, err)
	pp.Println("ENTITIES", entities)

	sources, err := pg_sack.DistinctSources(context.Background(), page)
	assert.NoError(t, err)
	pp.Println("SOURCES", sources)

	categories, err := pg_sack.DistinctCategories(context.Background(), page)
	assert.NoError(t, err)
	pp.Println("CATEGORIES", categories)
}

func TestQueryChatterStats(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()
	urls := []string{
		"https://www.slashgear.com/1896648/lifesaber-emergency-tool-usb-powered-features/",
		"https://www.wusa9.com/article/news/nation-world/trump-big-bill-may-have-political-cost/507-0c07fc4b-248b-4a3b-96c0-81d75a511228",
		"https://issuepay.app",
		"https://minutemirror.com.pk/lahore-high-court-halts-transfer-of-brazilian-monkeys-to-lahore-zoo-406544/",
		"https://jameshard.ing/pilot",
		"https://jobsbyreferral.com/",
		"https://llmapitest.com/",
		"https://htmlrev.com/",
	}
	cond := db.Condition{URLs: urls}
	page := db.Pagination{Limit: 5}
	chatters, err := pg_sack.QueryChatters(context.Background(), cond, page, nil)
	assert.NoError(t, err)
	pp.Println("CHATTERS", chatters)
	// aggregates not currently supported by the API
}

func TestQueryBeanExtensions(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()
	urls := []string{
		"https://www.slashgear.com/1896648/lifesaber-emergency-tool-usb-powered-features/",
		"https://www.wusa9.com/article/news/nation-world/trump-big-bill-may-have-political-cost/507-0c07fc4b-248b-4a3b-96c0-81d75a511228",
		"https://issuepay.app",
		"https://minutemirror.com.pk/lahore-high-court-halts-transfer-of-brazilian-monkeys-to-lahore-zoo-406544/",
		"https://jameshard.ing/pilot",
		"https://jobsbyreferral.com/",
		"https://llmapitest.com/",
		"https://htmlrev.com/",
	}
	// use the new query API to fetch beans matching the URLs
	cond := db.Condition{URLs: urls}
	page := db.Pagination{Limit: 5}
	beans, err := pg_sack.QueryLatestBeans(context.Background(), cond, page, nil)
	assert.NoError(t, err)
	pp.Println("BEAN EXTENSIONS", beans)
}

func TestQueryBeans(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()
	categories := []string{"public_policy_and_administration", "art_and_design"}

	cond := db.Condition{
		Categories: categories,
		Created:    time.Now().AddDate(0, 0, -3),
	}
	page := db.Pagination{Limit: 5}
	beans, err := pg_sack.QueryLatestBeans(context.Background(), cond, page, nil)
	if err != nil {
		t.Fatalf("QUERY BEANS ERROR: %v", err)
	}
	assert.Len(t, beans, 5)
}

func TestVectorSearch(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	distance := 0.4
	cond := db.Condition{
		Embedding: testQueryEmbedding,
		Distance:  &distance,
	}
	page := db.Pagination{Limit: 5}
	beans, err := pg_sack.QueryLatestBeans(context.Background(), cond, page, nil)
	if err != nil {
		t.Fatalf("VECTOR SEARCH ERROR: %v", err)
	}
	pp.Println("VECTOR SEARCH RESULTS", beans)
	assert.Len(t, beans, 5)
}

func TestQueryPropagation(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()
	urls := []string{
		"https://www.foxla.com/news/more-us-airlines-raise-baggage-fees-see-list",
		"https://www.slashgear.com/2143752/us-airlines-increased-fees-fuel-prices/",
		"http://andersource.dev/2026/03/29/tradeoff-sliders.html",
		"https://massivelyop.com/2026/04/08/perfect-ten-have-official-mmo-websites-improved-in-the-last-decade/",
		"https://phys.org/news/2026-04-uncharted-island-nautical.html",
	}
	results, err := pg_sack.QueryPropagation(context.Background(), urls)
	assert.NoError(t, err)
	assert.Len(t, results, len(urls))
	for i, r := range results {
		assert.Equal(t, urls[i], r.URL)
		assert.NotNil(t, r.Coverage)
		assert.NotNil(t, r.Mentions)
	}
	pp.Println("PROPAGATION", results)
}
