package espressoapi_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryTags(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	page, err := pg_cupboard.QueryTags(context.Background(), "academic", []string{db.SIP_KIND_EVENT}, db.PageRequest{Limit: 200})
	assert.NoError(t, err)
	assert.NotEmpty(t, page.Items)
	pp.Println("TAGS", page)
}

func TestScalarSearchEvents(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	filters := db.Filters{
		Tags:        test_scalar_tags,
		Regions:     []string{"us", "japan", "china"},
		CreatedFrom: testSearchFrom(),
	}
	page, err := pg_cupboard.QueryEvents(context.Background(), filters, db.PageRequest{Limit: 16})
	assert.NoError(t, err)
	assert.Greater(t, len(page.Items), 0)
	pp.Println("EVENTS", page)
}

func TestScalarSearchSignals(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	filters := db.Filters{
		Tags:        test_scalar_tags,
		CreatedFrom: testSearchFrom(),
	}
	page, err := pg_cupboard.QuerySignals(context.Background(), filters, db.PageRequest{Limit: 16})
	assert.NoError(t, err)
	assert.Greater(t, len(page.Items), 0)
	pp.Println("SIGNALS", page)
}

func TestVectorSearchEvents(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	filters := db.Filters{
		Embedding:   test_query_embedding,
		Tags:        test_scalar_tags,
		CreatedFrom: testSearchFrom(),
	}
	page, err := pg_cupboard.QueryEvents(context.Background(), filters, db.PageRequest{Limit: 5})
	assert.NoError(t, err)
	assert.Greater(t, len(page.Items), 0)
	pp.Println("EVENTS", page.Items)
}

func TestVectorSearchSignals(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	distance := 0.6
	filters := db.Filters{
		Embedding: test_query_embedding,
		Distance:  &distance,
	}
	page, err := pg_cupboard.QuerySignals(context.Background(), filters, db.PageRequest{Limit: 5})
	assert.NoError(t, err)
	assert.Greater(t, len(page.Items), 0)
	pp.Println("SIGNALS", page.Items)
}

func TestQuerySources(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	page, err := pg_cupboard.QuerySources(context.Background(), "", nil, db.PageRequest{Limit: 5})
	assert.NoError(t, err)
	assert.NotEmpty(t, page.Items)
	pp.Println("SOURCES", page)

	first := page.Items[0]
	source, err := pg_cupboard.GetSource(context.Background(), first.ID)
	assert.NoError(t, err)
	assert.False(t, source.IsZero())
	assert.Equal(t, first.ID, source.ID)
	pp.Println("SOURCE", source)
}

func TestGetEventAndRelations(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	list, err := pg_cupboard.QueryEvents(context.Background(), db.Filters{}, db.PageRequest{Limit: 1})
	require.NoError(t, err)
	require.NotEmpty(t, list.Items)
	event_id := list.Items[0].ID

	event, err := pg_cupboard.GetEvent(context.Background(), event_id)
	require.NoError(t, err)
	assert.False(t, event.IsZero())
	assert.Equal(t, db.SIP_KIND_EVENT, event.Kind)

	exists, err := pg_cupboard.EventExists(context.Background(), event_id)
	require.NoError(t, err)
	assert.True(t, exists)

	evidence, err := pg_cupboard.QueryEventEvidence(context.Background(), event_id, db.Filters{}, db.PageRequest{Limit: 5})
	assert.NoError(t, err)
	pp.Println("EVIDENCE", evidence)

	signals, err := pg_cupboard.QueryDerivedSignals(context.Background(), event_id, db.Filters{}, db.PageRequest{Limit: 5})
	assert.NoError(t, err)
	pp.Println("DERIVED_SIGNALS", signals)

	counts, err := pg_cupboard.CountRelations(context.Background(), event_id)
	assert.NoError(t, err)
	pp.Println("RELATION_COUNTS", counts)
}

func TestGetEventNotFound(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	event, err := pg_cupboard.GetEvent(context.Background(), uuid.New())
	assert.NoError(t, err)
	assert.True(t, event.IsZero())
}

func TestCursorPaginationEvents(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	first, err := pg_cupboard.QueryEvents(context.Background(), db.Filters{
		CreatedFrom: time.Now().AddDate(0, 0, -30),
	}, db.PageRequest{Limit: 2})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(first.Items), 1)

	if first.NextCursor == nil {
		t.Skip("not enough events for a second page")
	}

	second, err := pg_cupboard.QueryEvents(context.Background(), db.Filters{
		CreatedFrom: time.Now().AddDate(0, 0, -30),
	}, db.PageRequest{Limit: 2, Cursor: first.NextCursor})
	require.NoError(t, err)
	require.NotEmpty(t, second.Items)
	assert.NotEqual(t, first.Items[0].ID, second.Items[0].ID)
}
