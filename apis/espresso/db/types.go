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

const (
	EVENT_TAG_TYPE_REGION       = "region"
	EVENT_TAG_TYPE_PEOPLE       = "people"
	EVENT_TAG_TYPE_PRODUCT      = "product"
	EVENT_TAG_TYPE_COMPANY      = "company"
	EVENT_TAG_TYPE_STOCK_TICKER = "stock_ticker"
	EVENT_TAG_TYPE_EVENT_TYPE   = "event_type"
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
	Tags     []string        `db:"tags" json:"tags"`
	Digest   json.RawMessage `db:"digest" swaggertype:"object" json:"digest,omitempty"`
	URL      sql.NullString  `db:"url" json:"url,omitempty"`
	SourceID uuid.UUID       `db:"source" json:"source_id,omitempty" swaggertype:"string" format:"uuid"`
	BaseURL  sql.NullString  `db:"base_url" json:"base_url,omitempty"`
	// Distance is the embedding distance for semantic (relevance) queries. It is not persisted
	Distance sql.NullFloat64 `db:"distance" json:"-"`
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
	ID          uuid.UUID      `db:"id" json:"id,omitzero" swaggertype:"string" format:"uuid" example:"a1b2c3d4-e5f6-7890-abcd-ef1234567890"`
	BaseURL     string         `db:"base_url" json:"base_url,omitempty" example:"https://example.com"`
	DomainName  string         `db:"domain_name" json:"domain,omitempty" example:"example.com"`
	SiteName    sql.NullString `db:"site_name" json:"name,omitempty" example:"Example News"`
	Description sql.NullString `db:"description" json:"description,omitempty" example:"Independent business and policy coverage."`
	Favicon     sql.NullString `db:"favicon" json:"favicon_url,omitempty" example:"https://example.com/favicon.ico"`
	RSSFeed     sql.NullString `db:"rss_feed" json:"rss_feed_url,omitempty" example:"https://example.com/rss"`
}

func (source *Source) IsZero() bool {
	return source.ID == uuid.Nil && source.BaseURL == ""
}

type ExtendedSip struct {
	Sip
	DomainName  sql.NullString `db:"domain_name" json:"domain,omitempty" example:"example.com"`
	SiteName    sql.NullString `db:"site_name" json:"name,omitempty" example:"Example News"`
	Description sql.NullString `db:"description" json:"description,omitempty" example:"Independent business and policy coverage."`
	Favicon     sql.NullString `db:"favicon" json:"favicon_url,omitempty" example:"https://example.com/favicon.ico"`
	RSSFeed     sql.NullString `db:"rss_feed" json:"rss_feed_url,omitempty" example:"https://example.com/rss"`
}

func (sip *ExtendedSip) GetSource() *Source {
	return &Source{
		ID:          sip.SourceID,
		BaseURL:     sip.BaseURL.String,
		DomainName:  sip.DomainName.String,
		SiteName:    sip.SiteName,
		Description: sip.Description,
		Favicon:     sip.Favicon,
		RSSFeed:     sip.RSSFeed,
	}
}

// Relation links two sips by a named relationship (for example SAME_AS or DERIVED_FROM).
type Relation struct {
	FromID       uuid.UUID `db:"from_id" json:"from_id" swaggertype:"string" format:"uuid" example:"b07049b5-54c0-50b0-a620-d3aea3f8a173"`
	ToID         uuid.UUID `db:"to_id" json:"to_id" swaggertype:"string" format:"uuid" example:"9c3cc0a2-6eea-5290-9e9b-b5c462aeaa3a"`
	Relationship string    `db:"relationship" json:"relationship" example:"SAME_AS" enums:"SAME_AS,DERIVED_FROM"`
}

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

// TagValue is one exact value exposed by a discovery route.
type TagValue struct {
	Value string `db:"value" json:"value" toon:"value"`
	Type  string `db:"type" json:"type,omitempty" toon:"type,omitempty" example:"region"`
}

type Filters struct {
	IDs             []uuid.UUID
	SourceIDs       []uuid.UUID
	Kind            string
	CreatedFrom     time.Time
	CreatedTo       time.Time
	Tags            []string
	Entities        []string
	Categories      []string
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
	EventTag *TagValue  `json:"et,omitempty"`
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
	if c.Version != _CURSOR_VERSION || (c.ID == nil && c.TextKey == nil && c.EventTag == nil) {
		return nil, ErrInvalidCursor
	}
	return &c, nil
}
