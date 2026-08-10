package router

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/db"
	datautils "github.com/soumitsalman/data-utils"
)

// Pagination is the pagination block returned by every collection response.
type Pagination struct {
	Limit      int     `json:"limit"`
	NextCursor *string `json:"next_cursor"`
}

// ResponseMeta carries freshness metadata for a collection response.
type ResponseMeta struct {
	AsOf time.Time `json:"as_of"`
}

// CollectionResponse is the canonical envelope for list endpoints.
type PageResponse[T any] struct {
	Data       []T          `json:"data"`
	Pagination Pagination   `json:"pagination"`
	Meta       ResponseMeta `json:"meta"`
}

// DetailResponse is the canonical envelope for single-resource endpoints.
type ItemResponse[T any] struct {
	Data T `json:"data"`
}

// APIError is the stable public error shape.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("Error Code=%s, Message=%s", e.Code, e.Message)
}

// DigestDocument preserves arbitrary upstream digest members without imposing a closed Event or Signal response schema.
type DigestDocument map[string]any

func NewDigestDocument(sip *db.Sip) DigestDocument {
	fields, err := sip.MaterializeDigest()
	if err != nil {
		return nil
	}
	doc := DigestDocument(fields)
	doc["id"] = sip.ID
	doc["created_at"] = sip.Created
	doc["kind"] = sip.Kind
	if _, ok := doc["tags"]; !ok {
		doc["tags"] = sip.Tags
	}
	return doc
}

func NewDigestDocuments(sips []db.Sip) []DigestDocument {
	return datautils.Transform(sips, func(sip *db.Sip) DigestDocument {
		return NewDigestDocument(sip)
	})
}

func (digest DigestDocument) addEventDetails(counts db.RelationCounts) DigestDocument {
	digest["links"] = Links{
		Evidence: fmt.Sprintf("/events/%s/evidence", digest["id"]),
		Signals:  fmt.Sprintf("/events/%s/signals", digest["id"]),
		// Action:   fmt.Sprintf("/events/%s/actions", digest["id"]), // TODO: add actions link
	}
	digest["counts"] = Counts{
		Coverage: counts.SameAs,
		Actions:  counts.DerivedFrom,
		Signals:  counts.DerivedTo,
	}
	return digest
}

func (digest DigestDocument) addSignalDetails(counts db.RelationCounts) DigestDocument {
	digest["links"] = Links{
		Events: fmt.Sprintf("/signals/%s/events", digest["id"]),
	}
	digest["counts"] = Counts{
		Events: counts.DerivedFrom,
	}
	return digest
}

// encodeNextCursor wraps db.Cursor.Encode and returns nil string when cursor is nil.
func encodeNextCursor(c *db.Cursor) *string {
	if c == nil {
		return nil
	}
	s, err := c.Encode()
	if err != nil {
		return nil
	}
	return &s
}

/***********************************************************
THESE TYPES ARE USED PRIMARILY TO GENERATE THE OPENAPI SPEC
***********************************************************/

// SourceDocument preserves arbitrary upstream source members without imposing
// a closed Source response schema.
type SourceDocument map[string]any

type Links struct {
	Evidence string `json:"evidence,omitempty"`
	Actions  string `json:"actions,omitempty"`
	Events   string `json:"events,omitempty"`
	Signals  string `json:"signals,omitempty"`
}

type Counts struct {
	Coverage int64 `json:"coverage,omitzero"`
	Actions  int64 `json:"actions,omitzero"`
	Events   int64 `json:"events,omitzero"`
	Signals  int64 `json:"signals,omitzero"`
}

// SipEvidenceItem is the explicit bare-list response item for R03.
type EventEvidence struct {
	EventID  uuid.UUID `json:"event_id,omitzero"`
	Created  time.Time `json:"created,omitzero"`
	SourceID uuid.UUID `json:"source_id,omitzero"`
	URL      string    `json:"url,omitempty"`
	BaseURL  string    `json:"base_url,omitempty"`
}

func NewEventEvidence(sip *db.Sip) EventEvidence {
	return EventEvidence{
		EventID:  sip.ID,
		Created:  sip.Created,
		SourceID: sip.SourceID,
		URL:      sip.URL.String,
		BaseURL:  sip.BaseURL.String,
	}
}

// The concrete envelope types below exist so swag can generate named schemas for the
// OpenAPI spec. swag does not understand generic type parameters, so each collection
// and detail response needs a concrete wrapper. The runtime still uses the generic
// CollectionResponse[T] and DetailResponse[T] helpers; these are schema-only.
type SipCollectionResponse struct {
	Data       []DigestDocument `json:"data"`
	Pagination Pagination       `json:"pagination"`
	Meta       ResponseMeta     `json:"meta"`
}

type SourceCollectionResponse struct {
	Data       []SourceDocument `json:"data"`
	Pagination Pagination       `json:"pagination"`
	Meta       ResponseMeta     `json:"meta"`
}

type StringCollectionResponse struct {
	Data       []string     `json:"data"`
	Pagination Pagination   `json:"pagination"`
	Meta       ResponseMeta `json:"meta"`
}
