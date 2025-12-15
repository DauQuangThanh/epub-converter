// ------------------------------------------------------------------
// Developed by Dau Quang Thanh - 2025.
// Enterprise AI Solution Architect
//
// Happy Reading!
// ------------------------------------------------------------------

// Package parser provides input format parsers for the EPUB converter.
//
// The parser package implements parsers for Markdown, HTML, and PDF formats.
// Each parser converts input content into an intermediate Document representation
// that can be processed by the EPUB generator.
package parser

import (
	"github.com/dauquangthanh/epub-converter/internal/model"
)

// Parser defines the interface for input format parsers.
// All parsers convert their respective input formats into a common
// Document representation for EPUB generation.
type Parser interface {
	// Parse converts input content to a Document.
	// The content parameter contains the raw file content.
	// The basePath parameter is used to resolve relative paths (e.g., images).
	Parse(content []byte, basePath string) (*model.Document, error)

	// SupportedExtensions returns file extensions this parser handles.
	SupportedExtensions() []string
}

// Format represents supported input formats.
type Format string

const (
	FormatMarkdown Format = "markdown"
	FormatHTML     Format = "html"
	FormatPDF      Format = "pdf"
	FormatUnknown  Format = "unknown"
)

// String returns the string representation of the format.
func (f Format) String() string {
	return string(f)
}
