package model

import (
	"time"

	"github.com/google/uuid"
)

// Metadata contains Dublin Core metadata for the EPUB package document.
type Metadata struct {
	Title       string    // dc:title (required)
	Authors     []string  // dc:creator (can be multiple)
	Language    string    // dc:language (BCP 47, e.g., "en", "en-US")
	Identifier  string    // dc:identifier (UUID or ISBN)
	Description string    // dc:description
	Publisher   string    // dc:publisher
	Date        time.Time // dc:date (publication date)
	Rights      string    // dc:rights
	CoverImage  string    // Path to cover image resource
}

// NewMetadata creates a new Metadata with default values.
func NewMetadata() *Metadata {
	return &Metadata{
		Language: "en",
		Authors:  make([]string, 0),
	}
}

// EnsureIdentifier generates a UUID identifier if not already set.
func (m *Metadata) EnsureIdentifier() {
	if m.Identifier == "" {
		m.Identifier = "urn:uuid:" + uuid.New().String()
	}
}

// EnsureDefaults sets default values for unset fields.
func (m *Metadata) EnsureDefaults() {
	if m.Language == "" {
		m.Language = "en"
	}
	m.EnsureIdentifier()
	if m.Date.IsZero() {
		m.Date = time.Now()
	}
}

// Merge combines two Metadata objects, with override taking precedence.
// Fields from override replace fields in base if they are non-empty.
func (m *Metadata) Merge(override *Metadata) {
	if override == nil {
		return
	}
	if override.Title != "" {
		m.Title = override.Title
	}
	if len(override.Authors) > 0 {
		m.Authors = override.Authors
	}
	if override.Language != "" {
		m.Language = override.Language
	}
	if override.Identifier != "" {
		m.Identifier = override.Identifier
	}
	if override.Description != "" {
		m.Description = override.Description
	}
	if override.Publisher != "" {
		m.Publisher = override.Publisher
	}
	if !override.Date.IsZero() {
		m.Date = override.Date
	}
	if override.Rights != "" {
		m.Rights = override.Rights
	}
	if override.CoverImage != "" {
		m.CoverImage = override.CoverImage
	}
}

// Valid checks if required metadata fields are present.
func (m *Metadata) Valid() bool {
	return m.Title != ""
}

// PrimaryAuthor returns the first author or empty string if none.
func (m *Metadata) PrimaryAuthor() string {
	if len(m.Authors) > 0 {
		return m.Authors[0]
	}
	return ""
}
