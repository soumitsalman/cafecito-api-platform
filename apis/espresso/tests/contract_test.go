package espressoapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/db"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/router"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toon-format/toon-go"
	"gopkg.in/yaml.v3"
)

const testBackendKey = "test-contract-key"

func newHermeticRouter(t *testing.T, api_keys map[string]string) *httptest.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := router.NewRouter(nil, nil, api_keys)
	srv := httptest.NewServer(engine)
	t.Cleanup(srv.Close)
	return srv
}

func TestContractHealthUnauthenticated(t *testing.T) {
	srv := newHermeticRouter(t, map[string]string{"X-API-KEY": testBackendKey})
	status, body := routerGET(t, srv.URL, ROUTE_HEALTH, nil, "")
	requireStatus(t, http.StatusOK, status, body)
	var response map[string]string
	require.NoError(t, json.Unmarshal(body, &response))
	assert.Equal(t, "alive", response["status"])
}

func TestContractRESTRequiresAPIKey(t *testing.T) {
	srv := newHermeticRouter(t, map[string]string{"X-API-KEY": testBackendKey})
	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, nil, "")
	requireStatus(t, http.StatusUnauthorized, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_UNAUTHORIZED)
}

func TestContractRESTAcceptsAPIKey(t *testing.T) {
	srv := newHermeticRouter(t, map[string]string{"X-API-KEY": testBackendKey})
	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, url.Values{"cursor": {"not-a-valid-cursor"}}, testBackendKey)
	requireStatus(t, http.StatusBadRequest, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
}

func TestContractInvalidCursor(t *testing.T) {
	srv := newHermeticRouter(t, nil)
	status, body := routerGET(t, srv.URL, ROUTE_EVENTS, url.Values{"cursor": {"not-a-valid-cursor"}}, "")
	requireStatus(t, http.StatusBadRequest, status, body)
	assertExpectedAPIError(t, body, shared.API_ERROR_INVALID_REQUEST)
}

func TestContractEmptyCollectionEnvelope(t *testing.T) {
	page := router.PageResponse[map[string]any]{
		Data:       []map[string]any{},
		Pagination: router.NewPagination(20, 0, nil),
	}
	raw, err := json.Marshal(page)
	require.NoError(t, err)
	var env map[string]any
	require.NoError(t, json.Unmarshal(raw, &env))
	data, ok := env["data"].([]any)
	require.True(t, ok)
	assert.Empty(t, data)
	pagination, ok := env["pagination"].(map[string]any)
	require.True(t, ok)
	assert.NotContains(t, pagination, "cursor")
	assert.Nil(t, pagination["next_cursor"])
	assert.Equal(t, float64(20), pagination["limit"])
	assert.Equal(t, float64(0), pagination["num_results"])
}

func TestContractPaginationNextCursorAllFormats(t *testing.T) {
	next := "opaque-next"
	first := router.NewPagination(2, 2, &next)
	second := router.NewPagination(2, 1, nil)

	raw_json, err := json.Marshal(first)
	require.NoError(t, err)
	var json_first map[string]any
	require.NoError(t, json.Unmarshal(raw_json, &json_first))
	assert.NotContains(t, json_first, "cursor")
	assert.Equal(t, next, json_first["next_cursor"])

	raw_json, err = json.Marshal(second)
	require.NoError(t, err)
	var json_second map[string]any
	require.NoError(t, json.Unmarshal(raw_json, &json_second))
	assert.NotContains(t, json_second, "cursor")
	assert.Nil(t, json_second["next_cursor"])

	yaml_first, err := yaml.Marshal(first)
	require.NoError(t, err)
	var yaml_first_obj map[string]any
	require.NoError(t, yaml.Unmarshal(yaml_first, &yaml_first_obj))
	assert.NotContains(t, yaml_first_obj, "cursor")
	assert.Equal(t, next, yaml_first_obj["next_cursor"])

	yaml_second, err := yaml.Marshal(second)
	require.NoError(t, err)
	var yaml_second_obj map[string]any
	require.NoError(t, yaml.Unmarshal(yaml_second, &yaml_second_obj))
	assert.NotContains(t, yaml_second_obj, "cursor")
	assert.Nil(t, yaml_second_obj["next_cursor"])

	toon_first, err := toon.MarshalString(first)
	require.NoError(t, err)
	assert.NotRegexp(t, `(?m)^cursor:`, toon_first)
	assert.Contains(t, toon_first, "next_cursor")
	toon_second, err := toon.MarshalString(second)
	require.NoError(t, err)
	assert.NotRegexp(t, `(?m)^cursor:`, toon_second)
	assert.Contains(t, toon_second, "next_cursor")
}

func TestContractTypedErrorEnvelope(t *testing.T) {
	env := router.ErrorResponse{Error: router.APIError{Code: shared.API_ERROR_NOT_FOUND, Message: "EVENT with this ID not found"}}
	raw, err := json.Marshal(env)
	require.NoError(t, err)
	assertExpectedAPIError(t, raw, shared.API_ERROR_NOT_FOUND)
}

func TestContractStableEventCore(t *testing.T) {
	id := uuid.New()
	sip := db.Sip{
		ID:      id,
		Kind:    db.SIP_KIND_EVENT,
		Created: testSearchFrom(),
		Tags:    []string{"markets"},
		Digest:  mustDigest(t, map[string]any{"briefing": "Outlook cut.", "event_type": "earnings_guidance", "extra_pipeline": "ok"}),
	}
	doc := router.NewDigestDocumentForSip(&sip)
	require.Equal(t, id, doc["id"])
	require.Equal(t, db.SIP_KIND_EVENT, doc["kind"])
	require.Contains(t, doc, "created_at")
	require.Equal(t, []string{"markets"}, doc["tags"])
	assert.Equal(t, "Outlook cut.", doc["summary"])
	assert.NotContains(t, doc, "briefing")
	assert.Equal(t, "ok", doc["extra_pipeline"])
}
