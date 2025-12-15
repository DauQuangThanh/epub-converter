# to-epub Development Guidelines

Auto-generated from feature plans and standards. Last updated: 2025-12-15

## Project Overview

CLI tool to convert Markdown, HTML, and PDF files to EPUB 3+ format.

## Active Technologies
- N/A (file-based I/O only) (001-epub-converter-cli)

- **Language**: Go 1.24+
- **CLI Framework**: cobra + pflag
- **Markdown Parser**: goldmark (GFM-compliant)
- **HTML Parser**: golang.org/x/net/html
- **PDF Library**: ledongthuc/pdf + pdfcpu
- **Testing**: Go testing package + testify
- **Linting**: golangci-lint

## Project Structure

```text
cmd/toepub/main.go       # Entry point
internal/
├── cli/                 # CLI commands (root, convert, version)
├── parser/              # Format parsers (markdown, html, pdf)
├── epub/                # EPUB generation (builder, content, navigation)
├── converter/           # Orchestration (converter, image handling)
└── model/               # Data structures (document, toc, metadata)
tests/fixtures/          # Test input files
docs/                    # Documentation and standards
specs/                   # Feature specifications
```

## Code Style (Go)

- **Naming**: MixedCaps (exported = PascalCase, unexported = camelCase)
- **Constants**: MixedCaps (not SCREAMING_SNAKE_CASE)
- **Errors**: `Err` prefix for variables, `Error` suffix for types
- **Interfaces**: `-er` suffix for single-method interfaces
- **Packages**: short, lowercase, singular nouns
- **Receivers**: single letter or short abbreviation (`d *Document`)

## Commands

```bash
# Build
go build -o toepub ./cmd/toepub

# Test
go test ./...
go test -race -coverprofile=coverage.out ./...

# Lint
golangci-lint run

# Format
gofmt -w .
goimports -w .
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 64 | File not found |
| 65 | Format error |
| 66 | Not writable |
| 70 | Internal error |

## Key Documents

- `docs/standards.md` - Full coding standards
- `docs/standards-cheatsheet.md` - Quick reference
- `memory/ground-rules.md` - Core principles
- `specs/001-epub-converter-cli/` - Feature specification

## Ground Rules Summary

1. **EPUB 3+ Compliance**: All output must pass epubcheck
2. **CLI-First Design**: Full functionality via CLI, proper exit codes
3. **Input Format Fidelity**: Preserve semantic meaning from sources
4. **Maintainability Architecture**: Single responsibility, no circular deps
5. **Test-Driven Quality**: Unit, integration, and contract tests required

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->

## Recent Changes
- 001-epub-converter-cli: Added Go 1.24+
