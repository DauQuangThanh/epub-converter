package parser

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"github.com/dauquangthanh/epub-converter/internal/model"
)

// HTMLParser parses HTML content to Document model.
type HTMLParser struct{}

// NewHTMLParser creates a new HTML parser.
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

// Parse converts HTML content to a Document.
func (p *HTMLParser) Parse(content []byte, basePath string) (*model.Document, error) {
	doc := model.NewDocument()

	// Parse HTML
	htmlDoc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	// Extract metadata from head
	p.extractMetadata(htmlDoc, doc)

	// Extract headings for TOC
	headings := p.extractHeadings(htmlDoc)

	// Convert body to XHTML
	bodyContent := p.extractBody(htmlDoc)
	if bodyContent == "" {
		// If no body, use entire content
		bodyContent = string(content)
	}

	// Clean and convert to XHTML
	xhtmlContent := p.convertToXHTML(bodyContent)

	// Extract image references
	images := p.extractImageRefs(xhtmlContent, basePath)
	for _, img := range images {
		doc.AddResource(img)
	}

	// Rewrite image paths for EPUB
	xhtmlContent = p.rewriteImagePaths(xhtmlContent)

	// Strip JavaScript
	xhtmlContent = p.stripJavaScript(xhtmlContent)

	// Extract CSS
	css := p.extractCSS(htmlDoc, basePath)
	if css != "" {
		cssResource := model.Resource{
			ID:        "inline-css",
			FileName:  "styles/inline.css",
			MediaType: "text/css",
			Data:      []byte(css),
		}
		doc.AddResource(cssResource)
	}

	// Create chapter
	title := doc.Metadata.Title
	if title == "" && len(headings) > 0 {
		title = headings[0].Title
		doc.Metadata.Title = title
	}

	chapter := model.Chapter{
		ID:       "chapter-001",
		Title:    title,
		Level:    1,
		Content:  xhtmlContent,
		FileName: "content/chapter-001.xhtml",
		Order:    0,
	}
	doc.AddChapter(chapter)

	// Build TOC
	doc.TOC = *p.buildTOC(headings)

	return doc, nil
}

// SupportedExtensions returns file extensions this parser handles.
func (p *HTMLParser) SupportedExtensions() []string {
	return []string{".html", ".htm"}
}

// extractMetadata extracts metadata from HTML head.
func (p *HTMLParser) extractMetadata(doc *html.Node, mdoc *model.Document) {
	var findMeta func(*html.Node)
	findMeta = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil {
					mdoc.Metadata.Title = n.FirstChild.Data
				}
			case "meta":
				name := ""
				content := ""
				for _, attr := range n.Attr {
					switch attr.Key {
					case "name":
						name = attr.Val
					case "content":
						content = attr.Val
					}
				}
				switch strings.ToLower(name) {
				case "author":
					if content != "" {
						mdoc.Metadata.Authors = append(mdoc.Metadata.Authors, content)
					}
				case "description":
					mdoc.Metadata.Description = content
				case "language":
					mdoc.Metadata.Language = content
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findMeta(c)
		}
	}
	findMeta(doc)
}

// extractHeadings extracts heading elements for TOC.
func (p *HTMLParser) extractHeadings(doc *html.Node) []headingInfo {
	var headings []headingInfo

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			level := 0
			switch n.Data {
			case "h1":
				level = 1
			case "h2":
				level = 2
			case "h3":
				level = 3
			case "h4":
				level = 4
			case "h5":
				level = 5
			case "h6":
				level = 6
			}

			if level > 0 {
				text := p.extractText(n)
				id := p.getAttr(n, "id")
				if id == "" {
					id = generateHeadingID(text)
				}
				headings = append(headings, headingInfo{
					Level: level,
					Title: text,
					ID:    id,
				})
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	return headings
}

// extractText extracts text content from a node.
func (p *HTMLParser) extractText(n *html.Node) string {
	var text strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			text.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.TrimSpace(text.String())
}

// getAttr gets an attribute value from a node.
func (p *HTMLParser) getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// extractBody extracts body content as string.
func (p *HTMLParser) extractBody(doc *html.Node) string {
	var body *html.Node

	var find func(*html.Node)
	find = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(doc)

	if body == nil {
		return ""
	}

	var buf bytes.Buffer
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		html.Render(&buf, c)
	}
	return buf.String()
}

// convertToXHTML converts HTML to valid XHTML.
func (p *HTMLParser) convertToXHTML(content string) string {
	// Self-close void elements
	voidElements := []string{"br", "hr", "img", "input", "meta", "link", "area", "base", "col", "embed", "param", "source", "track", "wbr"}
	for _, elem := range voidElements {
		// Match <elem ...> not already self-closed and convert to <elem ... />
		// First, normalize any existing self-closed tags
		reAlreadyClosed := regexp.MustCompile(`<(` + elem + `)([^>]*)\s*/>`)
		content = reAlreadyClosed.ReplaceAllString(content, `<$1$2 />`)

		// Then, close unclosed void elements
		reUnclosed := regexp.MustCompile(`<(` + elem + `)([^/>]*[^/])>`)
		content = reUnclosed.ReplaceAllString(content, `<$1$2 />`)

		// Handle simple tags like <br> or <hr>
		reSimple := regexp.MustCompile(`<(` + elem + `)>`)
		content = reSimple.ReplaceAllString(content, `<$1 />`)
	}

	// Ensure lowercase tags (HTML5 is case-insensitive, XHTML requires lowercase)
	content = strings.ReplaceAll(content, "<BR>", "<br />")
	content = strings.ReplaceAll(content, "<HR>", "<hr />")
	content = strings.ReplaceAll(content, "<IMG", "<img")
	content = strings.ReplaceAll(content, "</IMG>", "")

	return content
}

// stripJavaScript removes script elements.
func (p *HTMLParser) stripJavaScript(content string) string {
	// Remove script tags and their content
	scriptRe := regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	content = scriptRe.ReplaceAllString(content, "")

	// Remove on* event handlers
	eventRe := regexp.MustCompile(`\s+on\w+="[^"]*"`)
	content = eventRe.ReplaceAllString(content, "")

	return content
}

// extractCSS extracts inline and style tag CSS.
func (p *HTMLParser) extractCSS(doc *html.Node, basePath string) string {
	var css strings.Builder

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "style" {
			if n.FirstChild != nil {
				css.WriteString(n.FirstChild.Data)
				css.WriteString("\n")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	return css.String()
}

// extractImageRefs finds image references in content.
func (p *HTMLParser) extractImageRefs(content string, basePath string) []model.Resource {
	var resources []model.Resource

	imgRe := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["']`)
	matches := imgRe.FindAllStringSubmatch(content, -1)

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

		if seen[src] {
			continue
		}
		seen[src] = true

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
		default:
			continue
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
		}
		resources = append(resources, resource)
	}

	return resources
}

// rewriteImagePaths updates image paths to EPUB-relative paths.
func (p *HTMLParser) rewriteImagePaths(content string) string {
	imgRe := regexp.MustCompile(`(<img[^>]+src=["'])([^"']+)(["'])`)
	return imgRe.ReplaceAllStringFunc(content, func(match string) string {
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

		baseName := filepath.Base(src)
		newSrc := "../images/" + baseName

		return parts[1] + newSrc + parts[3]
	})
}

// buildTOC creates table of contents from headings.
func (p *HTMLParser) buildTOC(headings []headingInfo) *model.TableOfContents {
	var entries []model.TOCEntry

	chapterFile := "content/chapter-001.xhtml"

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
