package gobeansack_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"github.com/soumitsalman/cafecito-api-platform/apis/beans/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeCursor(t *testing.T) {
	id := uuid.New()
	created := time.Date(2026, time.August, 17, 12, 0, 0, 0, time.UTC)
	trend_score := 91.7
	distance := 0.25
	text_key := "publisher.example"

	cases := []db.Cursor{
		{Version: 1, ID: &id, Created: &created},
		{Version: 1, ID: &id, TrendScore: &trend_score},
		{Version: 1, ID: &id, Distance: &distance},
		{Version: 1, TextKey: &text_key},
	}
	for _, original := range cases {
		encoded, err := original.Encode()
		require.NoError(t, err)
		require.NotEmpty(t, encoded)

		decoded, err := db.DecodeCursor(encoded)
		require.NoError(t, err)
		require.NotNil(t, decoded)
		assert.Equal(t, original.Version, decoded.Version)
		if original.ID != nil {
			require.NotNil(t, decoded.ID)
			assert.Equal(t, *original.ID, *decoded.ID)
		}
		if original.Created != nil {
			require.NotNil(t, decoded.Created)
			assert.True(t, original.Created.Equal(*decoded.Created))
		}
		if original.TrendScore != nil {
			require.NotNil(t, decoded.TrendScore)
			assert.Equal(t, *original.TrendScore, *decoded.TrendScore)
		}
		if original.Distance != nil {
			require.NotNil(t, decoded.Distance)
			assert.Equal(t, *original.Distance, *decoded.Distance)
		}
		if original.TextKey != nil {
			require.NotNil(t, decoded.TextKey)
			assert.Equal(t, *original.TextKey, *decoded.TextKey)
		}
	}
}

func TestDecodeCursorEmptyAndInvalid(t *testing.T) {
	decoded, err := db.DecodeCursor("")
	require.NoError(t, err)
	assert.Nil(t, decoded)

	decoded, err = db.DecodeCursor("   ")
	require.NoError(t, err)
	assert.Nil(t, decoded)

	decoded, err = db.DecodeCursor("not-a-valid-cursor")
	assert.ErrorIs(t, err, db.ErrInvalidCursor)
	assert.Nil(t, decoded)
}

func TestQueryBeans(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	page, err := pg_sack.QueryBeans(test_ctx, db.BeanFilters{
		Categories:  []string{test_categories[0]},
		CreatedFrom: testSearchFrom(),
	}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	require.Len(t, page.Items, 5)
	for _, bean := range page.Items {
		assert.NotEqual(t, uuid.Nil, bean.ID)
		assert.NotEmpty(t, bean.URL)
	}
	pp.Println("BEANS", page.Items)
}

func TestQueryBeansUnfilteredBrowse(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	page, err := pg_sack.QueryBeans(test_ctx, db.BeanFilters{}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	assert.NotEmpty(t, page.Items)
}

func TestQueryBeansByKindAndDateRange(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	page, err := pg_sack.QueryBeans(test_ctx, db.BeanFilters{
		Kind:        "news",
		CreatedFrom: testSearchFrom(),
		CreatedTo:   testSearchTo(),
	}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	assert.NotEmpty(t, page.Items)
	for _, bean := range page.Items {
		assert.Equal(t, "news", bean.Kind)
	}
}

func TestQueryBeansByIDs(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	seed, err := pg_sack.QueryBeans(test_ctx, db.BeanFilters{
		CreatedFrom: testSearchFrom(),
	}, db.PageRequest{Limit: 2}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	require.NotEmpty(t, seed.Items)

	ids := make([]uuid.UUID, 0, len(seed.Items))
	for _, bean := range seed.Items {
		ids = append(ids, bean.ID)
	}

	page, err := pg_sack.QueryBeans(test_ctx, db.BeanFilters{IDs: ids}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	require.Len(t, page.Items, len(ids))
	got := map[uuid.UUID]bool{}
	for _, bean := range page.Items {
		got[bean.ID] = true
	}
	for _, id := range ids {
		assert.True(t, got[id], "missing id %s", id)
	}
}

func TestQueryBeansEnrichmentFilters(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	page, err := pg_sack.QueryBeans(test_ctx, db.BeanFilters{
		Regions:     []string{test_regions[2]},
		Entities:    []string{test_entities[0]},
		Tags:        []string{test_tags[0]},
		Sentiments:  []string{test_sentiments[0]},
		CreatedFrom: testSearchFrom(),
	}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	pp.Println("BEANS_ENRICHMENT", page.Items)
}

func TestQueryTrendingBeans(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	page, err := pg_sack.QueryTrendingBeans(test_ctx, db.BeanFilters{
		CreatedFrom: testSearchFrom(),
	}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITH_TREND)
	pp.Println("TRENDING", page.Items)
	require.NoError(t, err)
	assert.NotEmpty(t, page.Items)
	for _, bean := range page.Items {
		assert.NotEqual(t, uuid.Nil, bean.ID)
		assert.True(t, bean.TrendScore.Valid)
	}

}

func TestVectorSearchBeans(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	distance := 0.4
	page, err := pg_sack.QueryBeans(test_ctx, db.BeanFilters{
		Embedding:   test_query_embedding,
		Distance:    distance,
		CreatedFrom: testSearchFrom(),
	}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	require.Len(t, page.Items, 5)
	for _, bean := range page.Items {
		assert.True(t, bean.Distance.Valid)
		assert.LessOrEqual(t, bean.Distance.Float64, distance)
	}
	pp.Println("VECTOR_SEARCH", page.Items)
}

func TestVectorSearchLatestBeans(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	distance := 0.4
	page, err := pg_sack.QueryLatestBeans(test_ctx, db.BeanFilters{
		Embedding: test_query_embedding,
		Distance:  distance,
	}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	assert.NotEmpty(t, page.Items)
	pp.Println("VECTOR_LATEST", page.Items)
}

func TestVectorSearchTrendingBeans(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	distance := 0.4
	page, err := pg_sack.QueryTrendingBeans(test_ctx, db.BeanFilters{
		Embedding: test_query_embedding,
		Distance:  distance,
	}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITH_TREND)
	require.NoError(t, err)
	assert.NotEmpty(t, page.Items)
	pp.Println("VECTOR_TRENDING", page.Items)
}

func TestGetBean(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	list, err := pg_sack.QueryBeans(test_ctx, db.BeanFilters{
		CreatedFrom: testSearchFrom(),
	}, db.PageRequest{Limit: 1}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	require.NotEmpty(t, list.Items)
	want := list.Items[0]

	bean, err := pg_sack.GetBean(test_ctx, want.ID, false)
	require.NoError(t, err)
	assert.False(t, bean.IsZero())
	assert.Equal(t, want.ID, bean.ID)
	assert.Equal(t, want.URL, bean.URL)
	assert.False(t, bean.Content.Valid)

	with_content, err := pg_sack.GetBean(test_ctx, want.ID, true)
	require.NoError(t, err)
	assert.Equal(t, want.ID, with_content.ID)
	pp.Println("BEAN", with_content)
}

func TestGetBeanNotFound(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	bean, err := pg_sack.GetBean(test_ctx, uuid.New(), false)
	assert.ErrorIs(t, err, db.ErrNonExistentID)
	assert.True(t, bean.IsZero())
}

func TestCursorPaginationBeans(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	filters := db.BeanFilters{CreatedFrom: time.Now().UTC().AddDate(0, 0, -30)}
	first, err := pg_sack.QueryBeans(test_ctx, filters, db.PageRequest{Limit: 2}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	require.NotEmpty(t, first.Items)
	if first.NextCursor == nil {
		t.Skip("not enough articles for a second page")
	}

	second, err := pg_sack.QueryBeans(test_ctx, filters, db.PageRequest{Limit: 2, Cursor: first.NextCursor}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	require.NotEmpty(t, second.Items)
	assert.NotEqual(t, first.Items[0].ID, second.Items[0].ID)
}

func TestQuerySources(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	page, err := pg_sack.QuerySources(test_ctx, db.SourceFilters{}, db.PageRequest{Limit: 5}, db.SOURCE_COLUMNS_ALL)
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	for _, source := range page.Items {
		assert.NotEqual(t, uuid.Nil, source.ID)
		assert.NotEmpty(t, source.BaseURL)
	}
	pp.Println("SOURCES", page.Items)

	filtered, err := pg_sack.QuerySources(test_ctx, db.SourceFilters{Q: test_source_query, Domains: test_domains}, db.PageRequest{Limit: 5}, db.SOURCE_COLUMNS_ALL)
	require.NoError(t, err)
	pp.Println("SOURCES_FILTERED", filtered.Items)
}

func TestGetSource(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	list, err := pg_sack.QuerySources(test_ctx, db.SourceFilters{}, db.PageRequest{Limit: 1}, db.SOURCE_COLUMNS_ALL)
	require.NoError(t, err)
	require.NotEmpty(t, list.Items)
	want := list.Items[0]

	source, err := pg_sack.GetSource(test_ctx, want.ID)
	require.NoError(t, err)
	assert.False(t, source.IsZero())
	assert.Equal(t, want.ID, source.ID)
	pp.Println("SOURCE", source)
}

func TestGetSourceNotFound(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	source, err := pg_sack.GetSource(test_ctx, uuid.New())
	assert.ErrorIs(t, err, db.ErrNonExistentID)
	assert.True(t, source.IsZero())
}

func TestQueryTags(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	for _, tag_type := range []string{"categories", "entities", "regions", "sentiments"} {
		t.Run(tag_type, func(t *testing.T) {
			page, err := pg_sack.QueryTags(test_ctx, "", tag_type, db.PageRequest{Limit: 5})
			require.NoError(t, err)
			require.NotEmpty(t, page.Items)
			assert.NotEmpty(t, page.Items[0])
			pp.Println(tag_type, page.Items)
		})
	}
}

func TestQueryTagsPrefix(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	page, err := pg_sack.QueryTags(test_ctx, "tech", "categories", db.PageRequest{Limit: 5})
	require.NoError(t, err)
	for _, value := range page.Items {
		assert.True(t, len(value) >= 4)
	}
	pp.Println("CATEGORY_PREFIX", page.Items)
}

func TestQuerySimilarBeans(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	seed, err := pg_sack.QueryTrendingBeans(test_ctx, db.BeanFilters{
		CreatedFrom: testSearchFrom(),
	}, db.PageRequest{Limit: 20}, db.BEAN_COLUMNS_WITH_TREND)
	require.NoError(t, err)
	require.NotEmpty(t, seed.Items)

	var page db.Page[db.Bean]
	for _, bean := range seed.Items {
		page, err = pg_sack.QuerySimilarBeans(test_ctx, bean.ID, db.BeanFilters{}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITHOUT_TREND)
		require.NoError(t, err)
		if len(page.Items) > 0 {
			break
		}
	}
	if len(page.Items) == 0 {
		t.Skip("no related coverage in the test window")
	}
	for _, bean := range page.Items {
		assert.NotEqual(t, uuid.Nil, bean.ID)
		assert.NotEmpty(t, bean.URL)
	}
	pp.Println("SIMILAR", page.Items)
}

func TestQuerySimilarBeansNotFound(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	page, err := pg_sack.QuerySimilarBeans(test_ctx, uuid.New(), db.BeanFilters{}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITHOUT_TREND)
	assert.ErrorIs(t, err, db.ErrNonExistentID)
	assert.Empty(t, page.Items)
}

func TestQuerySimilarBeansCursor(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	seed, err := pg_sack.QueryTrendingBeans(test_ctx, db.BeanFilters{
		CreatedFrom: testSearchFrom(),
	}, db.PageRequest{Limit: 20}, db.BEAN_COLUMNS_WITH_TREND)
	require.NoError(t, err)

	var first db.Page[db.Bean]
	var seed_id uuid.UUID
	for _, bean := range seed.Items {
		first, err = pg_sack.QuerySimilarBeans(test_ctx, bean.ID, db.BeanFilters{}, db.PageRequest{Limit: 2}, db.BEAN_COLUMNS_WITHOUT_TREND)
		require.NoError(t, err)
		if first.NextCursor != nil {
			seed_id = bean.ID
			break
		}
	}
	if first.NextCursor == nil {
		t.Skip("not enough related articles for a second page")
	}

	second, err := pg_sack.QuerySimilarBeans(test_ctx, seed_id, db.BeanFilters{}, db.PageRequest{Limit: 2, Cursor: first.NextCursor}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	require.NotEmpty(t, second.Items)
	assert.NotEqual(t, first.Items[0].ID, second.Items[0].ID)
}

func TestQueryMentions(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	seed, err := pg_sack.QueryTrendingBeans(test_ctx, db.BeanFilters{
		CreatedFrom: testSearchFrom(),
	}, db.PageRequest{Limit: 20}, db.BEAN_COLUMNS_WITH_TREND)
	require.NoError(t, err)
	require.NotEmpty(t, seed.Items)

	var page db.Page[db.Mention]
	for _, bean := range seed.Items {
		page, err = pg_sack.QueryMentions(test_ctx, bean.ID, db.MentionFilters{}, db.PageRequest{Limit: 5})
		require.NoError(t, err)
		if len(page.Items) > 0 {
			break
		}
	}
	if len(page.Items) == 0 {
		t.Skip("no mentions in the test window")
	}
	for _, mention := range page.Items {
		assert.NotEmpty(t, mention.URL)
		assert.NotEmpty(t, mention.Platform)
		assert.False(t, mention.Observed.IsZero())
	}
	pp.Println("MENTIONS", page.Items)
}

func TestQueryMentionsNotFound(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	page, err := pg_sack.QueryMentions(test_ctx, uuid.New(), db.MentionFilters{}, db.PageRequest{Limit: 5})
	assert.ErrorIs(t, err, db.ErrNonExistentID)
	assert.Empty(t, page.Items)
}

func TestQueryMentionsCursor(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	seed, err := pg_sack.QueryTrendingBeans(test_ctx, db.BeanFilters{
		ObservedFrom: testSearchFrom(),
	}, db.PageRequest{Limit: 20}, db.BEAN_COLUMNS_WITH_TREND)
	require.NoError(t, err)

	var first db.Page[db.Mention]
	var seed_id uuid.UUID
	for _, bean := range seed.Items {
		first, err = pg_sack.QueryMentions(test_ctx, bean.ID, db.MentionFilters{}, db.PageRequest{Limit: 2})
		require.NoError(t, err)
		if first.NextCursor != nil {
			seed_id = bean.ID
			break
		}
	}
	if first.NextCursor == nil {
		t.Skip("not enough mentions for a second page")
	}

	second, err := pg_sack.QueryMentions(test_ctx, seed_id, db.MentionFilters{}, db.PageRequest{Limit: 2, Cursor: first.NextCursor})
	require.NoError(t, err)
	require.NotEmpty(t, second.Items)
	assert.NotEqual(t, first.Items[0].URL, second.Items[0].URL)
}

func TestQueryClusters(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	page, err := pg_sack.QueryClusters(test_ctx, db.ClusterFilters{
		BeanFilters:  db.BeanFilters{CreatedFrom: testSearchFrom()},
		MinBeanCount: 2,
	}, db.PageRequest{Limit: 5})
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	for _, story := range page.Items {
		assert.NotEmpty(t, story.ID)
		assert.GreaterOrEqual(t, story.BeanCount, 2)
		assert.GreaterOrEqual(t, story.SourceCount, 1)
		assert.NotEmpty(t, story.TopArticles)
		assert.LessOrEqual(t, len(story.TopArticles), 3)
		assert.NotNil(t, story.Categories)
		assert.NotNil(t, story.Regions)
		assert.NotNil(t, story.Entities)
		assert.NotNil(t, story.Tags)
	}
	pp.Println("STORIES", page.Items)
}

func TestGetCluster(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	list, err := pg_sack.QueryClusters(test_ctx, db.ClusterFilters{
		BeanFilters:  db.BeanFilters{CreatedFrom: testSearchFrom()},
		MinBeanCount: 2,
	}, db.PageRequest{Limit: 1})
	require.NoError(t, err)
	require.NotEmpty(t, list.Items)
	want := list.Items[0]

	cluster, err := pg_sack.GetCluster(test_ctx, want.ID)
	require.NoError(t, err)
	assert.False(t, cluster.IsZero())
	assert.Equal(t, want.ID, cluster.ID)
	assert.GreaterOrEqual(t, cluster.BeanCount, 2)
	assert.NotEmpty(t, cluster.TopArticles)
	pp.Println("CLUSTER", cluster)
}

func TestGetClusterNotFound(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	cluster, err := pg_sack.GetCluster(test_ctx, uuid.UUID{})
	assert.ErrorIs(t, err, db.ErrNonExistentID)
	assert.True(t, cluster.IsZero())
}

func TestQueryStoryArticles(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	list, err := pg_sack.QueryClusters(test_ctx, db.ClusterFilters{
		BeanFilters:  db.BeanFilters{CreatedFrom: testSearchFrom()},
		MinBeanCount: 2,
	}, db.PageRequest{Limit: 1})
	require.NoError(t, err)
	require.NotEmpty(t, list.Items)
	cluster_id := list.Items[0].ID

	exists, err := pg_sack.ClusterExists(test_ctx, cluster_id)
	require.NoError(t, err)
	assert.True(t, exists)

	page, err := pg_sack.QueryBeans(test_ctx, db.BeanFilters{ClusterID: cluster_id}, db.PageRequest{Limit: 5}, db.BEAN_COLUMNS_WITHOUT_TREND)
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	for _, bean := range page.Items {
		assert.NotEqual(t, bean.ClusterID, uuid.Nil)
		assert.Equal(t, cluster_id, bean.ClusterID)
	}
	pp.Println("CLUSTER_ARTICLES", page.Items)
}

func TestQueryClustersCursor(t *testing.T) {
	pg_sack := setupTestDB()
	defer pg_sack.Close()

	filters := db.ClusterFilters{
		BeanFilters:  db.BeanFilters{CreatedFrom: testSearchFrom()},
		MinBeanCount: 2,
	}
	first, err := pg_sack.QueryClusters(test_ctx, filters, db.PageRequest{Limit: 2})
	require.NoError(t, err)
	require.NotEmpty(t, first.Items)
	if first.NextCursor == nil {
		t.Skip("not enough clusters for a second page")
	}

	second, err := pg_sack.QueryClusters(test_ctx, filters, db.PageRequest{Limit: 2, Cursor: first.NextCursor})
	require.NoError(t, err)
	require.NotEmpty(t, second.Items)
	assert.NotEqual(t, first.Items[0].ID, second.Items[0].ID)
}
