// Package converter orchestrates the document conversion pipeline.
package converter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dauquangthanh/epub-converter/internal/epub"
	"github.com/dauquangthanh/epub-converter/internal/model"
	"github.com/dauquangthanh/epub-converter/internal/parser"
)

// Common errors
var (
	ErrNoInput         = errors.New("no input files specified")
	ErrFileNotFound    = errors.New("file not found")
	ErrUnsupportedFmt  = errors.New("unsupported input format")
	ErrOutputNotWrite  = errors.New("output path not writable")
	ErrConversionFailed = errors.New("conversion failed")
)

// Options configures the conversion process.
type Options struct {
	OutputPath  string          // Output EPUB file path
	InputFormat string          // Force input format (md, html, pdf)
	CLIMetadata *model.Metadata // Metadata overrides from CLI flags
}

// Converter orchestrates the document conversion pipeline.
type Converter struct {
	parsers  map[parser.Format]parser.Parser
	builder  *epub.Builder
	imgHandler *ImageHandler
}

// New creates a new Converter with default parsers.
func New() *Converter {
	c := &Converter{
		parsers:    make(map[parser.Format]parser.Parser),
		builder:    epub.NewBuilder(),
		imgHandler: NewImageHandler(),
	}

	// Register default parsers
	c.RegisterParser(parser.FormatMarkdown, parser.NewMarkdownParser())
	c.RegisterParser(parser.FormatHTML, parser.NewHTMLParser())
	c.RegisterParser(parser.FormatPDF, parser.NewPDFParser())

	return c
}

// RegisterParser adds a parser for a specific format.
func (c *Converter) RegisterParser(format parser.Format, p parser.Parser) {
	c.parsers[format] = p
}

// Convert converts input files to EPUB format.
func (c *Converter) Convert(inputs []string, opts Options) (*model.ConversionResult, error) {
	start := time.Now()
	result := &model.ConversionResult{
		Success:  false,
		Warnings: make([]string, 0),
	}

	if len(inputs) == 0 {
		return result, ErrNoInput
	}

	// Expand directories and validate inputs
	files, err := c.expandInputs(inputs)
	if err != nil {
		return result, err
	}

	if len(files) == 0 {
		return result, fmt.Errorf("%w: no supported files found", ErrNoInput)
	}

	// Detect format from first file if not specified
	format := c.detectFormat(files[0], opts.InputFormat)
	if format == parser.FormatUnknown {
		return result, fmt.Errorf("%w: cannot detect format for %s", ErrUnsupportedFmt, files[0])
	}

	// Get parser for format
	p := c.getParser(format)
	if p == nil {
		return result, fmt.Errorf("%w: no parser for format %s", ErrUnsupportedFmt, format)
	}

	// Parse all input files
	doc := model.NewDocument()
	for i, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return result, fmt.Errorf("reading %s: %w", file, err)
		}

		basePath := filepath.Dir(file)
		parsedDoc, err := p.Parse(content, basePath)
		if err != nil {
			return result, fmt.Errorf("parsing %s: %w", file, err)
		}

		// Merge parsed content into main document
		c.mergeDocument(doc, parsedDoc, i)
	}

	// Apply CLI metadata overrides
	if opts.CLIMetadata != nil {
		doc.Metadata.Merge(opts.CLIMetadata)
	}

	// Ensure document has a title
	if doc.Metadata.Title == "" {
		// Use first input file name as title
		doc.Metadata.Title = strings.TrimSuffix(filepath.Base(files[0]), filepath.Ext(files[0]))
	}

	// Process cover image if specified
	if doc.Metadata.CoverImage != "" {
		if err := c.processCoverImage(doc, result); err != nil {
			result.AddWarning(fmt.Sprintf("Cover image: %s", err))
		}
	}

	// Process images
	c.processImages(doc, result)

	// Build EPUB
	epubData, err := c.builder.Build(doc)
	if err != nil {
		return result, fmt.Errorf("building EPUB: %w", err)
	}

	// Write output file
	outputPath := opts.OutputPath
	if outputPath == "" {
		outputPath = strings.TrimSuffix(filepath.Base(files[0]), filepath.Ext(files[0])) + ".epub"
	}

	if err := c.writeOutput(outputPath, epubData); err != nil {
		return result, err
	}

	// Build result
	result.Success = true
	result.OutputPath = outputPath
	result.Stats = model.ConversionStats{
		InputFormat:  format.String(),
		InputFiles:   len(files),
		ChapterCount: len(doc.Chapters),
		ImageCount:   len(doc.Resources),
		OutputSize:   int64(len(epubData)),
		Duration:     time.Since(start),
	}

	return result, nil
}

// ConvertContent converts raw content bytes to EPUB.
func (c *Converter) ConvertContent(content []byte, opts Options) (*model.ConversionResult, error) {
	start := time.Now()
	result := &model.ConversionResult{
		Success:  false,
		Warnings: make([]string, 0),
	}

	// Detect format
	format := c.detectFormatFromString(opts.InputFormat)
	if format == parser.FormatUnknown {
		format = parser.FormatMarkdown // Default to markdown
	}

	// Get parser
	p := c.getParser(format)
	if p == nil {
		return result, fmt.Errorf("%w: no parser for format %s", ErrUnsupportedFmt, format)
	}

	// Parse content
	doc, err := p.Parse(content, ".")
	if err != nil {
		return result, fmt.Errorf("parsing content: %w", err)
	}

	// Apply CLI metadata overrides
	if opts.CLIMetadata != nil {
		doc.Metadata.Merge(opts.CLIMetadata)
	}

	// Ensure document has a title
	if doc.Metadata.Title == "" {
		doc.Metadata.Title = "Untitled Document"
	}

	// Build EPUB
	epubData, err := c.builder.Build(doc)
	if err != nil {
		return result, fmt.Errorf("building EPUB: %w", err)
	}

	// Write output
	outputPath := opts.OutputPath
	if outputPath == "" {
		outputPath = "output.epub"
	}

	if err := c.writeOutput(outputPath, epubData); err != nil {
		return result, err
	}

	// Build result
	result.Success = true
	result.OutputPath = outputPath
	result.Stats = model.ConversionStats{
		InputFormat:  format.String(),
		InputFiles:   1,
		ChapterCount: len(doc.Chapters),
		ImageCount:   len(doc.Resources),
		OutputSize:   int64(len(epubData)),
		Duration:     time.Since(start),
	}

	return result, nil
}

// expandInputs expands directories and validates file existence.
func (c *Converter) expandInputs(inputs []string) ([]string, error) {
	var files []string

	for _, input := range inputs {
		info, err := os.Stat(input)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrFileNotFound, input)
		}

		if info.IsDir() {
			// Expand directory (non-recursive)
			dirFiles, err := c.expandDirectory(input)
			if err != nil {
				return nil, err
			}
			files = append(files, dirFiles...)
		} else {
			files = append(files, input)
		}
	}

	// Sort files alphabetically for consistent ordering
	sort.Strings(files)
	return files, nil
}

// expandDirectory lists supported files in a directory.
func (c *Converter) expandDirectory(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if c.isSupportedExtension(ext) {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	return files, nil
}

// isSupportedExtension checks if file extension is supported.
func (c *Converter) isSupportedExtension(ext string) bool {
	supported := []string{".md", ".markdown", ".html", ".htm", ".pdf"}
	for _, s := range supported {
		if ext == s {
			return true
		}
	}
	return false
}

// detectFormat determines the input format from file extension or explicit format.
func (c *Converter) detectFormat(file string, explicit string) parser.Format {
	if explicit != "" {
		return c.detectFormatFromString(explicit)
	}

	ext := strings.ToLower(filepath.Ext(file))
	switch ext {
	case ".md", ".markdown":
		return parser.FormatMarkdown
	case ".html", ".htm":
		return parser.FormatHTML
	case ".pdf":
		return parser.FormatPDF
	default:
		return parser.FormatUnknown
	}
}

// detectFormatFromString converts format string to Format type.
func (c *Converter) detectFormatFromString(s string) parser.Format {
	switch strings.ToLower(s) {
	case "md", "markdown":
		return parser.FormatMarkdown
	case "html", "htm":
		return parser.FormatHTML
	case "pdf":
		return parser.FormatPDF
	default:
		return parser.FormatUnknown
	}
}

// getParser returns the parser for the given format.
func (c *Converter) getParser(format parser.Format) parser.Parser {
	return c.parsers[format]
}

// mergeDocument merges a parsed document into the main document.
func (c *Converter) mergeDocument(main, parsed *model.Document, index int) {
	// Merge metadata (first file wins, except explicit overrides)
	if index == 0 {
		main.Metadata = parsed.Metadata
	}

	// Update chapter ordering for merged chapters
	offset := len(main.Chapters)
	for i, chapter := range parsed.Chapters {
		chapter.Order = offset + i
		chapter.ID = fmt.Sprintf("chapter-%03d", chapter.Order+1)
		chapter.FileName = fmt.Sprintf("content/chapter-%03d.xhtml", chapter.Order+1)
		main.AddChapter(chapter)
	}

	// Merge TOC entries
	main.TOC.Entries = append(main.TOC.Entries, parsed.TOC.Entries...)

	// Merge resources
	for _, res := range parsed.Resources {
		main.AddResource(res)
	}
}

// processCoverImage loads and embeds the cover image.
func (c *Converter) processCoverImage(doc *model.Document, result *model.ConversionResult) error {
	coverPath := doc.Metadata.CoverImage

	resource, err := c.imgHandler.ProcessImage(coverPath, ".")
	if err != nil {
		return err
	}

	// Mark as cover image
	resource.IsCover = true
	resource.ID = "cover-image"
	resource.FileName = "images/cover" + extensionFromMediaType(resource.MediaType)

	doc.AddResource(*resource)
	return nil
}

// extensionFromMediaType returns file extension for a MIME type.
func extensionFromMediaType(mediaType string) string {
	switch mediaType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "image/svg+xml":
		return ".svg"
	default:
		return ".bin"
	}
}

// processImages handles image resources in the document.
func (c *Converter) processImages(doc *model.Document, result *model.ConversionResult) {
	// Image processing will be handled by the image handler
	// For now, just count existing resources
	for range doc.Resources {
		// Resources are already processed by parser
	}
}

// writeOutput writes EPUB data to the output file.
func (c *Converter) writeOutput(path string, data []byte) error {
	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("%w: cannot create directory %s", ErrOutputNotWrite, dir)
		}
	}

	// Write to temp file first, then rename (atomic operation)
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("%w: %s", ErrOutputNotWrite, err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("%w: %s", ErrOutputNotWrite, err)
	}

	return nil
}
