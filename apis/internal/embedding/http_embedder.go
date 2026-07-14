package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type HTTPEmbedder struct {
	endpoint string
	apiKey   string
	model    string
	client   *http.Client
}

type httpEmbeddingRequest struct {
	Input          []string `json:"input"`
	Model          string   `json:"model,omitempty"`
	EncodingFormat string   `json:"encoding_format"`
}

type httpEmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
}

// NewHTTPEmbedder creates an embedder backed by llama-server's
// OpenAI-compatible /v1/embeddings HTTP endpoint.
func NewHTTPEmbedder(baseURL, apiKey, model string) *HTTPEmbedder {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return nil
	}

	endpoint := strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(endpoint, "/v1/embeddings") {
		endpoint += "/v1/embeddings"
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()

	return &HTTPEmbedder{
		endpoint: endpoint,
		apiKey:   apiKey,
		model:    model,
		client: &http.Client{
			Timeout:   _TIMEOUT,
			Transport: transport,
		},
	}
}

func (e *HTTPEmbedder) embed(ctx context.Context, inputs []string) [][]float32 {
	if len(inputs) == 0 {
		return nil
	}

	payload, err := json.Marshal(httpEmbeddingRequest{
		Input:          inputs,
		Model:          e.model,
		EncodingFormat: "float",
	})
	if err != nil {
		log.Error().Str("module", "EMBEDDER").Err(err).Msg("failed to encode llama embedding request")
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.endpoint, bytes.NewReader(payload))
	if err != nil {
		log.Error().Str("module", "EMBEDDER").Err(err).Msg("failed to create llama embedding request")
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	if e.apiKey != "" {
		if strings.HasPrefix(strings.ToLower(e.apiKey), "bearer ") {
			req.Header.Set("Authorization", e.apiKey)
		} else {
			req.Header.Set("Authorization", "Bearer "+e.apiKey)
		}
	}

	resp, err := e.client.Do(req)
	if err != nil {
		log.Error().Str("module", "EMBEDDER").Err(err).Msg("llama embedding request failed")
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		log.Error().
			Str("module", "EMBEDDER").
			Int("status", resp.StatusCode).
			Str("response", strings.TrimSpace(string(body))).
			Msg("llama embedding request returned an error")
		return nil
	}

	var decoded httpEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		log.Error().Str("module", "EMBEDDER").Err(err).Msg("failed to decode llama embedding response")
		return nil
	}
	if len(decoded.Data) != len(inputs) {
		log.Error().
			Str("module", "EMBEDDER").
			Err(fmt.Errorf("received %d embeddings for %d inputs", len(decoded.Data), len(inputs))).
			Msg("invalid llama embedding response")
		return nil
	}

	result := make([][]float32, len(inputs))
	for _, item := range decoded.Data {
		if item.Index < 0 || item.Index >= len(result) || len(item.Embedding) == 0 || result[item.Index] != nil {
			log.Error().Str("module", "EMBEDDER").Msg("invalid embedding item returned by llama-server")
			return nil
		}
		result[item.Index] = item.Embedding
	}

	return result
}

// EmbedQuery applies F2LLM's retrieval instruction before embedding the query.
func (e *HTTPEmbedder) EmbedQuery(ctx context.Context, query string) []float32 {
	embeddings := e.embed(ctx, []string{retrievalQueryInput(query)})
	if len(embeddings) == 0 {
		return nil
	}
	return embeddings[0]
}

// EmbedDocuments embeds documents without the retrieval query instruction.
func (e *HTTPEmbedder) EmbedDocuments(ctx context.Context, docs []string) [][]float32 {
	return e.embed(ctx, docs)
}

func (e *HTTPEmbedder) Close() error {
	if transport, ok := e.client.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
	return nil
}
