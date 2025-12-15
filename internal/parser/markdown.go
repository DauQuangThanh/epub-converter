// ------------------------------------------------------------------
// Developed by Dau Quang Thanh - 2025.
// Enterprise AI Solution Architect
//
// Happy Reading!
// ------------------------------------------------------------------

package parser

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"

	"github.com/dauquangthanh/epub-converter/internal/model"
)

// MarkdownParser parses Markdown content using goldmark with GFM support.
type MarkdownParser struct {
	md goldmark.Markdown
}

// NewMarkdownParser creates a new Markdown parser with GFM extensions.
func NewMarkdownParser() *MarkdownParser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // Tables, task lists, strikethrough, autolinks
			&frontmatter.Extender{}, // YAML/TOML front matter
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // Generate heading IDs
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),         // Generate XHTML for EPUB
			html.WithUnsafe(),        // Allow raw HTML in markdown
		),
	)

	return &MarkdownParser{md: md}
}

// Parse converts Markdown content to a Document.
func (p *MarkdownParser) Parse(content []byte, basePath string) (*model.Document, error) {
	doc := model.NewDocument()

	// Parse front matter and content
	var meta map[string]interface{}
	body := content

	// Try to extract front matter
	if bytes.HasPrefix(content, []byte("---")) {
		meta, body = p.extractFrontMatter(content)
	}

	// Apply front matter metadata
	p.applyMetadata(doc, meta)

	// Parse markdown to AST
	reader := text.NewReader(body)
	astDoc := p.md.Parser().Parse(reader)

	// Extract headings for TOC
	headings := p.extractHeadings(astDoc, body)

	// Render to XHTML
	var buf bytes.Buffer
	if err := p.md.Renderer().Render(&buf, body, astDoc); err != nil {
		return nil, fmt.Errorf("rendering markdown: %w", err)
	}

	htmlContent := buf.String()

	// Process image references
	images := p.extractImageRefs(htmlContent, basePath)
	for _, img := range images {
		doc.AddResource(img)
	}

	// Update image paths in content
	htmlContent = p.rewriteImagePaths(htmlContent)

	// Create chapters from headings or single chapter
	p.createChapters(doc, htmlContent, headings)

	// Build TOC
	doc.TOC = *p.buildTOC(headings, doc.Chapters)

	return doc, nil
}

// SupportedExtensions returns file extensions this parser handles.
func (p *MarkdownParser) SupportedExtensions() []string {
	return []string{".md", ".markdown"}
}

// extractFrontMatter parses YAML front matter from content.
func (p *MarkdownParser) extractFrontMatter(content []byte) (map[string]interface{}, []byte) {
	// Find front matter boundaries
	lines := bytes.Split(content, []byte("\n"))
	if len(lines) < 2 || string(bytes.TrimSpace(lines[0])) != "---" {
		return nil, content
	}

	endIdx := -1
	for i := 1; i < len(lines); i++ {
		if string(bytes.TrimSpace(lines[i])) == "---" {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		return nil, content
	}

	// Parse YAML front matter using goldmark-frontmatter
	// Create a new parser context to extract front matter
	ctx := parser.NewContext()
	reader := text.NewReader(content)
	p.md.Parser().Parse(reader, parser.WithContext(ctx))

	// Get front matter data
	fm := frontmatter.Get(ctx)
	if fm == nil {
		return nil, content
	}

	var meta map[string]interface{}
	if err := fm.Decode(&meta); err != nil {
		return nil, content
	}

	// Return body after front matter
	bodyStart := 0
	for i := 0; i <= endIdx; i++ {
		bodyStart += len(lines[i]) + 1
	}
	if bodyStart > len(content) {
		bodyStart = len(content)
	}

	return meta, content[bodyStart:]
}

// applyMetadata applies front matter values to document metadata.
func (p *MarkdownParser) applyMetadata(doc *model.Document, meta map[string]interface{}) {
	if meta == nil {
		return
	}

	if title, ok := meta["title"].(string); ok {
		doc.Metadata.Title = title
	}

	// Handle author as string or list
	switch author := meta["author"].(type) {
	case string:
		doc.Metadata.Authors = []string{author}
	case []interface{}:
		for _, a := range author {
			if s, ok := a.(string); ok {
				doc.Metadata.Authors = append(doc.Metadata.Authors, s)
			}
		}
	}

	if lang, ok := meta["language"].(string); ok {
		doc.Metadata.Language = lang
	}
	if lang, ok := meta["lang"].(string); ok {
		doc.Metadata.Language = lang
	}

	if desc, ok := meta["description"].(string); ok {
		doc.Metadata.Description = desc
	}

	if publisher, ok := meta["publisher"].(string); ok {
		doc.Metadata.Publisher = publisher
	}
}

// extractHeadings walks the AST to find all headings.
func (p *MarkdownParser) extractHeadings(doc ast.Node, source []byte) []headingInfo {
	var headings []headingInfo

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if h, ok := n.(*ast.Heading); ok {
			text := string(h.Text(source))
			id := generateHeadingID(text)

			headings = append(headings, headingInfo{
				Level: h.Level,
				Title: text,
				ID:    id,
			})
		}

		return ast.WalkContinue, nil
	})

	return headings
}

// headingInfo stores heading information for TOC building.
type headingInfo struct {
	Level int
	Title string
	ID    string
}

// generateHeadingID creates a URL-safe ID from heading text.
func generateHeadingID(text string) string {
	// Convert to lowercase
	id := strings.ToLower(text)

	// Replace spaces with hyphens
	id = strings.ReplaceAll(id, " ", "-")

	// Remove non-alphanumeric characters except hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	id = reg.ReplaceAllString(id, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	id = reg.ReplaceAllString(id, "-")

	// Trim leading/trailing hyphens
	id = strings.Trim(id, "-")

	if id == "" {
		id = "heading"
	}

	return id
}

// extractImageRefs finds all image references in the HTML content.
func (p *MarkdownParser) extractImageRefs(html string, basePath string) []model.Resource {
	var resources []model.Resource

	// Match img src attributes
	imgRe := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["']`)
	matches := imgRe.FindAllStringSubmatch(html, -1)

	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		src := match[1]

		// Skip remote URLs and data URIs
		if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") ||
			strings.HasPrefix(src, "data:") {
			continue
		}

		// Skip duplicates
		if seen[src] {
			continue
		}
		seen[src] = true

		// Create resource placeholder (actual loading done by converter)
		baseName := filepath.Base(src)
		ext := strings.ToLower(filepath.Ext(baseName))

		mediaType := ""
		switch ext {
		case ".png":
			mediaType = "image/png"
		case ".jpg", ".jpeg":
			mediaType = "image/jpeg"
		case ".gif":
			mediaType = "image/gif"
		case ".svg":
			mediaType = "image/svg+xml"
		case ".webp":
			mediaType = "image/png" // Will be converted
		default:
			continue // Skip unsupported formats
		}

		// Resolve source path relative to basePath
		sourcePath := src
		if !filepath.IsAbs(src) {
			sourcePath = filepath.Join(basePath, src)
		}

		resource := model.Resource{
			ID:         "img-" + sanitizeID(strings.TrimSuffix(baseName, ext)),
			FileName:   "images/" + baseName,
			MediaType:  mediaType,
			SourcePath: sourcePath, // Store resolved absolute path
			// Data will be loaded by converter
		}

		resources = append(resources, resource)
	}

	return resources
}

// rewriteImagePaths updates image paths to EPUB-relative paths.
func (p *MarkdownParser) rewriteImagePaths(html string) string {
	imgRe := regexp.MustCompile(`(<img[^>]+src=["'])([^"']+)(["'])`)
	return imgRe.ReplaceAllStringFunc(html, func(match string) string {
		parts := imgRe.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}

		src := parts[2]

		// Skip remote URLs and data URIs
		if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") ||
			strings.HasPrefix(src, "data:") {
			return match
		}

		// Rewrite to EPUB path
		baseName := filepath.Base(src)
		newSrc := "../images/" + baseName

		return parts[1] + newSrc + parts[3]
	})
}

// createChapters creates chapters from content and headings.
func (p *MarkdownParser) createChapters(doc *model.Document, content string, headings []headingInfo) {
	if len(headings) == 0 {
		// Single chapter for entire content
		chapter := model.Chapter{
			ID:       "chapter-001",
			Title:    doc.Metadata.Title,
			Level:    1,
			Content:  content,
			FileName: "content/chapter-001.xhtml",
			Order:    0,
		}
		doc.AddChapter(chapter)
		return
	}

	// For now, create a single chapter with all content
	// TODO: Split content at h1/h2 boundaries for multi-chapter support
	title := headings[0].Title
	if doc.Metadata.Title == "" {
		doc.Metadata.Title = title
	}

	chapter := model.Chapter{
		ID:       "chapter-001",
		Title:    title,
		Level:    headings[0].Level,
		Content:  content,
		FileName: "content/chapter-001.xhtml",
		Order:    0,
	}
	doc.AddChapter(chapter)
}

// buildTOC creates table of contents from headings.
func (p *MarkdownParser) buildTOC(headings []headingInfo, chapters []model.Chapter) *model.TableOfContents {
	var entries []model.TOCEntry

	if len(chapters) == 0 {
		return model.NewTableOfContents()
	}

	// Map headings to TOC entries
	chapterFile := chapters[0].FileName

	for _, h := range headings {
		entry := model.TOCEntry{
			Title: h.Title,
			Href:  chapterFile + "#" + h.ID,
			Level: h.Level,
		}
		entries = append(entries, entry)
	}

	return model.BuildFromHeadings(entries)
}

// sanitizeID converts a string to a valid XML ID.
func sanitizeID(s string) string {
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else {
			result.WriteRune('-')
		}
	}
	return strings.Trim(result.String(), "-")
}
