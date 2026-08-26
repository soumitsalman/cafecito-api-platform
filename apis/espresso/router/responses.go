package router

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/db"
	datautils "github.com/soumitsalman/data-utils"
)

// Pagination is the pagination block returned by every collection response.
// Send next_cursor unchanged as the next request's cursor. The response does not echo a cursor field.
type Pagination struct {
	Limit      int     `json:"limit" yaml:"limit" toon:"limit"`
	NumResults int     `json:"num_results" yaml:"num_results" toon:"num_results"`
	NextCursor *string `json:"next_cursor" yaml:"next_cursor" toon:"next_cursor" extensions:"x-nullable"`
}

// NewPagination builds the public pagination object for JSON, YAML, and TOON.
func NewPagination(limit int, num_results int, next_cursor *string) Pagination {
	return Pagination{
		Limit:      limit,
		NumResults: num_results,
		NextCursor: next_cursor,
	}
}

// ResponseMeta carries freshness metadata for a collection response.
type ResponseMeta struct {
	AsOf time.Time `json:"as_of" yaml:"as_of" toon:"as_of"`
}

// PageResponse is the canonical envelope for list endpoints.
type PageResponse[T any] struct {
	Data       []T          `json:"data" yaml:"data" toon:"data"`
	Pagination Pagination   `json:"pagination" yaml:"pagination" toon:"pagination"`
	Meta       ResponseMeta `json:"meta" yaml:"meta" toon:"meta"`
}

// ItemResponse is the canonical envelope for single-resource endpoints.
type ItemResponse[T any] struct {
	Data T `json:"data" yaml:"data" toon:"data"`
}

// APIError is the typed object inside the public error envelope.
type APIError struct {
	Code    string `json:"code" yaml:"code" toon:"code"`
	Message string `json:"message" yaml:"message" toon:"message"`
}

// ErrorResponse is the canonical envelope for error responses: {"error":{"code","message"}}.
type ErrorResponse struct {
	Error APIError `json:"error" yaml:"error" toon:"error"`
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

// Stable Event/Signal fields present on every public record: id, kind, created_at, tags.
// Conditional fields include summary, source, links, and counts.
// All other keys are extension fields; clients must ignore unknown extension fields.

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
	ID          uuid.UUID `json:"id" yaml:"id" toon:"id" swaggertype:"string" format:"uuid"`
	BaseURL     string    `json:"url" yaml:"url" toon:"url"`
	DomainName  string    `json:"domain" yaml:"domain" toon:"domain"`
	SiteName    *string   `json:"name" yaml:"name" toon:"name"`
	Description *string   `json:"description,omitempty" yaml:"description,omitempty" toon:"description,omitempty"`
	Favicon     *string   `json:"favicon_url,omitempty" yaml:"favicon_url,omitempty" toon:"favicon_url,omitempty"`
	RSSFeed     *string   `json:"rss_feed_url,omitempty" yaml:"rss_feed_url,omitempty" toon:"rss_feed_url,omitempty"`
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
	Evidence string `json:"evidence,omitempty" yaml:"evidence,omitempty" toon:"evidence,omitempty"`
	Actions  string `json:"actions,omitempty" yaml:"actions,omitempty" toon:"actions,omitempty"`
	Events   string `json:"events,omitempty" yaml:"events,omitempty" toon:"events,omitempty"`
	Signals  string `json:"signals,omitempty" yaml:"signals,omitempty" toon:"signals,omitempty"`
}

type Counts struct {
	Evidence *int64 `json:"evidence,omitempty" yaml:"evidence,omitempty" toon:"evidence,omitempty"`
	Actions  *int64 `json:"actions,omitempty" yaml:"actions,omitempty" toon:"actions,omitempty"`
	Events   *int64 `json:"events,omitempty" yaml:"events,omitempty" toon:"events,omitempty"`
	Signals  *int64 `json:"signals,omitempty" yaml:"signals,omitempty" toon:"signals,omitempty"`
}

// SipEvidenceItem is the explicit bare-list response item for R03.
type EventEvidence struct {
	ID       uuid.UUID  `json:"id" yaml:"id" toon:"id"`
	Kind     string     `json:"kind" yaml:"kind" toon:"kind"`
	Created  time.Time  `json:"created_at" yaml:"created_at" toon:"created_at"`
	Tags     []string   `json:"tags" yaml:"tags" toon:"tags"`
	SourceID *uuid.UUID `json:"source_id" yaml:"source_id" toon:"source_id"`
	URL      string     `json:"url" yaml:"url" toon:"url"`
	BaseURL  string     `json:"base_url" yaml:"base_url" toon:"base_url"`
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
	Data       []db.Tag     `json:"data" binding:"required"`
	Pagination Pagination   `json:"pagination" binding:"required"`
	Meta       ResponseMeta `json:"meta"`
}

type EventEvidenceCollectionResponse struct {
	Data       []EventEvidence `json:"data" binding:"required"`
	Pagination Pagination      `json:"pagination" binding:"required"`
	Meta       ResponseMeta    `json:"meta" binding:"required"`
}

// sipDocument is the shared OpenAPI field set for Event and Signal records.
// Kind is declared on EventDocument and SignalDocument so each keeps its own enum.
type sipDocument struct {
	ID              uuid.UUID `json:"id" yaml:"id" toon:"id" binding:"required" swaggertype:"string" format:"uuid"`
	CreatedAt       time.Time `json:"created_at" yaml:"created_at" toon:"created_at" binding:"required"`
	Tags            []string  `json:"tags" yaml:"tags" toon:"tags" binding:"required"`
	KeyPoints       []string  `json:"key_points,omitempty" yaml:"key_points,omitempty" toon:"key_points,omitempty"`
	Drivers         []string  `json:"drivers,omitempty" yaml:"drivers,omitempty" toon:"drivers,omitempty"`
	Impacts         []string  `json:"impacts,omitempty" yaml:"impacts,omitempty" toon:"impacts,omitempty"`
	ImpactedDomains []string  `json:"impacted_domains,omitempty" yaml:"impacted_domains,omitempty" toon:"impacted_domains,omitempty"`
	ImpactLevel     string    `json:"impact_level,omitempty" yaml:"impact_level,omitempty" toon:"impact_level,omitempty"`
	FutureOutlook   string    `json:"future_outlook,omitempty" yaml:"future_outlook,omitempty" toon:"future_outlook,omitempty"`
	Summary         string    `json:"summary,omitempty" yaml:"summary,omitempty" toon:"summary,omitempty"`
	Categories      []string  `json:"categories,omitempty" yaml:"categories,omitempty" toon:"categories,omitempty"`
	Companies       []string  `json:"companies,omitempty" yaml:"companies,omitempty" toon:"companies,omitempty"`
	Regions         []string  `json:"regions,omitempty" yaml:"regions,omitempty" toon:"regions,omitempty"`
	People          []string  `json:"people,omitempty" yaml:"people,omitempty" toon:"people,omitempty"`
	Products        []string  `json:"products,omitempty" yaml:"products,omitempty" toon:"products,omitempty"`
}

type SignalDocument struct {
	sipDocument
	Kind       string `json:"kind" yaml:"kind" toon:"kind" binding:"required" enums:"signal"`
	Confidence string `json:"confidence,omitempty" yaml:"confidence,omitempty" toon:"confidence,omitempty"`
}

type SignalCollectionResponse struct {
	Data       []SignalDocument `json:"data" binding:"required"`
	Pagination Pagination       `json:"pagination" binding:"required"`
	Meta       ResponseMeta     `json:"meta" binding:"required"`
}

type SignalDetail struct {
	SignalDocument
	Links  *Links  `json:"links,omitempty" yaml:"links,omitempty" toon:"links,omitempty"`
	Counts *Counts `json:"counts,omitempty" yaml:"counts,omitempty" toon:"counts,omitempty"`
}

type SignalDetailResponse struct {
	Data SignalDetail `json:"data" binding:"required"`
}

// EventDocument is the public Event schema for OpenAPI generation.
// Runtime Event records remain extensible maps; clients must ignore unknown extension fields.
type EventDocument struct {
	sipDocument
	Kind         string          `json:"kind" yaml:"kind" toon:"kind" binding:"required" enums:"event"`
	EventType    string          `json:"event_type,omitempty" yaml:"event_type,omitempty" toon:"event_type,omitempty"`
	MacroContext string          `json:"macro_context,omitempty" yaml:"macro_context,omitempty" toon:"macro_context,omitempty"`
	Source       *SourceDocument `json:"source,omitempty" yaml:"source,omitempty" toon:"source,omitempty"`
}

type EventCollectionResponse struct {
	Data       []EventDocument `json:"data" binding:"required"`
	Pagination Pagination      `json:"pagination" binding:"required"`
	Meta       ResponseMeta    `json:"meta" binding:"required"`
}

type EventDetail struct {
	EventDocument
	Links  *Links  `json:"links,omitempty" yaml:"links,omitempty" toon:"links,omitempty"`
	Counts *Counts `json:"counts,omitempty" yaml:"counts,omitempty" toon:"counts,omitempty"`
}

type EventDetailResponse struct {
	Data EventDetail `json:"data" binding:"required"`
}

type SourceItemResponse struct {
	Data SourceDocument `json:"data" binding:"required"`
}

type DiscoveryValueCollectionResponse struct {
	Data       []db.Tag     `json:"data" binding:"required"`
	Pagination Pagination   `json:"pagination" binding:"required"`
	Meta       ResponseMeta `json:"meta" binding:"required"`
}
