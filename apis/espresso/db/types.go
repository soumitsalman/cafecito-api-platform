package db

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SipKind is the stored kind of a sip record.
const (
	SIP_KIND_ACTION = "action"
	SIP_KIND_EVENT  = "event"
	SIP_KIND_SIGNAL = "signal"
)

// Sip is the storage-shaped unit of information in Espresso.
//
// The database layer returns storage-shaped records and leaves public JSON naming to
// the router. Digest is kept as json.RawMessage so malformed or unexpected digest fields
// are detected during explicit router decoding rather than silently coerced.
type Sip struct {
	ID       uuid.UUID       `db:"id" json:"id" swaggertype:"string" format:"uuid" example:"339366bc-464d-582f-8132-6875ccc814d2"`
	Kind     string          `db:"kind" json:"kind" example:"event"`
	Created  time.Time       `db:"created" json:"created_at" example:"2026-05-19T06:00:00Z"`
	SourceID uuid.UUID       `db:"source" json:"source_id,omitempty" swaggertype:"string" format:"uuid"`
	Tags     []string        `db:"tags" json:"tags"`
	Digest   json.RawMessage `db:"digest" swaggertype:"object" json:"digest,omitempty"`
	URL      sql.NullString  `db:"url" json:"url,omitempty"`
	BaseURL  sql.NullString  `db:"base_url" json:"base_url,omitempty"`

	// Distance is the embedding distance for semantic (relevance) queries. It is not persisted
	// and not exposed in public JSON; it is only used to build relevance cursors.
	Distance float64 `db:"distance" json:"-"`
}

func (sip *Sip) IsZero() bool {
	return sip.ID == uuid.Nil && sip.Kind == ""
}

func (sip *Sip) MaterializeDigest() (map[string]any, error) {
	var fields map[string]any
	if err := json.Unmarshal(sip.Digest, &fields); err != nil {
		return nil, err
	}
	return fields, nil
}

// Source describes a content publisher tracked in the database.
// Optional metadata is represented with pointers; missing optional metadata is not an error.
type Source struct {
	ID          uuid.UUID `db:"id" json:"id" swaggertype:"string" format:"uuid" example:"a1b2c3d4-e5f6-7890-abcd-ef1234567890"`
	BaseURL     string    `db:"base_url" json:"base_url" example:"https://example.com"`
	DomainName  *string   `db:"domain_name" json:"domain,omitempty" example:"example.com"`
	SiteName    *string   `db:"site_name" json:"name,omitempty" example:"Example News"`
	Description *string   `db:"description" json:"description,omitempty" example:"Independent business and policy coverage."`
	Favicon     *string   `db:"favicon" json:"favicon_url,omitempty" example:"https://example.com/favicon.ico"`
	RSSFeed     *string   `db:"rss_feed" json:"rss_feed_url,omitempty" example:"https://example.com/rss"`
}

func (source *Source) IsZero() bool {
	return source.ID == uuid.Nil && source.BaseURL == ""
}

// Relation links two sips by a named relationship (for example SAME_AS or DERIVED_FROM).
type Relation struct {
	FromID       uuid.UUID `db:"from_id" json:"from_id" swaggertype:"string" format:"uuid" example:"b07049b5-54c0-50b0-a620-d3aea3f8a173"`
	ToID         uuid.UUID `db:"to_id" json:"to_id" swaggertype:"string" format:"uuid" example:"9c3cc0a2-6eea-5290-9e9b-b5c462aeaa3a"`
	Relationship string    `db:"relationship" json:"relationship" example:"SAME_AS" enums:"SAME_AS,DERIVED_FROM"`
}

// EventEvidence is the narrow storage projection used by the bare evidence route.
// type EventEvidence struct {
// 	ID       uuid.UUID  `db:"id" json:"id"`
// 	Created  time.Time  `db:"created" json:"created"`
// 	SourceID *uuid.UUID `db:"source" json:"source_id,omitempty"`
// 	URL      *string    `db:"url" json:"url,omitempty"`
// 	BaseURL  *string    `db:"base_url" json:"base_url,omitempty"`
// }

type RelationCounts struct {
	SameAs      int64 `db:"same_as_count"`
	DerivedFrom int64 `db:"derived_from_count"`
	DerivedTo   int64 `db:"derived_to_count"`
}

// Page is a page of results plus the cursor to continue scanning, if more rows remain.
type Page[T any] struct {
	Items      []T
	NextCursor *Cursor
}

// const (
// 	SortRecent    string = "recent"
// 	SortRelevance string = "relevance"
// )

// const (
// 	SummaryCreatedDay string = "created_day"
// 	SummaryEventType  string = "event_type"
// 	SummaryImpact     string = "impact_level"
// 	SummarySource     string = "source"
// 	SummaryTag        string = "tag"
// 	SummaryRegion     string = "region"
// )

// // EventSummaryRow is one aggregate bucket returned by SummarizeEvents.
// type EventSummaryRow struct {
// 	Key           string `db:"key" json:"key"`
// 	EventCount    int64  `db:"event_count" json:"event_count"`
// 	CoverageCount int64  `db:"coverage_count" json:"coverage_count"`
// }

type Filters struct {
	IDs             []uuid.UUID
	SourceIDs       []uuid.UUID
	Kind            string
	CreatedFrom     time.Time
	CreatedTo       time.Time
	Tags            []string
	Companies       []string
	People          []string
	Products        []string
	Regions         []string
	EventTypes      []string
	ImpactLevels    []string
	ImpactedDomains []string
	Embedding       []float32
	Distance        *float64
	// Sort            string
}

// Cursor is the opaque, versioned pagination cursor. It encodes the last sort key and UUID
// so clients can resume a deterministic keyset scan without inspecting or constructing it.
type Cursor struct {
	Version  int        `json:"v"`
	Sort     *string    `json:"s,omitempty"`
	ID       *uuid.UUID `json:"id"`
	Created  *time.Time `json:"c,omitempty"`
	Distance *float64   `json:"d,omitempty"`
	TextKey  *string    `json:"k,omitempty"`
}

// PageRequest is the cursor-based page request. Callers (router) own limit defaults and caps.
type PageRequest struct {
	Limit  int
	Cursor *Cursor
}

// Cursor version. Bump when the encoded shape changes in a backward-incompatible way.
const _CURSOR_VERSION = 1

// ErrInvalidCursor is returned when a client-supplied cursor cannot be decoded or is malformed.
var ErrInvalidCursor = errors.New("invalid or malformed cursor")

// EncodeCursor serializes a cursor into an opaque, URL-safe string.
func (c *Cursor) Encode() (string, error) {
	if c == nil {
		return "", nil
	}
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// DecodeCursor parses an opaque cursor string. Empty string yields nil.
func DecodeCursor(raw string) (*Cursor, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	b, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, ErrInvalidCursor
	}
	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, ErrInvalidCursor
	}
	if c.Version != _CURSOR_VERSION || (c.ID == nil && c.TextKey == nil) {
		return nil, ErrInvalidCursor
	}
	return &c, nil
}
