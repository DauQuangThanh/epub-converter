# Coding Standards & Conventions: to-epub

**Version**: 1.0 | **Date**: 2025-12-15 | **Status**: Active
**Maintained by**: Development Team | **Last Reviewed**: 2025-12-15

**Note**: This document defines the coding standards, naming conventions, and best practices for the to-epub EPUB converter CLI. All developers must follow these standards to ensure consistency, maintainability, and quality across the codebase.

---

## Document Control

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-12-15 | Development Team | Initial standards document |
| | | | |

**Related Documents**:

- Ground Rules: `memory/ground-rules.md`
- Feature Specifications: `specs/001-epub-converter-cli/spec.md`
- Implementation Plan: `specs/001-epub-converter-cli/design.md`

---

## Table of Contents

1. [Introduction](#1-introduction)
2. [UI Naming Conventions](#2-ui-naming-conventions) (N/A - CLI Tool)
3. [Code Naming Conventions](#3-code-naming-conventions)
4. [File and Directory Structure](#4-file-and-directory-structure)
5. [API Design Standards](#5-api-design-standards)
6. [Database Standards](#6-database-standards) (N/A)
7. [Testing Standards](#7-testing-standards)
8. [Git Workflow](#8-git-workflow)
9. [Documentation Standards](#9-documentation-standards)
10. [Code Style Guide](#10-code-style-guide)
11. [Enforcement](#11-enforcement)
12. [Appendices](#12-appendices)

---

## 1. Introduction

### 1.1 Purpose

This document establishes comprehensive coding standards and naming conventions for the to-epub EPUB converter CLI tool. Following these standards ensures:

- **Consistency**: Code looks uniform across the codebase
- **Maintainability**: Code is easier to understand and modify
- **Collaboration**: Team members can read and work with each other's code
- **Quality**: Automated tools can enforce standards
- **Onboarding**: New team members can quickly understand conventions

### 1.2 Scope

These standards apply to:

- All Go source code in the repository
- All test files and fixtures
- All documentation
- All configuration files
- CLI interface design

### 1.3 Technology Stack

**Language**: Go 1.24+
**CLI Framework**: cobra + pflag
**Markdown Parser**: goldmark (GFM-compliant)
**HTML Parser**: golang.org/x/net/html
**PDF Library**: ledongthuc/pdf + pdfcpu
**Image Processing**: image/png, image/jpeg, image/gif (stdlib) + golang.org/x/image/webp
**Testing**: Go testing package + testify
**Linting**: golangci-lint
**Formatting**: gofmt / goimports

### 1.4 How to Use This Document

- **Developers**: Follow these standards in all code you write
- **Code Reviewers**: Verify adherence to these standards in PRs
- **Team Leads**: Enforce standards and update this document as needed
- **New Team Members**: Read this document during onboarding

---

## 2. UI Naming Conventions

**N/A - No UI Layer**

This project is a command-line interface (CLI) tool with no graphical user interface. Skip to Section 3 for code naming conventions.

For CLI output formatting conventions, see [Section 5.2 CLI Output Standards](#52-cli-output-standards).

---

## 3. Code Naming Conventions

### 3.1 Variables

#### 3.1.1 Local Variables

**Convention**: camelCase, descriptive nouns or noun phrases

```go
// ✅ Good - camelCase, descriptive
userName := "John"
totalAmount := 100.50
isActive := true
itemsList := []Item{}
currentIndex := 0
documentContent := ""
chapterTitle := ""

// ✅ Good - short names acceptable in small scope
for i, item := range items {
    // i is acceptable for loop index
}

// ❌ Bad - unclear, abbreviated
usr := "John"
amt := 100.50
active := true  // Missing 'is' prefix for boolean
l := []Item{}   // Too short outside loop
```

#### 3.1.2 Constants

**Convention**: Go uses MixedCaps for exported constants, mixedCaps for unexported

```go
// ✅ Good - exported constants (PascalCase)
const (
    MaxRetryCount    = 3
    DefaultTimeout   = 30 * time.Second
    EPUBVersion      = "3.0"
    MIMETypeXHTML    = "application/xhtml+xml"
    MIMETypePNG      = "image/png"
    MIMETypeJPEG     = "image/jpeg"
)

// ✅ Good - unexported constants (camelCase)
const (
    maxFileSize      = 100 * 1024 * 1024  // 100MB
    defaultLanguage  = "en"
    outputExtension  = ".epub"
)

// ✅ Good - grouped by purpose
const (
    // Exit codes
    ExitSuccess       = 0
    ExitGeneralError  = 1
    ExitInvalidArgs   = 2
    ExitFileNotFound  = 64
    ExitFormatError   = 65
    ExitNotWritable   = 66
    ExitInternalError = 70
)

// ❌ Bad - SCREAMING_SNAKE_CASE (not Go convention)
const MAX_RETRY_COUNT = 3
const DEFAULT_TIMEOUT = 30
```

#### 3.1.3 Package-Level Variables

**Convention**: Avoid package-level variables; use dependency injection instead

```go
// ✅ Good - pass dependencies explicitly
type Converter struct {
    parser Parser
    builder EPUBBuilder
    logger  *log.Logger
}

func NewConverter(parser Parser, builder EPUBBuilder, logger *log.Logger) *Converter {
    return &Converter{
        parser:  parser,
        builder: builder,
        logger:  logger,
    }
}

// ⚠️ Acceptable - package-level for initialization
var (
    version   = "dev"     // Set via ldflags at build time
    buildDate = "unknown"
)

// ❌ Bad - mutable package-level state
var globalConfig Config
var currentDocument *Document
```

### 3.2 Functions and Methods

#### 3.2.1 Function Names

**Convention**: MixedCaps, verb-based for actions, noun-based for getters

```go
// ✅ Good - verb + noun, descriptive
func ParseMarkdown(content []byte) (*Document, error) { }
func BuildEPUB(doc *Document) ([]byte, error) { }
func ValidateMetadata(meta Metadata) error { }
func FormatOutput(result *Result, format string) string { }
func WriteToFile(path string, data []byte) error { }

// ✅ Good - getters use noun form (no Get prefix in Go)
func (m *Metadata) Title() string { return m.title }
func (m *Metadata) Authors() []string { return m.authors }

// ✅ Good - exported functions (PascalCase)
func Convert(input string, options Options) (*Result, error) { }

// ✅ Good - unexported functions (camelCase)
func parseHeadings(content string) []Heading { }
func buildNavDocument(toc *TOC) string { }

// ❌ Bad - noun-based actions, Get prefix
func Total(items []Item) float64 { }        // Missing verb
func GetTitle() string { }                   // Go doesn't use Get prefix
func DoConversion() { }                      // Too generic
```

#### 3.2.2 Boolean Functions

**Convention**: Prefix with Is, Has, Should, Can, or use adjective form

```go
// ✅ Good - boolean intent clear
func IsValid(data []byte) bool { }
func HasPermission(user *User, action string) bool { }
func ShouldRetry(err error) bool { }
func CanConvert(format string) bool { }

// ✅ Good - adjective form for methods
func (d *Document) Empty() bool { }
func (m *Metadata) Complete() bool { }
func (p *Path) Writable() bool { }

// ❌ Bad - unclear return type
func Valid(data []byte) bool { }        // Missing prefix
func CheckPermission() bool { }         // Check is action, not boolean
func Retry(count int) bool { }          // Unclear intent
```

#### 3.2.3 Method Receivers

**Convention**: Single-letter or short abbreviation of type name

```go
// ✅ Good - consistent short receiver names
func (d *Document) AddChapter(chapter Chapter) { }
func (d *Document) Validate() error { }

func (c *Converter) Convert(input string) (*Result, error) { }
func (c *Converter) SetOptions(opts Options) { }

func (m *Metadata) SetTitle(title string) { }
func (m *Metadata) Title() string { }

// ✅ Good - longer names for clarity when needed
func (mb *MetadataBuilder) Build() *Metadata { }
func (ep *EPUBPackager) Package() error { }

// ❌ Bad - inconsistent or verbose receivers
func (doc *Document) AddChapter() { }       // Inconsistent with 'd'
func (self *Converter) Convert() { }        // Don't use 'self' or 'this'
func (converter *Converter) Convert() { }   // Too verbose
```

### 3.3 Types

#### 3.3.1 Struct Names

**Convention**: PascalCase, nouns, exported types for public API

```go
// ✅ Good - PascalCase, descriptive nouns
type Document struct {
    Metadata  Metadata
    Chapters  []Chapter
    Resources []Resource
    TOC       TableOfContents
}

type Chapter struct {
    ID       string
    Title    string
    Level    int
    Content  string
    FileName string
    Order    int
}

type ConversionResult struct {
    Success    bool
    OutputPath string
    Warnings   []string
    Error      error
    Stats      ConversionStats
}

// ❌ Bad - camelCase, abbreviations
type document struct { }     // Unexported when should be exported
type Chap struct { }         // Abbreviation
type ConvRes struct { }      // Unclear abbreviation
```

#### 3.3.2 Interface Names

**Convention**: Verb-er suffix for single-method interfaces, descriptive for multi-method

```go
// ✅ Good - single-method interfaces with -er suffix
type Parser interface {
    Parse(content []byte) (*Document, error)
}

type Builder interface {
    Build(doc *Document) ([]byte, error)
}

type Validator interface {
    Validate() error
}

type Reader interface {
    Read(p []byte) (n int, err error)
}

// ✅ Good - multi-method interfaces with descriptive names
type EPUBGenerator interface {
    SetMetadata(meta Metadata)
    AddChapter(chapter Chapter)
    AddResource(resource Resource)
    Generate() ([]byte, error)
}

type FormatConverter interface {
    CanConvert(format string) bool
    Convert(input []byte) (*Document, error)
    SupportedFormats() []string
}

// ❌ Bad - 'I' prefix, generic names
type IParser interface { }     // Don't use I prefix
type Doer interface { }        // Too generic
```

#### 3.3.3 Type Aliases and Custom Types

```go
// ✅ Good - meaningful type aliases
type MediaType string
type ExitCode int
type ChapterID string

const (
    MediaTypeXHTML MediaType = "application/xhtml+xml"
    MediaTypePNG   MediaType = "image/png"
    MediaTypeJPEG  MediaType = "image/jpeg"
)

// ✅ Good - functional types
type ParserFunc func([]byte) (*Document, error)
type ValidatorFunc func(*Document) error
type OutputFormatter func(*Result) string
```

### 3.4 Packages

**Convention**: Short, lowercase, singular nouns

```go
// ✅ Good - short, clear package names
package parser     // Parsing logic
package epub       // EPUB generation
package cli        // CLI command handling
package model      // Data models
package converter  // Conversion orchestration

// ✅ Good - import paths reflect structure
import (
    "github.com/dauquangthanh/epub-converter/internal/parser"
    "github.com/dauquangthanh/epub-converter/internal/epub"
    "github.com/dauquangthanh/epub-converter/internal/cli"
)

// ❌ Bad - plural, multi-word, generic
package parsers    // Use singular
package epub_builder // Use single word or combine
package util       // Too generic - be specific
package common     // Too generic - split into specific packages
```

### 3.5 Error Variables and Types

**Convention**: Err prefix for error variables, Error suffix for error types

```go
// ✅ Good - error variables with Err prefix
var (
    ErrFileNotFound     = errors.New("file not found")
    ErrInvalidFormat    = errors.New("invalid input format")
    ErrOutputNotWritable = errors.New("output path not writable")
    ErrParseFailed      = errors.New("failed to parse input")
    ErrEPUBGeneration   = errors.New("failed to generate EPUB")
)

// ✅ Good - custom error types with Error suffix
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

type ParseError struct {
    Line    int
    Column  int
    Message string
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("parse error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}
```

---

## 4. File and Directory Structure

### 4.1 File Naming

#### 4.1.1 Source Code Files

**Convention**: lowercase, snake_case for multi-word names

```
# ✅ Good - lowercase, descriptive
parser.go
markdown.go
html_parser.go
pdf_extractor.go
epub_builder.go
content_document.go
navigation.go
metadata.go

# ❌ Bad - camelCase, PascalCase
markdownParser.go
MarkdownParser.go
epub-builder.go
```

#### 4.1.2 Test Files

**Convention**: `{source}_test.go` adjacent to source file

```
internal/parser/
├── parser.go
├── parser_test.go
├── markdown.go
├── markdown_test.go
├── html.go
├── html_test.go
├── pdf.go
└── pdf_test.go
```

#### 4.1.3 Special Files

```
# Root level
go.mod              # Module definition
go.sum              # Dependency checksums
Makefile            # Build commands
README.md           # Project documentation
LICENSE             # License file
.gitignore          # Git ignore rules
.golangci.yml       # Linter configuration
.editorconfig       # Editor settings
CLAUDE.md           # AI agent context
```

### 4.2 Directory Structure

#### 4.2.1 Project Root Structure

```
toepub/
├── cmd/
│   └── toepub/
│       └── main.go              # Entry point only
├── internal/                    # Private application code
│   ├── cli/                     # CLI command handling
│   │   ├── root.go              # Root command, global flags
│   │   ├── convert.go           # Convert command
│   │   ├── version.go           # Version command
│   │   └── output.go            # Output formatting
│   ├── parser/                  # Input format parsers
│   │   ├── parser.go            # Parser interface
│   │   ├── markdown.go          # Markdown parser
│   │   ├── html.go              # HTML parser
│   │   └── pdf.go               # PDF parser
│   ├── epub/                    # EPUB generation
│   │   ├── builder.go           # EPUB package builder
│   │   ├── content.go           # XHTML content generation
│   │   ├── navigation.go        # nav.xhtml generation
│   │   ├── package.go           # OPF package document
│   │   └── metadata.go          # Dublin Core metadata
│   ├── converter/               # Conversion orchestration
│   │   ├── converter.go         # Main conversion logic
│   │   └── image.go             # Image handling
│   └── model/                   # Data structures
│       ├── document.go          # Intermediate representation
│       ├── toc.go               # Table of contents
│       └── metadata.go          # Book metadata
├── tests/                       # Test resources
│   ├── fixtures/                # Test input files
│   │   ├── markdown/
│   │   ├── html/
│   │   └── pdf/
│   └── golden/                  # Expected output files
├── docs/                        # Documentation
│   ├── standards.md             # This file
│   └── examples/
├── specs/                       # Feature specifications
│   └── 001-epub-converter-cli/
├── memory/                      # Project memory
│   └── ground-rules.md
├── .rainbow/                    # Rainbow framework
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── CLAUDE.md
```

#### 4.2.2 Package Organization Rules

1. **cmd/**: Contains only `main.go` files that wire up the application
2. **internal/**: All private application code - prevents external imports
3. **One package = One concern**: Each package handles a single responsibility
4. **No circular imports**: Design packages to avoid circular dependencies

```
# Dependency direction (allowed imports)
cmd/toepub → internal/cli
internal/cli → internal/converter, internal/model
internal/converter → internal/parser, internal/epub, internal/model
internal/parser → internal/model
internal/epub → internal/model
internal/model → (no internal dependencies)
```

### 4.3 Module Organization

**Convention**: Group imports in three sections

```go
import (
    // Standard library
    "context"
    "errors"
    "fmt"
    "io"
    "os"

    // Third-party packages
    "github.com/spf13/cobra"
    "github.com/yuin/goldmark"

    // Internal packages
    "github.com/dauquangthanh/epub-converter/internal/model"
    "github.com/dauquangthanh/epub-converter/internal/parser"
)
```

---

## 5. API Design Standards

### 5.1 Internal Package APIs

#### 5.1.1 Constructor Functions

**Convention**: New prefix, return pointer for structs

```go
// ✅ Good - New prefix, returns pointer
func NewConverter(opts ConverterOptions) *Converter {
    return &Converter{
        parser:  opts.Parser,
        builder: opts.Builder,
        logger:  opts.Logger,
    }
}

func NewMarkdownParser(opts MarkdownOptions) *MarkdownParser {
    return &MarkdownParser{
        extensions: opts.Extensions,
    }
}

// ✅ Good - functional options pattern for complex construction
func NewEPUBBuilder(opts ...BuilderOption) *EPUBBuilder {
    b := &EPUBBuilder{
        version: "3.0",
    }
    for _, opt := range opts {
        opt(b)
    }
    return b
}

type BuilderOption func(*EPUBBuilder)

func WithVersion(v string) BuilderOption {
    return func(b *EPUBBuilder) {
        b.version = v
    }
}
```

#### 5.1.2 Method Signatures

**Convention**: Accept interfaces, return concrete types

```go
// ✅ Good - accept interface, return concrete
func (c *Converter) Convert(r io.Reader) (*Result, error) {
    // ...
}

// ✅ Good - use small interfaces
func ParseFrom(r io.Reader) (*Document, error) {
    // ...
}

// ✅ Good - error as last return value
func BuildEPUB(doc *Document) ([]byte, error) {
    // ...
}

// ❌ Bad - error not last
func BuildEPUB(doc *Document) (error, []byte) {
    // ...
}
```

### 5.2 CLI Output Standards

#### 5.2.1 Human-Readable Output

**Convention**: Structured, informative output with visual indicators

```
# Success output
Converting document.md...
✓ Created document.epub (145 KB)
  - 5 chapters
  - 12 images
  - Duration: 1.2s

# Warning output
Converting document.md...
⚠ Warning: Image not found: missing.png
⚠ Warning: Unsupported image format: photo.webp (converted to PNG)
✓ Created document.epub (145 KB)
  - 5 chapters
  - 11 images (1 converted)
  - Duration: 1.5s

# Error output (stderr)
Error: Cannot read input file: document.md
  → File does not exist

Suggestion: Check the file path and try again.
```

#### 5.2.2 JSON Output

**Convention**: Consistent JSON structure with success, data, error patterns

```json
// Success response
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

// Error response
{
  "success": false,
  "error": {
    "code": 64,
    "message": "Cannot read input file: document.md",
    "detail": "File does not exist"
  }
}
```

### 5.3 Exit Codes

**Convention**: Follow BSD/sysexits.h conventions

| Code | Constant | Meaning |
|------|----------|---------|
| 0 | ExitSuccess | Success |
| 1 | ExitGeneralError | General error |
| 2 | ExitInvalidArgs | Invalid arguments or flags |
| 64 | ExitFileNotFound | Input file not found |
| 65 | ExitFormatError | Input format not supported |
| 66 | ExitNotWritable | Output path not writable |
| 70 | ExitInternalError | Internal software error |

---

## 6. Database Standards

**N/A - No Database**

This project is a file-based CLI tool with no database persistence. All data is read from input files and written to output EPUB files.

---

## 7. Testing Standards

### 7.1 Test File Organization

**Convention**: Test files adjacent to source, `_test.go` suffix

```
internal/parser/
├── markdown.go
├── markdown_test.go
├── html.go
├── html_test.go
└── testdata/           # Package-level test data
    ├── simple.md
    └── complex.md
```

### 7.2 Test Function Naming

**Convention**: `Test{Function}_{Scenario}` for clarity

```go
// ✅ Good - descriptive test names
func TestParseMarkdown_WithValidInput(t *testing.T) { }
func TestParseMarkdown_WithEmptyContent(t *testing.T) { }
func TestParseMarkdown_WithInvalidUTF8(t *testing.T) { }
func TestParseMarkdown_WithFrontMatter(t *testing.T) { }

func TestConvert_MarkdownToEPUB_Success(t *testing.T) { }
func TestConvert_HTMLToEPUB_WithImages(t *testing.T) { }
func TestConvert_PDFToEPUB_ImageBasedPDF_ReturnsError(t *testing.T) { }

// ✅ Good - table-driven tests
func TestParseHeadingLevel(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected int
    }{
        {"h1 heading", "# Title", 1},
        {"h2 heading", "## Subtitle", 2},
        {"h6 heading", "###### Deep", 6},
        {"not a heading", "Regular text", 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := parseHeadingLevel(tt.input)
            if result != tt.expected {
                t.Errorf("parseHeadingLevel(%q) = %d, want %d", tt.input, result, tt.expected)
            }
        })
    }
}

// ❌ Bad - unclear test names
func TestParse1(t *testing.T) { }
func TestMarkdown(t *testing.T) { }
func Test_convert(t *testing.T) { }
```

### 7.3 Test Structure

**Convention**: Arrange-Act-Assert (AAA) pattern

```go
func TestConverter_Convert_WithValidMarkdown(t *testing.T) {
    // Arrange
    input := []byte("# Hello World\n\nThis is a test.")
    converter := NewConverter(DefaultOptions())

    // Act
    result, err := converter.Convert(input)

    // Assert
    if err != nil {
        t.Fatalf("Convert() unexpected error: %v", err)
    }
    if !result.Success {
        t.Errorf("Convert() success = false, want true")
    }
    if result.Stats.Chapters != 1 {
        t.Errorf("Convert() chapters = %d, want 1", result.Stats.Chapters)
    }
}
```

### 7.4 Test Categories

```go
// Unit tests - test individual functions
func TestParseMarkdown_ExtractsHeadings(t *testing.T) { }

// Integration tests - test component interactions
func TestConverter_FullPipeline_MarkdownToEPUB(t *testing.T) { }

// Contract tests - verify EPUB validity
func TestGeneratedEPUB_PassesEpubcheck(t *testing.T) {
    // Generate EPUB
    epub := generateTestEPUB()

    // Run epubcheck validation
    err := runEpubcheck(epub)
    if err != nil {
        t.Errorf("Generated EPUB failed epubcheck: %v", err)
    }
}
```

### 7.5 Test Fixtures

**Convention**: Use `testdata/` directory within package

```
internal/parser/testdata/
├── markdown/
│   ├── simple.md           # Basic markdown
│   ├── gfm_tables.md       # GFM table syntax
│   ├── front_matter.md     # YAML front matter
│   └── images/             # Image references
├── html/
│   ├── simple.html
│   └── with_css.html
└── expected/
    └── simple_toc.json     # Expected TOC structure
```

### 7.6 Test Helpers

**Convention**: Prefix with `test` or use `t.Helper()`

```go
// ✅ Good - helper function marked with t.Helper()
func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func loadTestFixture(t *testing.T, name string) []byte {
    t.Helper()
    data, err := os.ReadFile(filepath.Join("testdata", name))
    if err != nil {
        t.Fatalf("failed to load fixture %s: %v", name, err)
    }
    return data
}

// ✅ Good - test setup function
func setupTestConverter(t *testing.T) *Converter {
    t.Helper()
    return NewConverter(ConverterOptions{
        Logger: log.New(io.Discard, "", 0),
    })
}
```

---

## 8. Git Workflow

### 8.1 Branch Naming

**Convention**: `{type}/{feature-id}-{brief-description}`

```bash
# ✅ Good - structured branch names
feature/001-epub-converter-cli
feature/002-add-pdf-support
bugfix/001-fix-image-embedding
hotfix/security-patch-v1.0.1
release/v1.0.0
chore/update-dependencies

# ❌ Bad - unclear branches
fix
new-feature
john-work
temp
```

**Branch Types**:

- `feature/` - New features (use feature ID from specs/)
- `bugfix/` - Bug fixes
- `hotfix/` - Critical production fixes
- `release/` - Release preparation
- `chore/` - Maintenance tasks (deps, CI, docs)

### 8.2 Commit Messages

**Convention**: Conventional Commits format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types**:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no logic changes)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `perf`: Performance improvements

**Scopes** (for this project):

- `cli`: CLI command handling
- `parser`: Input parsers (markdown, html, pdf)
- `epub`: EPUB generation
- `converter`: Conversion orchestration
- `model`: Data models
- `deps`: Dependencies

**Examples**:

```bash
# ✅ Good - clear, conventional commits
feat(parser): add GFM table parsing support
fix(epub): resolve image path resolution in subdirectories
docs(readme): add installation instructions for macOS
refactor(converter): extract image handling to separate module
test(parser): add edge case tests for malformed markdown
chore(deps): update goldmark to v1.6.0
perf(epub): optimize ZIP compression for large images

# With body and footer
feat(cli): add --cover flag for custom cover image

Allows users to specify a cover image via CLI flag.
The image is validated for supported formats (JPEG, PNG)
before being embedded in the EPUB.

Closes #45
```

### 8.3 Pull Request Standards

**Convention**: Descriptive title + template body

```markdown
## Summary
Brief description of changes (1-3 bullet points)

## Test Plan
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] Generated EPUBs pass epubcheck

## Checklist
- [ ] Code follows project standards
- [ ] Tests added for new functionality
- [ ] Documentation updated if needed
- [ ] No breaking changes (or documented)
```

---

## 9. Documentation Standards

### 9.1 Code Comments

**Convention**: Use comments to explain "why", not "what"

```go
// ✅ Good - explains why
func (p *MarkdownParser) Parse(content []byte) (*Document, error) {
    // Use goldmark with GFM extension to support tables and task lists
    // as specified in ground-rules (III. Input Format Fidelity)
    md := goldmark.New(
        goldmark.WithExtensions(extension.GFM),
    )

    // ...
}

// ✅ Good - explains non-obvious behavior
func normalizeImagePath(src string, basePath string) string {
    // Handle both Unix and Windows path separators
    // EPUB spec requires forward slashes in paths
    normalized := filepath.ToSlash(src)

    // ...
}

// ❌ Bad - states the obvious
func (p *MarkdownParser) Parse(content []byte) (*Document, error) {
    // Create a new goldmark instance
    md := goldmark.New()

    // Parse the content
    doc := md.Parse(content)

    // Return the document
    return doc, nil
}
```

### 9.2 Package Documentation

**Convention**: Package comment in `doc.go` or first file alphabetically

```go
// Package parser provides input format parsers for the EPUB converter.
//
// The parser package implements parsers for Markdown, HTML, and PDF formats.
// Each parser converts input content into an intermediate Document representation
// that can be processed by the EPUB generator.
//
// Example usage:
//
//     parser := parser.NewMarkdownParser(opts)
//     doc, err := parser.Parse(content)
//     if err != nil {
//         log.Fatal(err)
//     }
//
// Parsers implement the Parser interface, allowing for consistent usage
// across different input formats.
package parser
```

### 9.3 Function Documentation

**Convention**: Document all exported functions

```go
// Parse converts Markdown content into a Document representation.
//
// The parser supports GitHub Flavored Markdown (GFM) including:
//   - Tables
//   - Task lists
//   - Strikethrough
//   - Autolinks
//
// YAML front matter is extracted and used to populate document metadata.
// Relative image paths are resolved against the provided base path.
//
// Returns an error if the content is not valid UTF-8 or if parsing fails.
func (p *MarkdownParser) Parse(content []byte) (*Document, error) {
    // ...
}
```

### 9.4 README Structure

```markdown
# toepub

Convert Markdown, HTML, and PDF to EPUB 3+ format.

## Features

- Markdown to EPUB (GFM support)
- HTML to EPUB (CSS preservation)
- PDF to EPUB (text extraction)

## Installation

### From Source
\`\`\`bash
go install github.com/dauquangthanh/epub-converter/cmd/toepub@latest
\`\`\`

### Pre-built Binaries
Download from [Releases](releases).

## Usage

\`\`\`bash
# Basic conversion
toepub convert document.md

# With metadata
toepub convert document.md --title "My Book" --author "John Doe"
\`\`\`

## Documentation

- [CLI Reference](docs/cli-reference.md)
- [Examples](docs/examples/)

## Development

\`\`\`bash
make build
make test
make lint
\`\`\`

## License

MIT
```

---

## 10. Code Style Guide

### 10.1 Formatting

**Convention**: Use `gofmt` / `goimports` (enforced by CI)

```go
// ✅ Good - gofmt formatted
func (c *Converter) Convert(input string) (*Result, error) {
    doc, err := c.parser.Parse(input)
    if err != nil {
        return nil, fmt.Errorf("parse failed: %w", err)
    }

    epub, err := c.builder.Build(doc)
    if err != nil {
        return nil, fmt.Errorf("build failed: %w", err)
    }

    return &Result{
        Success: true,
        Data:    epub,
    }, nil
}
```

### 10.2 Line Length

**Convention**: No strict limit, but prefer ~100 characters for readability

```go
// ✅ Good - break long function calls
result, err := converter.ConvertWithOptions(
    inputPath,
    outputPath,
    ConverterOptions{
        Title:    title,
        Author:   author,
        Language: language,
    },
)

// ✅ Good - break long conditionals
if err != nil &&
    !errors.Is(err, ErrFileNotFound) &&
    !errors.Is(err, ErrInvalidFormat) {
    return err
}
```

### 10.3 Error Handling

**Convention**: Handle errors explicitly, wrap with context

```go
// ✅ Good - wrap errors with context
func (c *Converter) Convert(inputPath string) (*Result, error) {
    content, err := os.ReadFile(inputPath)
    if err != nil {
        return nil, fmt.Errorf("reading input file: %w", err)
    }

    doc, err := c.parser.Parse(content)
    if err != nil {
        return nil, fmt.Errorf("parsing content: %w", err)
    }

    epub, err := c.builder.Build(doc)
    if err != nil {
        return nil, fmt.Errorf("building EPUB: %w", err)
    }

    return &Result{Success: true, Data: epub}, nil
}

// ✅ Good - use errors.Is for comparison
if errors.Is(err, ErrFileNotFound) {
    return ExitFileNotFound
}

// ❌ Bad - ignore errors
content, _ := os.ReadFile(inputPath)

// ❌ Bad - no context in error
if err != nil {
    return nil, err
}
```

### 10.4 Blank Lines

**Convention**: Use blank lines to separate logical blocks

```go
func (c *Converter) Convert(input string) (*Result, error) {
    // Validation
    if input == "" {
        return nil, ErrEmptyInput
    }

    // Parse input
    doc, err := c.parser.Parse(input)
    if err != nil {
        return nil, fmt.Errorf("parse failed: %w", err)
    }

    // Build EPUB
    epub, err := c.builder.Build(doc)
    if err != nil {
        return nil, fmt.Errorf("build failed: %w", err)
    }

    // Return result
    return &Result{Success: true, Data: epub}, nil
}
```

---

## 11. Enforcement

### 11.1 Automated Tools

**Linting**: golangci-lint

```yaml
# .golangci.yml
linters:
  enable:
    - errcheck      # Check error returns
    - govet         # Go vet checks
    - staticcheck   # Static analysis
    - unused        # Unused code detection
    - gosimple      # Code simplification
    - ineffassign   # Ineffectual assignments
    - misspell      # Spelling mistakes
    - gofmt         # Format checking
    - goimports     # Import organization

linters-settings:
  errcheck:
    check-type-assertions: true
  govet:
    check-shadowing: true

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
```

**Formatting**: gofmt / goimports (built into Go)

### 11.2 Pre-commit Hooks

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Format check
if ! gofmt -l . | grep -q '^'; then
    echo "Code is properly formatted"
else
    echo "Code needs formatting. Run: gofmt -w ."
    exit 1
fi

# Lint check
if golangci-lint run; then
    echo "Linting passed"
else
    echo "Linting failed"
    exit 1
fi

# Test check
if go test ./...; then
    echo "Tests passed"
else
    echo "Tests failed"
    exit 1
fi
```

### 11.3 CI/CD Integration

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Run tests
        run: go test -race -coverprofile=coverage.out ./...
      - name: Check coverage
        run: go tool cover -func=coverage.out

  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Build
        run: go build -v ./...
```

### 11.4 Code Review Checklist

- [ ] Follows Go naming conventions (exported/unexported)
- [ ] Uses gofmt formatting
- [ ] Passes golangci-lint checks
- [ ] Has tests with descriptive names
- [ ] Follows file and directory structure standards
- [ ] Errors are wrapped with context
- [ ] Has meaningful commit messages
- [ ] Documentation updated if needed

---

## 12. Appendices

### 12.1 Glossary

| Term | Definition |
|------|------------|
| PascalCase | Capitalized first letter of each word: `UserProfile` |
| camelCase | Lowercase first letter, capitalized subsequent words: `userProfile` |
| snake_case | Lowercase with underscores: `user_profile` |
| kebab-case | Lowercase with hyphens: `user-profile` |
| MixedCaps | Go's term for PascalCase/camelCase (no underscores) |
| Exported | PascalCase names visible outside package |
| Unexported | camelCase names private to package |

### 12.2 Quick Reference Checklist

**Code Naming**:

- [ ] Exported names: PascalCase
- [ ] Unexported names: camelCase
- [ ] Constants: MixedCaps (not SCREAMING_SNAKE)
- [ ] Errors: Err prefix for variables, Error suffix for types
- [ ] Interfaces: -er suffix for single-method
- [ ] Packages: short, lowercase, singular

**File Naming**:

- [ ] Source files: lowercase, snake_case.go
- [ ] Test files: {source}_test.go
- [ ] Test data: testdata/ directory

**Testing**:

- [ ] Test functions: Test{Function}_{Scenario}
- [ ] Table-driven tests for multiple cases
- [ ] AAA pattern (Arrange-Act-Assert)

**Git**:

- [ ] Branches: type/feature-id-description
- [ ] Commits: Conventional Commits format
- [ ] PRs: Descriptive titles

### 12.3 Tool Configuration Files

**Editor Config** (`.editorconfig`):

```ini
root = true

[*]
indent_style = tab
indent_size = 4
end_of_line = lf
charset = utf-8
trim_trailing_whitespace = true
insert_final_newline = true

[*.md]
trim_trailing_whitespace = false

[*.yml]
indent_style = space
indent_size = 2

[Makefile]
indent_style = tab
```

**golangci-lint** (`.golangci.yml`): See Section 11.1

### 12.4 Resources

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Naming Conventions](https://talks.golang.org/2014/names.slide)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [EPUB 3.3 Specification](https://www.w3.org/TR/epub-33/)

### 12.5 Change History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-12-15 | Development Team | Initial standards document |
| | | | |

---

**END OF STANDARDS DOCUMENT**

---

## Maintenance Notes

This document should be:

- **Reviewed quarterly** by the team
- **Updated** when new technologies are adopted
- **Referenced** in code reviews
- **Shared** with new team members during onboarding
- **Enforced** through automated tools (golangci-lint, gofmt)

For questions or suggestions, update this document via PR.
