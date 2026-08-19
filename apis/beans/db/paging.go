package db

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Cursor version. Bump when the encoded shape changes in a backward-incompatible way.
const _CURSOR_VERSION = 1

var ErrInvalidCursor = errors.New("invalid or malformed cursor")

// Cursor is the opaque, versioned pagination cursor. It encodes the last sort key and UUID
// so clients can resume a deterministic keyset scan without inspecting or constructing it.
type Cursor struct {
	Version    int        `json:"v"`
	Sort       string     `json:"s,omitempty"`
	ID         *uuid.UUID `json:"id"`
	Created    *time.Time `json:"c,omitempty"`
	TrendScore *float64   `json:"ts,omitempty"`
	Distance   *float64   `json:"d,omitempty"`
	TextKey    *string    `json:"k,omitempty"`
}

// PageRequest is the cursor-based page request. Callers (router) own limit defaults and caps.
type PageRequest struct {
	Limit  int
	Cursor *Cursor
}

// Page is a page of results plus the cursor to continue scanning, if more rows remain.
type Page[T any] struct {
	Items      []T
	NextCursor *Cursor
}

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
