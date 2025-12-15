# Implementation Plan: EPUB Converter CLI

**Branch**: `001-epub-converter-cli` | **Date**: 2025-12-15 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `specs/001-epub-converter-cli/spec.md`

## Summary

Build a CLI tool in Go that converts Markdown, HTML, and PDF files to EPUB 3+ format. The tool prioritizes ease of use (simple CLI interface), high output quality (epubcheck-validated EPUBs), and maintainability (modular architecture with clear separation between parsers, EPUB generator, and CLI layer).

## Technical Context

**Language/Version**: Go 1.24+
**Primary Dependencies**:
- CLI: cobra (CLI framework) + pflag (flag parsing)
- Markdown: goldmark (GFM-compliant parser)
- HTML: golang.org/x/net/html (HTML parsing)
- PDF: pdfcpu or unipdf (PDF text extraction)
- EPUB: go-epub or custom EPUB 3 generator
- Images: image/png, image/jpeg, image/gif (stdlib) + golang.org/x/image/webp

**Storage**: N/A (file-based I/O only)
**Testing**: Go testing package + testify for assertions
**Target Platform**: Cross-platform (Linux, macOS, Windows) - single binary distribution
**Project Type**: Single CLI application
**Performance Goals**: <5s for 50-page Markdown/HTML, <30s for 50-page PDF
**Constraints**: <100MB memory for typical documents, streaming for large files
**Scale/Scope**: Single-user CLI tool, files up to 100MB

## Ground-rules Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Evidence |
|-----------|--------|----------|
| I. EPUB 3+ Standards Compliance | ✅ PASS | FR-005 requires epubcheck validation; design uses EPUB 3 semantic markup |
| II. CLI-First Design | ✅ PASS | FR-010-014 define CLI behavior; cobra provides --help, exit codes |
| III. Input Format Fidelity | ✅ PASS | FR-001-003, FR-006-008 specify preservation requirements; goldmark for GFM |
| IV. Maintainability-First Architecture | ✅ PASS | Modular structure: parsers/, epub/, cli/ with clear interfaces |
| V. Test-Driven Quality | ✅ PASS | Unit tests per parser, integration tests with epubcheck validation |

**Quality Standards Alignment:**
- Output Quality: EPUB 3+ with proper nav.xhtml, UTF-8 throughout
- Performance: Targets align with SC-001 to SC-003
- Error Handling: FR-010, FR-011, FR-021 define error behaviors

## Project Structure

### Documentation (this feature)

```text
specs/001-epub-converter-cli/
├── design.md            # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (CLI contract)
└── tasks.md             # Phase 2 output (/rainbow.taskify)
```

### Source Code (repository root)

```text
cmd/
└── toepub/
    └── main.go          # Entry point

internal/
├── cli/
│   ├── root.go          # Root command, global flags
│   ├── convert.go       # Convert command implementation
│   └── output.go        # JSON/human-readable output formatting
├── parser/
│   ├── parser.go        # Parser interface
│   ├── markdown.go      # Markdown → intermediate representation
│   ├── html.go          # HTML → intermediate representation
│   └── pdf.go           # PDF → intermediate representation
├── epub/
│   ├── builder.go       # EPUB package builder
│   ├── content.go       # XHTML content document generation
│   ├── navigation.go    # nav.xhtml (TOC, landmarks)
│   ├── package.go       # OPF package document
│   └── metadata.go      # Dublin Core metadata handling
├── converter/
│   ├── converter.go     # Orchestrates parser → EPUB pipeline
│   └── image.go         # Image handling (embed, convert WebP)
└── model/
    ├── document.go      # Intermediate document representation
    ├── toc.go           # Table of contents structure
    └── metadata.go      # Book metadata structure

tests/
├── unit/
│   ├── parser/
│   ├── epub/
│   └── converter/
├── integration/
│   └── convert_test.go  # End-to-end conversion tests
├── contract/
│   └── epubcheck_test.go # EPUB validation tests
└── fixtures/
    ├── markdown/
    ├── html/
    └── pdf/
```

**Structure Decision**: Single project structure using Go's `cmd/` + `internal/` convention. The `internal/` package prevents external imports and enforces module boundaries. Each parser is isolated; the converter orchestrates the pipeline; EPUB generation is self-contained.

## Complexity Tracking

No ground-rules violations requiring justification.
