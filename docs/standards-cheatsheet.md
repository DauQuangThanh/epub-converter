# to-epub Standards Cheatsheet

Quick reference for coding standards. See [standards.md](./standards.md) for full details.

## Go Naming

| Element | Convention | Example |
|---------|------------|---------|
| Exported | PascalCase | `ParseMarkdown`, `Document` |
| Unexported | camelCase | `parseHeadings`, `maxFileSize` |
| Constants | MixedCaps | `MaxRetryCount`, `EPUBVersion` |
| Errors (var) | Err prefix | `ErrFileNotFound` |
| Errors (type) | Error suffix | `ValidationError` |
| Interfaces | -er suffix | `Parser`, `Builder` |
| Packages | lowercase | `parser`, `epub`, `model` |
| Receivers | short | `(d *Document)`, `(c *Converter)` |

## Boolean Naming

```go
// Functions: Is, Has, Should, Can prefix
func IsValid(data []byte) bool
func HasPermission(user *User) bool
func ShouldRetry(err error) bool
func CanConvert(format string) bool

// Methods: adjective form
func (d *Document) Empty() bool
func (m *Metadata) Complete() bool
```

## File Naming

```
source_file.go          # lowercase, snake_case
source_file_test.go     # test file adjacent to source
testdata/               # test fixtures directory
```

## Project Structure

```
cmd/toepub/main.go      # Entry point
internal/cli/           # CLI commands
internal/parser/        # Format parsers
internal/epub/          # EPUB generation
internal/converter/     # Orchestration
internal/model/         # Data structures
```

## Imports

```go
import (
    // Standard library
    "fmt"
    "os"

    // Third-party
    "github.com/spf13/cobra"

    // Internal
    "github.com/dauquangthanh/epub-converter/internal/model"
)
```

## Error Handling

```go
// Always wrap with context
if err != nil {
    return nil, fmt.Errorf("parsing content: %w", err)
}

// Use errors.Is for comparison
if errors.Is(err, ErrFileNotFound) {
    return ExitFileNotFound
}
```

## Testing

```go
// Naming: Test{Function}_{Scenario}
func TestParseMarkdown_WithValidInput(t *testing.T) { }

// Pattern: Arrange-Act-Assert
func TestConvert_Success(t *testing.T) {
    // Arrange
    input := []byte("# Test")
    converter := NewConverter()

    // Act
    result, err := converter.Convert(input)

    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}
```

## Git

```bash
# Branches
feature/001-epub-converter-cli
bugfix/fix-image-embedding

# Commits (Conventional Commits)
feat(parser): add GFM table support
fix(epub): resolve image path issue
docs(readme): update installation
test(parser): add edge case tests
```

## Exit Codes

| Code | Constant | Meaning |
|------|----------|---------|
| 0 | ExitSuccess | Success |
| 1 | ExitGeneralError | General error |
| 2 | ExitInvalidArgs | Invalid arguments |
| 64 | ExitFileNotFound | File not found |
| 65 | ExitFormatError | Format error |
| 66 | ExitNotWritable | Not writable |
| 70 | ExitInternalError | Internal error |

## Tools

```bash
# Format
gofmt -w .
goimports -w .

# Lint
golangci-lint run

# Test
go test ./...
go test -race -coverprofile=coverage.out ./...
```
