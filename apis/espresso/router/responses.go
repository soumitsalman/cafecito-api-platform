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
	Limit      int     `json:"limit" toon:"limit" binding:"required"`
	NumResults int     `json:"num_results" toon:"num_results" binding:"required"`
	Cursor     *string `json:"page" toon:"page" binding:"required"`
	NextCursor *string `json:"next_page" toon:"next_page" binding:"required"`
}

// ResponseMeta carries freshness metadata for a collection response.
type ResponseMeta struct {
	AsOf time.Time `json:"as_of" toon:"as_of" binding:"required"`
}

// PageResponse is the canonical envelope for list endpoints.
type PageResponse[T any] struct {
	Data       []T          `json:"data" toon:"data"`
	Pagination Pagination   `json:"pagination" toon:"pagination"`
	Meta       ResponseMeta `json:"meta" toon:"meta"`
}

// ItemResponse is the canonical envelope for single-resource endpoints.
type ItemResponse[T any] struct {
	Data T `json:"data" toon:"data"`
}

// APIError is the stable public error shape.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("Error Code=%s, Message=%s", e.Code, e.Message)
}

// ErrorResponse is the canonical envelope for error responses.
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// DigestDocument preserves arbitrary upstream digest members without imposing a closed Event or Signal response schema.
type DigestDocument map[string]any

func NewDigestDocumentForSip(sip *db.Sip) DigestDocument {
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
	if briefing, ok := doc["briefing"]; ok {
		doc["summary"] = briefing
		delete(doc, "briefing")
	}
	return doc
}

func NewDigestDocumentForExtendedSip(sip *db.ExtendedSip) DigestDocument {
	doc := NewDigestDocumentForSip(&sip.Sip)
	doc["url"] = nil
	doc["base_url"] = nil
	doc["source_id"] = nil

	// override default values for non-signals
	if sip.Kind != db.SIP_KIND_SIGNAL {
		if sip.URL.Valid {
			doc["url"] = sip.URL.String
		}
		if sip.BaseURL.Valid {
			doc["base_url"] = sip.BaseURL.String
		}
		if sip.SourceID != uuid.Nil {
			doc["source_id"] = sip.SourceID
			doc["source"] = NewSourceDocument(sip.GetSource())
		}
	}
	return doc
}

func NewDigestDocuments(sips []db.Sip) []DigestDocument {
	return datautils.Transform(sips, func(sip *db.Sip) DigestDocument {
		return NewDigestDocumentForSip(sip)
	})
}

func (digest DigestDocument) addEventDetails(counts db.RelationCounts) DigestDocument {
	digest["links"] = Links{
		Evidence: fmt.Sprintf("/events/%s/evidence", digest["id"]),
		Signals:  fmt.Sprintf("/events/%s/signals", digest["id"]),
	}
	digest["counts"] = Counts{
		Evidence: &counts.SameAs,
		Signals:  &counts.DerivedTo,
	}
	return digest
}

func (digest DigestDocument) addSignalDetails(counts db.RelationCounts) DigestDocument {
	digest["links"] = Links{
		Events: fmt.Sprintf("/signals/%s/events", digest["id"]),
	}
	digest["counts"] = Counts{
		Events: &counts.DerivedFrom,
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

// SourceDocument is the stable public Source response shape. Optional source
// metadata is represented explicitly as null when it is unavailable.
type SourceDocument struct {
	ID          uuid.UUID `json:"id" toon:"id" swaggertype:"string" format:"uuid"`
	BaseURL     string    `json:"url" toon:"url"`
	DomainName  string    `json:"domain" toon:"domain"`
	SiteName    *string   `json:"name" toon:"name"`
	Description *string   `json:"description,omitempty" toon:"description,omitempty"`
	Favicon     *string   `json:"favicon_url,omitempty" toon:"favicon_url,omitempty"`
	RSSFeed     *string   `json:"rss_feed_url,omitempty" toon:"rss_feed_url,omitempty"`
}

func NewSourceDocument(source *db.Source) *SourceDocument {
	doc := &SourceDocument{
		ID:         source.ID,
		BaseURL:    source.BaseURL,
		DomainName: source.DomainName,
	}
	if source.SiteName.Valid {
		doc.SiteName = &source.SiteName.String
	}
	if source.Description.Valid {
		doc.Description = &source.Description.String
	}
	if source.Favicon.Valid {
		doc.Favicon = &source.Favicon.String
	}
	if source.RSSFeed.Valid {
		doc.RSSFeed = &source.RSSFeed.String
	}
	return doc
}

func NewSourceDocuments(sources []db.Source) []SourceDocument {
	return datautils.Transform(sources, func(source *db.Source) SourceDocument {
		return *NewSourceDocument(source)
	})
}

type Links struct {
	Evidence string `json:"evidence,omitempty" toon:"evidence,omitempty"`
	Actions  string `json:"actions,omitempty" toon:"actions,omitempty"`
	Events   string `json:"events,omitempty" toon:"events,omitempty"`
	Signals  string `json:"signals,omitempty" toon:"signals,omitempty"`
}

type Counts struct {
	Evidence *int64 `json:"evidence,omitempty" toon:"evidence,omitempty"`
	Actions  *int64 `json:"actions,omitempty" toon:"actions,omitempty"`
	Events   *int64 `json:"events,omitempty" toon:"events,omitempty"`
	Signals  *int64 `json:"signals,omitempty" toon:"signals,omitempty"`
}

// SipEvidenceItem is the explicit bare-list response item for R03.
type EventEvidence struct {
	ID       uuid.UUID  `json:"id" toon:"id"`
	Kind     string     `json:"kind" toon:"kind"`
	Created  time.Time  `json:"created_at" toon:"created_at"`
	Tags     []string   `json:"tags" toon:"tags"`
	SourceID *uuid.UUID `json:"source_id" toon:"source_id"`
	URL      string     `json:"url" toon:"url"`
	BaseURL  string     `json:"base_url" toon:"base_url"`
}

func NewEventEvidence(sip *db.Sip) EventEvidence {
	doc := EventEvidence{
		ID:      sip.ID,
		Kind:    sip.Kind,
		Created: sip.Created,
		Tags:    sip.Tags,
		URL:     sip.URL.String,
		BaseURL: sip.BaseURL.String,
	}
	if sip.SourceID != uuid.Nil {
		doc.SourceID = &sip.SourceID
	} else {
		doc.SourceID = nil
	}
	return doc
}

// The concrete envelope types below exist so swag can generate named schemas for the
// OpenAPI spec. swag does not understand generic type parameters, so each collection
// and detail response needs a concrete wrapper. The runtime still uses the generic
// CollectionResponse[T] and DetailResponse[T] helpers; these are schema-only.
type SipCollectionResponse struct {
	Data       []DigestDocument `json:"data"`
	Pagination Pagination       `json:"pagination" binding:"required"`
	Meta       ResponseMeta     `json:"meta" binding:"required"`
}

type SourceCollectionResponse struct {
	Data       []SourceDocument `json:"data" binding:"required"`
	Pagination Pagination       `json:"pagination" binding:"required"`
	Meta       ResponseMeta     `json:"meta" binding:"required"`
}

type StringCollectionResponse struct {
	Data       []string     `json:"data"`
	Pagination Pagination   `json:"pagination"`
	Meta       ResponseMeta `json:"meta"`
}

type TagValueCollectionResponse struct {
	Data       []db.TagValue `json:"data" binding:"required"`
	Pagination Pagination    `json:"pagination" binding:"required"`
	Meta       ResponseMeta  `json:"meta"`
}

type EventEvidenceCollectionResponse struct {
	Data       []EventEvidence `json:"data" binding:"required"`
	Pagination Pagination      `json:"pagination" binding:"required"`
	Meta       ResponseMeta    `json:"meta" binding:"required"`
}

// EventDocument is the flattened, extensible public Event collection/detail payload.
type EventDocument map[string]any

// SignalDocument is the flattened, extensible public Signal collection/detail payload.
type SignalDocument map[string]any

type EventCollectionResponse struct {
	Data       []EventDocument `json:"data" binding:"required"`
	Pagination Pagination      `json:"pagination" binding:"required"`
	Meta       ResponseMeta    `json:"meta" binding:"required"`
}

type EventDetailResponse struct {
	Data EventDocument `json:"data" binding:"required"`
}

type SignalCollectionResponse struct {
	Data       []SignalDocument `json:"data" binding:"required"`
	Pagination Pagination       `json:"pagination" binding:"required"`
	Meta       ResponseMeta     `json:"meta" binding:"required"`
}

type SignalDetailResponse struct {
	Data SignalDocument `json:"data" binding:"required"`
}

type SipItemResponse struct {
	Data DigestDocument `json:"data" binding:"required"`
}

type SourceItemResponse struct {
	Data SourceDocument `json:"data" binding:"required"`
}

type DiscoveryValueCollectionResponse struct {
	Data       []db.TagValue `json:"data" binding:"required"`
	Pagination Pagination    `json:"pagination" binding:"required"`
	Meta       ResponseMeta  `json:"meta" binding:"required"`
}
