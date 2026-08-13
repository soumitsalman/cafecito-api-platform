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
	page, err := pg_cupboard.QuerySips(context.Background(), filters, db.PageRequest{Limit: 16})
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
	page, err := pg_cupboard.QuerySips(context.Background(), filters, db.PageRequest{Limit: 16})
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
	page, err := pg_cupboard.QuerySips(context.Background(), filters, db.PageRequest{Limit: 5})
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
	page, err := pg_cupboard.QuerySips(context.Background(), filters, db.PageRequest{Limit: 5})
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

	list, err := pg_cupboard.QuerySips(context.Background(), db.Filters{Kind: db.SIP_KIND_EVENT}, db.PageRequest{Limit: 1})
	require.NoError(t, err)
	require.NotEmpty(t, list.Items)
	event_id := list.Items[0].ID

	event, err := pg_cupboard.GetSip(context.Background(), event_id, db.SIP_KIND_EVENT)
	require.NoError(t, err)
	assert.False(t, event.IsZero())
	assert.Equal(t, db.SIP_KIND_EVENT, event.Kind)

	exists, err := pg_cupboard.SipExists(context.Background(), event_id, db.SIP_KIND_EVENT)
	require.NoError(t, err)
	assert.True(t, exists)

	evidence, err := pg_cupboard.QuerySameSips(context.Background(), event_id, db.Filters{}, db.PageRequest{Limit: 5})
	assert.NoError(t, err)
	pp.Println("EVIDENCE", evidence)

	signals, err := pg_cupboard.QueryDerivedSips(context.Background(), event_id, db.Filters{}, db.PageRequest{Limit: 5})
	assert.NoError(t, err)
	pp.Println("DERIVED_SIGNALS", signals)

	counts, err := pg_cupboard.CountRelations(context.Background(), event_id)
	assert.NoError(t, err)
	pp.Println("RELATION_COUNTS", counts)
}

func TestGetEventNotFound(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	event, err := pg_cupboard.GetSip(context.Background(), uuid.New(), db.SIP_KIND_EVENT)
	assert.NoError(t, err)
	assert.True(t, event.IsZero())
}

func TestCursorPaginationEvents(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	first, err := pg_cupboard.QuerySips(context.Background(), db.Filters{
		Kind:        db.SIP_KIND_EVENT,
		CreatedFrom: time.Now().AddDate(0, 0, -30),
	}, db.PageRequest{Limit: 2})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(first.Items), 1)

	if first.NextCursor == nil {
		t.Skip("not enough events for a second page")
	}

	second, err := pg_cupboard.QuerySips(context.Background(), db.Filters{
		Kind:        db.SIP_KIND_EVENT,
		CreatedFrom: testSearchFrom(),
	}, db.PageRequest{Limit: 2, Cursor: first.NextCursor})
	require.NoError(t, err)
	require.NotEmpty(t, second.Items)
	assert.NotEqual(t, first.Items[0].ID, second.Items[0].ID)
}

func TestDiscoveryQueries(t *testing.T) {
	pg_cupboard := setupTestDB()
	defer pg_cupboard.Close()

	for _, query := range []func() (db.Page[db.TagValue], error){
		func() (db.Page[db.TagValue], error) {
			return pg_cupboard.QueryEventTags(context.Background(), "", []string{db.EVENT_TAG_TYPE_COMPANY, db.EVENT_TAG_TYPE_PEOPLE}, db.PageRequest{Limit: 5})
		},
		func() (db.Page[db.TagValue], error) {
			return pg_cupboard.QueryEventTags(context.Background(), "", []string{db.EVENT_TAG_TYPE_REGION}, db.PageRequest{Limit: 5})
		},
		func() (db.Page[db.TagValue], error) {
			return pg_cupboard.QueryEventTags(context.Background(), "", []string{db.EVENT_TAG_TYPE_EVENT_TYPE}, db.PageRequest{Limit: 5})
		},
	} {
		page, err := query()
		require.NoError(t, err)
		require.NotEmpty(t, page.Items)
		assert.NotEmpty(t, page.Items[0].Value)
	}
}
