package internal_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/soumitsalman/cafecito-api-platform/apis/internal/embedding"
)

const retrievalQueryPrefix = "Instruct: Given a question, retrieve passages that can help answer the question.\nQuery: "

type httpEmbeddingRequest struct {
	Input          []string `json:"input"`
	Model          string   `json:"model,omitempty"`
	EncodingFormat string   `json:"encoding_format"`
}

func retrievalQueryInput(query string) string {
	return retrievalQueryPrefix + query
}

func TestHTTPEmbedderPrefixesQueriesOnly(t *testing.T) {
	t.Helper()

	var requests []httpEmbeddingRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/embeddings" {
			t.Fatalf("unexpected endpoint: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("unexpected authorization header: %q", got)
		}

		var request httpEmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		requests = append(requests, request)

		w.Header().Set("Content-Type", "application/json")
		if len(request.Input) == 1 {
			_, _ = w.Write([]byte(`{"data":[{"index":0,"embedding":[1,0]}]}`))
			return
		}
		_, _ = w.Write([]byte(`{"data":[{"index":1,"embedding":[0,1]},{"index":0,"embedding":[1,0]}]}`))
	}))
	defer server.Close()

	embedder := embedding.NewHTTPEmbedder(server.URL, "secret", "F2LLM-v2-80M")
	defer embedder.Close()

	queryEmbedding := embedder.EmbedQuery(context.Background(), "vector databases")
	if !reflect.DeepEqual(queryEmbedding, []float32{1, 0}) {
		t.Fatalf("unexpected query embedding: %v", queryEmbedding)
	}

	documentEmbeddings := embedder.EmbedDocuments(context.Background(), []string{"first", "second"})
	if !reflect.DeepEqual(documentEmbeddings, [][]float32{{1, 0}, {0, 1}}) {
		t.Fatalf("unexpected document embeddings: %v", documentEmbeddings)
	}

	if len(requests) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(requests))
	}
	if got, want := requests[0].Input, []string{retrievalQueryInput("vector databases")}; !reflect.DeepEqual(got, want) {
		t.Fatalf("query input mismatch: got %q, want %q", got, want)
	}
	if got, want := requests[1].Input, []string{"first", "second"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("document input mismatch: got %q, want %q", got, want)
	}
	if requests[0].Model != "F2LLM-v2-80M" || requests[0].EncodingFormat != "float" {
		t.Fatalf("unexpected request options: %+v", requests[0])
	}
}
