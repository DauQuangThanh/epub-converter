# toepub

A CLI tool to convert Markdown, HTML, and PDF files to EPUB 3+ format.

## Features

- **Markdown Conversion**: Full GFM (GitHub Flavored Markdown) support including tables, task lists, and code blocks
- **HTML Conversion**: HTML5 to XHTML conversion with CSS extraction and JavaScript stripping
- **PDF Conversion**: Text extraction with heading detection and structure preservation
- **Metadata Override**: Set title, author, language, and cover image via CLI flags
- **Multiple Output Formats**: Human-readable and JSON output for CI/CD integration
- **EPUB 3.3 Compliant**: Generates valid EPUB 3.3 files with proper navigation

## Installation

```bash
go install github.com/dauquangthanh/epub-converter/cmd/toepub@latest
```

Or build from source:

```bash
git clone https://github.com/dauquangthanh/epub-converter.git
cd epub-converter
go build -o toepub ./cmd/toepub
```

## Usage

### Basic Conversion

```bash
# Convert Markdown to EPUB
toepub convert document.md

# Convert HTML to EPUB
toepub convert page.html -o output.epub

# Convert PDF to EPUB
toepub convert document.pdf --title "My Book"
```

### Metadata Flags

```bash
toepub convert document.md \
  --title "My Book Title" \
  --author "John Doe" \
  --language "en" \
  --cover cover.jpg \
  --output mybook.epub
```

### Reading from Stdin

```bash
cat document.md | toepub convert - --input-format md -o output.epub
```

### JSON Output

```bash
toepub convert document.md --format json
```

Output:
```json
{
  "success": true,
  "output_path": "document.epub",
  "stats": {
    "input_format": "markdown",
    "input_files": 1,
    "chapter_count": 5,
    "image_count": 3,
    "output_size": 45678,
    "duration_ms": 125
  },
  "warnings": []
}
```

### Multiple Files

```bash
# Convert multiple files into a single EPUB
toepub convert chapter1.md chapter2.md chapter3.md -o book.epub

# Convert all Markdown files in a directory
toepub convert ./chapters/ -o book.epub
```

## CLI Reference

```
toepub convert <input...> [flags]

Flags:
  -o, --output string        Output EPUB file path
  -f, --format string        Output format: human (default), json
  -t, --title string         Override document title
  -a, --author string        Override document author (repeatable)
  -l, --language string      Override document language (default "en")
  -c, --cover string         Cover image path
      --input-format string  Force input format: md, html, pdf
  -h, --help                 Help for convert
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 64 | File not found |
| 65 | Format/parsing error |
| 66 | Output not writable |
| 70 | Internal error |

## Input Formats

### Markdown

- Full GFM support (tables, task lists, strikethrough, autolinks)
- YAML front matter for metadata:

```markdown
---
title: My Document
author: Jane Doe
language: en
---

# Chapter 1
Content here...
```

### HTML

- HTML5 input with automatic XHTML conversion
- Metadata extraction from `<title>` and `<meta>` tags
- CSS extraction from `<style>` tags
- JavaScript automatically stripped

### PDF

- Text extraction with structure preservation
- Heading detection based on font size
- Note: Complex layouts and scanned PDFs may have limited support

## Development

### Prerequisites

- Go 1.24+
- golangci-lint (for linting)

### Commands

```bash
# Build
go build -o toepub ./cmd/toepub

# Test
go test ./...

# Test with coverage
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Lint
golangci-lint run
```

### Project Structure

```
cmd/toepub/          # CLI entry point
internal/
├── cli/             # CLI commands
├── parser/          # Format parsers (markdown, html, pdf)
├── epub/            # EPUB generation
├── converter/       # Conversion orchestration
└── model/           # Data structures
tests/fixtures/      # Test input files
```

## License

GPL-3.0 License - see [LICENSE](LICENSE) for details.
Author: Dau Quang Thanh
