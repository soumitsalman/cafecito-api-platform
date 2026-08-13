package espressoapi_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/db"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustDigest(t *testing.T, m map[string]any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(m)
	require.NoError(t, err)
	return b
}

func TestNewDigestDocumentIncludesSipFields(t *testing.T) {
	id := uuid.New()
	created := time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC)
	sip := db.Sip{
		ID:      id,
		Kind:    db.SIP_KIND_EVENT,
		Created: created,
		Tags:    []string{"markets"},
		Digest: mustDigest(t, map[string]any{
			"briefing":      "Example Semiconductor cut its annual outlook.",
			"event_type":    "earnings_guidance",
			"stock_tickers": []any{"EXSC"},
		}),
	}

	doc := router.NewDigestDocumentForSip(&sip)
	require.NotNil(t, doc)
	assert.Equal(t, id, doc["id"])
	assert.Equal(t, created, doc["created_at"])
	assert.Equal(t, []string{"markets"}, doc["tags"])
	assert.Equal(t, "Example Semiconductor cut its annual outlook.", doc["briefing"])
	assert.Equal(t, "earnings_guidance", doc["event_type"])
	assert.NotContains(t, doc, "kind")
}

func TestNewDigestDocumentRejectsNonObject(t *testing.T) {
	doc := router.NewDigestDocumentForSip(&db.Sip{Digest: json.RawMessage(`[]`)})
	assert.Nil(t, doc)
}

func TestMaterializeDigest(t *testing.T) {
	sip := db.Sip{Digest: mustDigest(t, map[string]any{"briefing": "hello"})}
	fields, err := sip.MaterializeDigest()
	require.NoError(t, err)
	assert.Equal(t, "hello", fields["briefing"])

	_, err = (&db.Sip{Digest: json.RawMessage(`[]`)}).MaterializeDigest()
	assert.Error(t, err)
}

func TestCursorEncodeDecodeRoundTrip(t *testing.T) {
	id := uuid.New()
	created := time.Date(2026, 7, 29, 13, 0, 0, 0, time.UTC)
	c := &db.Cursor{Version: 1, ID: &id, Created: &created}

	encoded, err := c.Encode()
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := db.DecodeCursor(encoded)
	require.NoError(t, err)
	require.NotNil(t, decoded)
	require.NotNil(t, decoded.ID)
	assert.Equal(t, id, *decoded.ID)
	require.NotNil(t, decoded.Created)
	assert.True(t, created.Equal(*decoded.Created))
}

func TestCursorEncodeDecodeRelevance(t *testing.T) {
	id := uuid.New()
	distance := 0.42
	c := &db.Cursor{Version: 1, ID: &id, Distance: &distance}

	encoded, err := c.Encode()
	require.NoError(t, err)

	decoded, err := db.DecodeCursor(encoded)
	require.NoError(t, err)
	require.NotNil(t, decoded.ID)
	assert.Equal(t, id, *decoded.ID)
	require.NotNil(t, decoded.Distance)
	assert.InDelta(t, distance, *decoded.Distance, 1e-9)
}

func TestDecodeCursorRejectsGarbage(t *testing.T) {
	_, err := db.DecodeCursor("not-valid-base64-or-json!!")
	assert.ErrorIs(t, err, db.ErrInvalidCursor)

	decoded, err := db.DecodeCursor("")
	require.NoError(t, err)
	assert.Nil(t, decoded)
}

func TestNewDigestDocuments(t *testing.T) {
	sips := []db.Sip{
		{
			ID:      uuid.New(),
			Kind:    db.SIP_KIND_SIGNAL,
			Created: time.Now().UTC(),
			Digest:  mustDigest(t, map[string]any{"thesis": "rates stay high"}),
		},
	}
	docs := router.NewDigestDocuments(sips)
	require.Len(t, docs, 1)
	assert.Equal(t, "rates stay high", docs[0]["thesis"])
	assert.Equal(t, sips[0].ID, docs[0]["id"])
}

func TestNewSourceDocumentUsesPublicFieldNames(t *testing.T) {
	id := uuid.New()
	doc := router.NewSourceDocument(&db.Source{ID: id, BaseURL: "https://example.com"})
	raw, err := json.Marshal(doc)
	require.NoError(t, err)

	var fields map[string]any
	require.NoError(t, json.Unmarshal(raw, &fields))
	assert.Equal(t, id.String(), fields["id"])
	assert.Equal(t, "https://example.com", fields["url"])
	assert.Contains(t, fields, "domain")
	assert.Contains(t, fields, "name")
	assert.Contains(t, fields, "description")
	assert.Contains(t, fields, "favicon_url")
	assert.Contains(t, fields, "rss_feed_url")
	assert.NotContains(t, fields, "base_url")
}
