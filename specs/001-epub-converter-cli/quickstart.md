# Quickstart: EPUB Converter CLI

**Date**: 2025-12-15
**Feature**: 001-epub-converter-cli

## Installation

### From Source

```bash
# Clone repository
git clone https://github.com/dauquangthanh/epub-converter.git
cd toepub

# Build
go build -o toepub ./cmd/toepub

# Install to PATH (optional)
go install ./cmd/toepub
```

### Pre-built Binary

```bash
# Download for your platform
curl -L https://github.com/dauquangthanh/epub-converter/releases/latest/download/toepub-$(uname -s)-$(uname -m).tar.gz | tar xz

# Move to PATH
sudo mv toepub /usr/local/bin/
```

---

## Basic Usage

### Convert Markdown to EPUB

```bash
# Simple conversion
toepub convert document.md

# Output: document.epub in current directory
```

### Convert HTML to EPUB

```bash
toepub convert article.html
```

### Convert PDF to EPUB

```bash
toepub convert book.pdf
```

---

## Common Tasks

### Set Book Metadata

```bash
toepub convert document.md \
  --title "My Book Title" \
  --author "John Doe" \
  --language en-US
```

### Add Cover Image

```bash
toepub convert document.md --cover cover.jpg -o book.epub
```

### Combine Multiple Files

```bash
# Files combined in alphabetical order
toepub convert chapter1.md chapter2.md chapter3.md -o book.epub

# Or use a directory
toepub convert ./chapters/ -o book.epub
```

### Custom Output Path

```bash
toepub convert document.md --output ~/Books/mybook.epub
```

---

## Using Markdown Front Matter

Add YAML front matter to your Markdown file for automatic metadata:

```markdown
---
title: My Book Title
author: Jane Author
language: en-US
---

# Chapter 1

Your content here...
```

Then convert:

```bash
toepub convert document.md
# Metadata automatically extracted from front matter
```

CLI flags override front matter when both are present.

---

## JSON Output for Scripting

```bash
# Get structured output
toepub convert document.md --format json

# Example output:
# {
#   "success": true,
#   "output": "/path/to/document.epub",
#   "stats": { ... }
# }
```

Use with `jq`:

```bash
# Check if successful
toepub convert doc.md -f json | jq '.success'

# Get output path
toepub convert doc.md -f json | jq -r '.output'
```

---

## Reading from Stdin

```bash
# Pipe markdown content
cat README.md | toepub convert - -o readme.epub

# Specify format when auto-detection not possible
curl -s https://example.com/doc | toepub convert - --input-format html -o doc.epub
```

---

## Verification

After conversion, verify your EPUB:

```bash
# Using epubcheck (requires Java)
epubcheck document.epub

# Or online at https://validator.w3.org/ebooks/
```

---

## Troubleshooting

### "File not found" Error

```bash
# Check file exists
ls -la document.md

# Use absolute path
toepub convert /full/path/to/document.md
```

### "Unsupported image format" Warning

The converter supports JPEG, PNG, GIF, and SVG. WebP images are automatically converted to PNG. Other formats are skipped with a warning.

```bash
# Check which images were skipped
toepub convert document.md --format json | jq '.warnings'
```

### "PDF text extraction failed" Error

This usually means the PDF is scanned/image-based. OCR is required first:

```bash
# Use external OCR tool first
ocrmypdf scanned.pdf text-based.pdf
toepub convert text-based.pdf
```

---

## Examples

### Technical Documentation

```bash
# Convert project docs
toepub convert \
  --title "Project Documentation" \
  --author "Development Team" \
  README.md \
  docs/getting-started.md \
  docs/api.md \
  docs/faq.md \
  -o project-docs.epub
```

### Blog Post Collection

```bash
# Convert blog posts directory
toepub convert \
  --title "2024 Blog Posts" \
  --author "Blog Author" \
  --cover blog-cover.png \
  ./posts/ \
  -o blog-2024.epub
```

### Single Article

```bash
# Quick conversion with all defaults
toepub convert article.html
# Creates article.epub
```

---

## Getting Help

```bash
# General help
toepub --help

# Command-specific help
toepub convert --help

# Version info
toepub version
```

---

## Next Steps

- Read the full [CLI Contract](./contracts/cli-contract.md) for all options
- See [Data Model](./data-model.md) for internal structure details
- Check [Research](./research.md) for technology decisions
