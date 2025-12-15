# Data Model: EPUB Converter CLI

**Date**: 2025-12-15
**Feature**: 001-epub-converter-cli

## Overview

This document defines the internal data structures used by the EPUB Converter CLI. The converter uses an intermediate document representation to decouple input parsing from EPUB generation.

---

## Core Entities

### 1. Document

The intermediate representation of parsed content, independent of source format.

```go
// Document represents parsed content ready for EPUB generation
type Document struct {
    Metadata  Metadata       // Book metadata
    Chapters  []Chapter      // Content chapters
    Resources []Resource     // Images, stylesheets, fonts
    TOC       TableOfContents // Navigation structure
}
```

**Fields:**

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| Metadata | Metadata | Book publication information | Title required |
| Chapters | []Chapter | Ordered content sections | At least one chapter |
| Resources | []Resource | Embedded media files | Valid file references |
| TOC | TableOfContents | Navigation hierarchy | Auto-generated from headings |

---

### 2. Metadata

Dublin Core metadata for the EPUB package document.

```go
// Metadata contains book publication information
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
```

**Fields:**

| Field | Type | Required | Source Priority |
|-------|------|----------|-----------------|
| Title | string | Yes | CLI flag > Front matter > Filename |
| Authors | []string | No | CLI flag > Front matter |
| Language | string | No | CLI flag > Front matter > "en" default |
| Identifier | string | Auto | Generated UUID if not provided |
| CoverImage | string | No | CLI flag > Front matter |

**Validation Rules:**
- Title must be non-empty
- Language must be valid BCP 47 code
- Identifier auto-generated as UUID v4 if not provided

---

### 3. Chapter

A content section of the book, typically corresponding to one XHTML file in the EPUB.

```go
// Chapter represents a content section
type Chapter struct {
    ID       string  // Unique identifier (e.g., "chapter-01")
    Title    string  // Chapter title (for TOC)
    Level    int     // Heading level (1-6)
    Content  string  // XHTML content
    FileName string  // Output filename (e.g., "chapter-01.xhtml")
    Order    int     // Spine order
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| ID | string | Unique identifier, used in manifest and spine |
| Title | string | Display title for TOC entry |
| Level | int | Depth in TOC hierarchy (1 = top level) |
| Content | string | Valid XHTML content |
| FileName | string | File path within EPUB (OEBPS/content/) |
| Order | int | Reading order position |

**State Transitions:**
```
[Parsed] → [Validated] → [Written to EPUB]
```

---

### 4. Resource

Embedded media file (image, stylesheet, font).

```go
// Resource represents an embedded media file
type Resource struct {
    ID        string // Unique identifier for manifest
    FileName  string // Path within EPUB
    MediaType string // MIME type (e.g., "image/png")
    Data      []byte // File contents
    IsCover   bool   // True if this is the cover image
}
```

**Supported Media Types:**

| Type | Extensions | MediaType |
|------|------------|-----------|
| JPEG | .jpg, .jpeg | image/jpeg |
| PNG | .png | image/png |
| GIF | .gif | image/gif |
| SVG | .svg | image/svg+xml |
| WebP | .webp | Converted to image/png |
| CSS | .css | text/css |

**Validation Rules:**
- MediaType must be EPUB-compatible
- WebP converted to PNG before storage
- Unsupported formats generate warning, skipped

---

### 5. TableOfContents

Hierarchical navigation structure.

```go
// TableOfContents represents the navigation hierarchy
type TableOfContents struct {
    Entries []TOCEntry
}

// TOCEntry is a single navigation item
type TOCEntry struct {
    Title    string     // Display text
    Href     string     // Link to content (e.g., "chapter-01.xhtml")
    Level    int        // Hierarchy depth (1-6)
    Children []TOCEntry // Nested entries
}
```

**Generation Rules:**
- Built automatically from heading structure (h1-h6)
- h1 creates top-level entry
- h2-h6 create nested entries based on level
- Flat list for documents without clear hierarchy

---

### 6. ConversionResult

Output of the conversion process.

```go
// ConversionResult contains conversion outcome
type ConversionResult struct {
    Success    bool     // True if conversion completed
    OutputPath string   // Path to generated EPUB
    Warnings   []string // Non-fatal issues encountered
    Error      error    // Fatal error (if Success is false)
    Stats      ConversionStats
}

// ConversionStats contains metrics
type ConversionStats struct {
    InputFormat   string        // "markdown", "html", "pdf"
    InputFiles    int           // Number of input files
    ChapterCount  int           // Number of chapters generated
    ImageCount    int           // Number of images embedded
    OutputSize    int64         // EPUB file size in bytes
    Duration      time.Duration // Processing time
}
```

---

## Input Format Models

### MarkdownInput

```go
// MarkdownInput represents parsed Markdown source
type MarkdownInput struct {
    FrontMatter map[string]interface{} // YAML/TOML metadata
    Body        []byte                  // Markdown content
    Images      []ImageReference        // Referenced images
}
```

### HTMLInput

```go
// HTMLInput represents parsed HTML source
type HTMLInput struct {
    Title   string           // From <title> tag
    Meta    map[string]string // From <meta> tags
    Body    string           // <body> content
    CSS     []string         // Inline and linked CSS
    Images  []ImageReference // <img> references
}
```

### PDFInput

```go
// PDFInput represents extracted PDF content
type PDFInput struct {
    Pages   []PDFPage // Extracted pages
    Images  []ImageReference
}

// PDFPage represents content from one PDF page
type PDFPage struct {
    Text     string      // Extracted text
    Headings []PDFHeading // Detected headings
}

// PDFHeading represents a detected heading
type PDFHeading struct {
    Text     string
    Level    int     // Inferred from font size
    FontSize float64 // Original font size
}
```

---

## Relationships

```
Document
├── Metadata (1:1)
├── Chapters (1:N)
│   └── references Resources via image URLs
├── Resources (1:N)
│   └── includes CoverImage from Metadata
└── TOC (1:1)
    └── entries link to Chapters
```

---

## Entity Lifecycle

### Document Processing Pipeline

```
1. INPUT PHASE
   MarkdownInput / HTMLInput / PDFInput
           ↓
2. PARSING PHASE
   Parser extracts content, metadata, images
           ↓
3. DOCUMENT PHASE
   Document (intermediate representation)
           ↓
4. GENERATION PHASE
   EPUB Builder creates package
           ↓
5. OUTPUT PHASE
   ConversionResult with EPUB file
```

### Resource States

```
[Referenced] → [Downloaded/Read] → [Validated] → [Converted?] → [Embedded]
     ↓               ↓                  ↓
  [Missing]      [Error]         [Unsupported]
  (warning)      (warning)        (warning)
```

---

## Validation Summary

| Entity | Validation |
|--------|------------|
| Document | Has title, at least one chapter |
| Metadata | Title non-empty, language valid BCP 47 |
| Chapter | Content is valid XHTML, filename unique |
| Resource | MediaType EPUB-compatible, data not empty |
| TOCEntry | Href points to valid chapter |
