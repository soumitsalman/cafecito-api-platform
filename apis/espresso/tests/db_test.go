package espressoapi_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/k0kubun/pp"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/db"
	"github.com/soumitsalman/cafecito-api-platform/apis/internal/config"
	"github.com/stretchr/testify/assert"
)

var test_ctx = context.Background()

func TestGetTags(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()
	page := db.Pagination{Limit: 5, Offset: 10}
	tags, err := pg_cupboard.GetTags(context.Background(), page)
	assert.NoError(t, err)
	pp.Println("TAGS", tags)
}

func TestRelatedSips(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()
	cond := db.Condition{
		IDs:          testRelatedIDs,
		Relationship: "SAME_AS",
	}
	page := db.Pagination{}
	sips, err := pg_cupboard.QueryRelatedSips(context.Background(), cond, page)
	assert.NoError(t, err)
	assert.Greater(t, len(sips), 0)
	pp.Println("RELATED SIPS", sips)
}

func TestScalarSearchSips(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()
	cond := db.Condition{
		Tags:    testScalarTags,
		Created: testSearchFrom(),
		Kinds:   db.EVENTS,
	}
	page := db.Pagination{}
	sips, err := pg_cupboard.QuerySips(context.Background(), cond, page)
	assert.NoError(t, err)
	assert.Greater(t, len(sips), 0)
	pp.Println("SIPS", sips)
}

func TestTextSearchSips(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()
	cond := db.Condition{
		Tags:    testTextTags,
		FTS:     true,
		Created: testSearchFrom(),
	}
	page := db.Pagination{}
	sips, err := pg_cupboard.QuerySips(context.Background(), cond, page)
	assert.NoError(t, err)
	assert.Greater(t, len(sips), 0)
	pp.Println("SIPS", sips)
}

func TestVectorSearchSips(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()
	distance := 0.4
	cond := db.Condition{
		Embedding: testQueryEmbedding,
		Distance:  &distance,
	}
	page := db.Pagination{Limit: 5}
	sips, err := pg_cupboard.QuerySips(context.Background(), cond, page)
	assert.NoError(t, err)
	assert.Greater(t, len(sips), 0)
	pp.Println("SIPS", sips)
}

func TestBuildVectorSQLUsesHNSWCandidateQuery(t *testing.T) {
	distance := 0.0
	conditions := db.Condition{
		Kinds:     []string{"signal"},
		Created:   time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
		Tags:      []string{"markets"},
		Embedding: []float32{0.1, 0.2, 0.3},
		Distance:  &distance,
	}

	query, params := db.BuildVectorSQL(conditions, db.Pagination{Limit: 16, Offset: 4})

	expected_parts := []string{
		"WITH nearest_results AS MATERIALIZED",
		"sips.kind = ANY(@kinds)",
		"sips.created >= @created",
		"sips.tags && @tags",
		"ORDER BY sips.embedding <=> @embedding ASC",
		"LIMIT @candidate_limit",
		"WHERE distance <= @distance",
		"ORDER BY distance ASC",
		"LIMIT @limit",
		"OFFSET @offset",
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
	if got := params["distance"]; got != 0.0 {
		t.Errorf("distance = %v, want 0", got)
	}
}

func TestBuildScalarSQLDoesNotHandleEmbedding(t *testing.T) {
	distance := 0.25
	query, params := db.BuildScalarSQL(db.SIPS, db.Condition{
		Embedding: []float32{0.1, 0.2, 0.3},
		Distance:  &distance,
	}, db.Pagination{Limit: 16}, "id, created, digest")

	if strings.Contains(query, "<=>") {
		t.Fatalf("scalar query unexpectedly contains vector distance:\n%s", query)
	}
	if _, ok := params["embedding"]; ok {
		t.Fatal("scalar query unexpectedly includes embedding parameter")
	}
}

func TestBuildVectorCountSQLKeepsExactDistanceFilterSeparate(t *testing.T) {
	distance := 0.25
	query, params := db.BuildVectorCountSQL(db.Condition{
		Embedding: []float32{0.1, 0.2, 0.3},
		Distance:  &distance,
	})

	if !strings.Contains(query, "(sips.embedding <=> @embedding) <= @distance") {
		t.Fatalf("vector count query does not contain distance filter:\n%s", query)
	}
	if strings.Contains(query, "ORDER BY") {
		t.Fatalf("vector count query unexpectedly contains ordering:\n%s", query)
	}
	if got := params["distance"]; got != distance {
		t.Errorf("distance = %v, want %v", got, distance)
	}
}

func TestCountSipsRejectsEmbeddingWithoutDistance(t *testing.T) {
	pg_cupboard := &db.Cupboard{}
	count, err := pg_cupboard.CountSips(context.Background(), db.Condition{
		Embedding: []float32{0.1, 0.2, 0.3},
	})

	assert.ErrorIs(t, err, db.ErrVectorDistanceRequired)
	assert.Zero(t, count)
}
