package embedding

import (
	"context"
	"time"
)

const _TIMEOUT = 10 * time.Minute
const _QUERY_PREFIX = "Instruct: Given a question, retrieve passages that can help answer the question.\nQuery: "

type Embedder interface {
	EmbedQuery(ctx context.Context, query string) []float32
	EmbedDocuments(ctx context.Context, docs []string) [][]float32
	Close() error
}

func retrievalQueryInput(query string) string {
	return _QUERY_PREFIX + query
}
