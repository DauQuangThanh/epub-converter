# CLI Contract: EPUB Converter CLI

**Date**: 2025-12-15
**Feature**: 001-epub-converter-cli

## Overview

This document defines the command-line interface contract for the `toepub` converter tool.

---

## Command Structure

```
toepub [command] [flags]

Commands:
  convert    Convert input file(s) to EPUB format
  version    Print version information
  help       Help about any command

Global Flags:
  -h, --help      Help for toepub
  -v, --version   Print version and exit
```

---

## Convert Command

### Synopsis

```
toepub convert <input>... [flags]
toepub convert <directory> [flags]
cat file.md | toepub convert - [flags]
```

### Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| input | string(s) | Yes | Input file path(s), directory, or "-" for stdin |

### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| --output | -o | string | `<input>.epub` | Output file path |
| --format | -f | string | human | Output format: `human` or `json` |
| --title | -t | string | (from source) | Override book title |
| --author | -a | string | (from source) | Override author name |
| --language | -l | string | en | Book language (BCP 47 code) |
| --cover | -c | string | (none) | Cover image path |
| --input-format | | string | (auto) | Force input format: `md`, `html`, `pdf` |
| --help | -h | | | Help for convert command |

### Input Handling

**File Input:**
```bash
# Single file
toepub convert document.md

# Multiple files (combined into single EPUB, alphabetical order)
toepub convert chapter1.md chapter2.md chapter3.md

# Directory (all supported files, alphabetical order, non-recursive)
toepub convert ./docs/

# Mixed
toepub convert intro.md ./chapters/ appendix.md
```

**Stdin Input:**
```bash
# Pipe content
cat document.md | toepub convert -

# Stdin with explicit format
cat document.html | toepub convert - --input-format html
```

### Output Path Resolution

| Scenario | Default Output |
|----------|----------------|
| Single file `doc.md` | `./doc.epub` |
| Multiple files | `./output.epub` (or first file name) |
| Directory `./docs/` | `./docs.epub` |
| Stdin | `./output.epub` |
| With `-o path` | Specified path |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (invalid input, conversion failure) |
| 2 | Invalid arguments or flags |
| 64 | Input file not found |
| 65 | Input format not supported or detected |
| 66 | Output path not writable |
| 70 | Internal error |

---

## Output Formats

### Human-Readable (Default)

**Success:**
```
Converting document.md...
✓ Created document.epub (145 KB)
  - 5 chapters
  - 12 images
  - Duration: 1.2s
```

**With Warnings:**
```
Converting document.md...
⚠ Warning: Image not found: missing.png
⚠ Warning: Unsupported image format: photo.webp (converted to PNG)
✓ Created document.epub (145 KB)
  - 5 chapters
  - 11 images (1 converted)
  - Duration: 1.5s
```

**Error:**
```
Error: Cannot read input file: document.md
  → File does not exist

Suggestion: Check the file path and try again.
```

### JSON Format (`--format json`)

**Success:**
```json
{
  "success": true,
  "output": "/path/to/document.epub",
  "stats": {
    "input_format": "markdown",
    "input_files": 1,
    "chapters": 5,
    "images": 12,
    "output_size": 148480,
    "duration_ms": 1234
  },
  "warnings": []
}
```

**With Warnings:**
```json
{
  "success": true,
  "output": "/path/to/document.epub",
  "stats": {
    "input_format": "markdown",
    "input_files": 1,
    "chapters": 5,
    "images": 11,
    "output_size": 148480,
    "duration_ms": 1543
  },
  "warnings": [
    "Image not found: missing.png",
    "Unsupported image format: photo.webp (converted to PNG)"
  ]
}
```

**Error:**
```json
{
  "success": false,
  "error": {
    "code": 64,
    "message": "Cannot read input file: document.md",
    "detail": "File does not exist"
  }
}
```

---

## Version Command

### Synopsis

```
toepub version
toepub --version
toepub -v
```

### Output

```
toepub version 1.0.0
EPUB 3.3 compliant output
Built with Go 1.21
```

---

## Help Output

### Root Help

```
toepub - Convert Markdown, HTML, and PDF to EPUB 3+

Usage:
  toepub [command]

Available Commands:
  convert     Convert input file(s) to EPUB format
  help        Help about any command
  version     Print version information

Flags:
  -h, --help      help for toepub
  -v, --version   version for toepub

Use "toepub [command] --help" for more information about a command.
```

### Convert Help

```
Convert input file(s) to EPUB 3+ format.

Supports Markdown (.md), HTML (.html, .htm), and PDF (.pdf) input.
Multiple files or directories are combined into a single EPUB.

Usage:
  toepub convert <input>... [flags]

Examples:
  # Convert single Markdown file
  toepub convert document.md

  # Convert with custom output path
  toepub convert document.md -o book.epub

  # Convert multiple files
  toepub convert chapter1.md chapter2.md chapter3.md

  # Convert directory
  toepub convert ./docs/

  # Set metadata
  toepub convert document.md --title "My Book" --author "John Doe"

  # Add cover image
  toepub convert document.md --cover cover.jpg

  # JSON output for scripting
  toepub convert document.md --format json

  # From stdin
  cat document.md | toepub convert -

Flags:
  -a, --author string        Override author name
  -c, --cover string         Cover image path
  -f, --format string        Output format: human or json (default "human")
  -h, --help                 help for convert
      --input-format string  Force input format: md, html, pdf
  -l, --language string      Book language (BCP 47 code) (default "en")
  -o, --output string        Output file path
  -t, --title string         Override book title
```

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| TOEPUB_DEFAULT_LANGUAGE | Default language code | en |
| NO_COLOR | Disable colored output | (unset) |

---

## Stdin/Stdout Behavior

| Scenario | stdin | stdout | stderr |
|----------|-------|--------|--------|
| Normal conversion | Not used | Progress/result | Errors, warnings |
| Stdin input (`-`) | Content | Progress/result | Errors, warnings |
| JSON format | Same | JSON result only | Nothing (errors in JSON) |
| Piped output | Same | Minimal | Progress, errors |

---

## Examples

### Basic Usage

```bash
# Convert Markdown to EPUB
toepub convert README.md

# Convert HTML with metadata
toepub convert article.html --title "Article Title" --author "Jane Doe"

# Convert PDF
toepub convert document.pdf -o output.epub
```

### Multi-File Book

```bash
# Combine chapters into single book
toepub convert \
  --title "My Book" \
  --author "Author Name" \
  --cover cover.jpg \
  --language en-US \
  intro.md \
  chapters/ \
  appendix.md \
  -o mybook.epub
```

### Scripting

```bash
# Check conversion success
if toepub convert doc.md --format json | jq -e '.success'; then
  echo "Conversion succeeded"
else
  echo "Conversion failed"
fi

# Process multiple files
for f in *.md; do
  toepub convert "$f" --format json
done
```

### Pipeline

```bash
# Convert from curl
curl -s https://example.com/doc.html | toepub convert - --input-format html -o doc.epub

# Convert generated markdown
pandoc README.rst -t markdown | toepub convert - -o readme.epub
```
