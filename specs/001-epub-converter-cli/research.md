# Research: EPUB Converter CLI

**Date**: 2025-12-15
**Feature**: 001-epub-converter-cli

## Executive Summary

This document consolidates research findings for building a Go CLI tool that converts Markdown, HTML, and PDF to EPUB 3+ format. Key findings:

1. **No mature Go EPUB 3 library exists** - Custom EPUB generation recommended
2. **Goldmark** is the clear choice for Markdown parsing (GFM support, XHTML output)
3. **PDF extraction has limitations** - Open source options are basic; commercial (UniPDF) offers better quality
4. **EPUB structure** can be generated using Go's archive/zip + templates

---

## 1. EPUB 3 Generation

### Decision: Custom EPUB 3 Generator with Templates

### Rationale
No existing Go library provides comprehensive EPUB 3.3 support. Available libraries (go-shiori/go-epub) offer basic EPUB 3.0 but lack:
- nav.xhtml with landmarks (`epub:type="landmarks"`)
- epub:type semantic attributes
- Full epubcheck validation support

### Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| go-shiori/go-epub | Working EPUB 3.0 foundation | Missing landmarks, minimal maintenance (last release Dec 2023) | REJECTED |
| Pandoc via go-pandoc | Mature EPUB support | External dependency, validation issues | REJECTED |
| Custom generator | Full control, EPUB 3.3 compliant | Development effort | **SELECTED** |

### Implementation Approach
- Use Go's `archive/zip` for EPUB (ZIP) creation
- Use `html/template` for generating XHTML content, nav.xhtml, content.opf
- Follow EPUB 3.3 specification strictly
- Validate with epubcheck during testing

### EPUB 3 Required Structure
```
mimetype                    # First file, uncompressed
META-INF/
└── container.xml           # Points to content.opf
OEBPS/
├── content.opf             # Package document
├── nav.xhtml               # Navigation (TOC + landmarks)
├── content/                # XHTML content files
├── images/                 # Embedded images
└── styles/                 # CSS stylesheets
```

---

## 2. Markdown Parser

### Decision: Goldmark with Extensions

### Rationale
Goldmark is the industry standard Go Markdown parser (used by Hugo). It provides:
- Full GFM support via `extension.GFM`
- Native XHTML output via `html.WithXHTML()`
- Excellent front matter support via goldmark-frontmatter
- CommonMark compliant
- Active maintenance (latest release July 2025)

### Alternatives Considered

| Library | GFM Support | XHTML | Front Matter | Maintenance | Decision |
|---------|-------------|-------|--------------|-------------|----------|
| goldmark | Excellent | Native | Via extension | Active | **SELECTED** |
| gomarkdown | Very Good | Native | Limited | Active | Alternative |
| blackfriday | Good | Limited | None | Inactive | REJECTED |

### Configuration
```go
md := goldmark.New(
    goldmark.WithExtensions(
        extension.GFM,  // Tables, task lists, strikethrough, autolinks
        &frontmatter.Extender{},  // YAML front matter
    ),
    goldmark.WithRendererOptions(
        html.WithXHTML(),  // XHTML output for EPUB
    ),
)
```

### Dependencies
- `github.com/yuin/goldmark` - Main parser
- `go.abhg.dev/goldmark/frontmatter` - YAML/TOML front matter

---

## 3. HTML Parser

### Decision: golang.org/x/net/html + CSS Inlining

### Rationale
Go's extended standard library provides robust HTML parsing. For EPUB:
- Parse HTML to AST for manipulation
- Convert to XHTML (self-closing tags, lowercase elements)
- Inline or extract CSS for EPUB compatibility

### Implementation Notes
- Use `golang.org/x/net/html` for parsing
- Convert HTML5 to XHTML via AST transformation
- Handle remote images: download and embed
- Strip JavaScript (not executable in EPUB)

### Dependencies
- `golang.org/x/net/html` - HTML parsing
- Custom CSS handler or `github.com/tdewolff/parse/css` for CSS parsing

---

## 4. PDF Text Extraction

### Decision: ledongthuc/pdf (Open Source) with Fallback Strategy

### Rationale
PDF text extraction in Go has significant limitations for open-source options:

| Library | Text Quality | Images | License | Decision |
|---------|--------------|--------|---------|----------|
| UniPDF | Excellent | Yes | Commercial ($3k/yr) | Too expensive for OSS |
| pdfcpu | None (content streams only) | Excellent | Apache 2.0 | For images only |
| ledongthuc/pdf | Basic | Limited | Open Source | **SELECTED** |

### Implementation Strategy
1. Use `ledongthuc/pdf` for text extraction with style metadata
2. Use `pdfcpu` for image extraction
3. Build custom heading detection using font size heuristics
4. Accept lower quality for PDF conversion (per P3 priority)

### Limitations to Document
- Complex PDF layouts may not extract perfectly
- Scanned/image-based PDFs will produce error (per spec)
- Multi-column text may merge incorrectly

### Dependencies
- `github.com/ledongthuc/pdf` - Text extraction
- `github.com/pdfcpu/pdfcpu` - Image extraction

---

## 5. Image Processing

### Decision: Standard Library + golang.org/x/image

### Rationale
EPUB supports JPEG, PNG, GIF, SVG. WebP needs conversion to PNG.

### Implementation
- Use `image/png`, `image/jpeg`, `image/gif` from stdlib
- Use `golang.org/x/image/webp` for WebP decoding
- Encode WebP as PNG for EPUB inclusion
- SVG passed through unchanged

### Dependencies
- Standard library image packages
- `golang.org/x/image/webp` - WebP support

---

## 6. CLI Framework

### Decision: Cobra

### Rationale
Cobra is the de facto standard for Go CLIs (used by kubectl, hugo, docker CLI). It provides:
- Automatic help generation
- Flag parsing
- Subcommand support
- Shell completions

### Command Structure
```
toepub convert <input> [flags]
  --output, -o     Output file path
  --format, -f     Output format (human, json)
  --title          Override title metadata
  --author         Override author metadata
  --language       Override language metadata
  --cover          Cover image path
  --input-format   Force input format (md, html, pdf)
```

### Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/pflag` - POSIX flag parsing (included with cobra)

---

## 7. Final Dependency List

```go
// go.mod
module github.com/dauquangthanh/epub-converter

go 1.24

require (
    // CLI
    github.com/spf13/cobra v1.8.0

    // Markdown
    github.com/yuin/goldmark v1.7.12
    go.abhg.dev/goldmark/frontmatter v0.2.0

    // HTML
    golang.org/x/net v0.23.0

    // PDF
    github.com/ledongthuc/pdf v0.0.0-20240201131950-da5b75280b06
    github.com/pdfcpu/pdfcpu v0.8.0

    // Images
    golang.org/x/image v0.15.0

    // Testing
    github.com/stretchr/testify v1.9.0
)
```

---

## 8. Architecture Alignment

This design aligns with ground-rules:

| Principle | Alignment |
|-----------|-----------|
| EPUB 3+ Compliance | Custom generator ensures full EPUB 3.3 support |
| CLI-First | Cobra provides standard CLI patterns |
| Input Format Fidelity | Goldmark preserves GFM; HTML converted to XHTML |
| Maintainability | Modular structure: parsers/, epub/, cli/ |
| Test-Driven | Integration tests with epubcheck validation |

---

## 9. Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| PDF extraction quality | High | Medium | Document limitations; PDF is P3 priority |
| EPUB validation failures | Medium | High | Test with epubcheck during development |
| WebP conversion quality | Low | Low | PNG encoding is lossless |
| Large file memory | Medium | Medium | Stream processing for >10MB files |

---

## Sources

- EPUB 3.3 Specification: https://www.w3.org/TR/epub-33/
- goldmark: https://github.com/yuin/goldmark
- goldmark-frontmatter: https://go.abhg.dev/goldmark/frontmatter
- go-shiori/go-epub: https://github.com/go-shiori/go-epub
- ledongthuc/pdf: https://github.com/ledongthuc/pdf
- pdfcpu: https://github.com/pdfcpu/pdfcpu
- cobra: https://github.com/spf13/cobra
