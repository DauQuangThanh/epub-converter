// Package model provides data structures for the EPUB converter.
package model

import "time"

// Document represents parsed content ready for EPUB generation.
// It serves as the intermediate representation between input parsers
// and the EPUB builder.
type Document struct {
	Metadata  Metadata        // Book publication information
	Chapters  []Chapter       // Content chapters in reading order
	Resources []Resource      // Embedded media files (images, stylesheets)
	TOC       TableOfContents // Navigation hierarchy
}

// NewDocument creates a new Document with initialized slices.
func NewDocument() *Document {
	return &Document{
		Chapters:  make([]Chapter, 0),
		Resources: make([]Resource, 0),
		TOC:       TableOfContents{Entries: make([]TOCEntry, 0)},
	}
}

// AddChapter appends a chapter to the document.
func (d *Document) AddChapter(chapter Chapter) {
	d.Chapters = append(d.Chapters, chapter)
}

// AddResource appends a resource to the document.
func (d *Document) AddResource(resource Resource) {
	d.Resources = append(d.Resources, resource)
}

// Valid checks if the document has required fields.
func (d *Document) Valid() bool {
	return d.Metadata.Title != "" && len(d.Chapters) > 0
}

// Chapter represents a content section of the book.
// Each chapter typically corresponds to one XHTML file in the EPUB.
type Chapter struct {
	ID       string // Unique identifier (e.g., "chapter-01")
	Title    string // Chapter title for TOC display
	Level    int    // Heading level (1-6) for hierarchy
	Content  string // XHTML content
	FileName string // Output filename (e.g., "chapter-01.xhtml")
	Order    int    // Reading order position in spine
}

// Resource represents an embedded media file (image, stylesheet, font).
type Resource struct {
	ID        string // Unique identifier for manifest
	FileName  string // Path within EPUB (e.g., "images/photo.png")
	MediaType string // MIME type (e.g., "image/png")
	Data      []byte // File contents
	IsCover   bool   // True if this is the cover image
}

// ConversionResult contains the outcome of a conversion operation.
type ConversionResult struct {
	Success    bool             // True if conversion completed successfully
	OutputPath string           // Path to generated EPUB file
	Warnings   []string         // Non-fatal issues encountered
	Error      error            // Fatal error if Success is false
	Stats      ConversionStats  // Conversion metrics
}

// ConversionStats contains metrics about the conversion process.
type ConversionStats struct {
	InputFormat  string        // Source format: "markdown", "html", "pdf"
	InputFiles   int           // Number of input files processed
	ChapterCount int           // Number of chapters generated
	ImageCount   int           // Number of images embedded
	OutputSize   int64         // EPUB file size in bytes
	Duration     time.Duration // Processing time
}

// AddWarning appends a warning message to the result.
func (r *ConversionResult) AddWarning(msg string) {
	r.Warnings = append(r.Warnings, msg)
}
