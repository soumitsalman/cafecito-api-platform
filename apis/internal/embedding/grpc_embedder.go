package embedding

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCEmbedder struct {
	conn   *grpc.ClientConn
	client EmbedClient
}

// NewGRPCEmbedder creates a new embedder that connects to Hugging Face TEI via gRPC
// baseURL should be in format "hostname:port" (e.g., "localhost:10000")
// model and apiKey parameters are ignored for TEI but kept for backward compatibility
func NewGRPCEmbedder(baseURL, apiKey, model string) *GRPCEmbedder {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return nil
	}

	conn, err := grpc.NewClient(
		baseURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error().Str("module", "EMBEDDER").Err(err).Msg("failed to connect to TEI gRPC server")
	}

	return &GRPCEmbedder{
		conn:   conn,
		client: NewEmbedClient(conn),
	}
}

func (e *GRPCEmbedder) embedInput(ctx context.Context, input string) []float32 {
	ctx, cancel := context.WithTimeout(ctx, _TIMEOUT)
	defer cancel()

	normalize := true
	resp, err := e.client.Embed(ctx, &EmbedRequest{
		Inputs:    input,
		Normalize: &normalize,
	})
	if err != nil {
		log.Error().Str("module", "EMBEDDER").Err(err).Msg("failed to embed input")
		return nil
	}

	if len(resp.Embeddings) == 0 {
		log.Error().Str("module", "EMBEDDER").Msg("empty embedding response from TEI")
		return nil
	}

	return resp.Embeddings
}

// EmbedQuery applies F2LLM's retrieval instruction before embedding the query.
func (e *GRPCEmbedder) EmbedQuery(ctx context.Context, query string) []float32 {
	return e.embedInput(ctx, retrievalQueryInput(query))
}

// EmbedDocuments embeds multiple documents by calling Embed for each
// (TEI's gRPC Embed service processes one input at a time; use EmbedStream for batching)
func (e *GRPCEmbedder) EmbedDocuments(ctx context.Context, docs []string) [][]float32 {
	if len(docs) == 0 {
		return nil
	}

	result := make([][]float32, 0, len(docs))
	for _, doc := range docs {
		embedding := e.embedInput(ctx, doc)
		if embedding != nil {
			result = append(result, embedding)
		}
	}

	return result
}

// Close closes the gRPC connection
func (e *GRPCEmbedder) Close() error {
	if e.conn != nil {
		return e.conn.Close()
	}
	return nil
}
