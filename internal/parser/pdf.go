package parser

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/ledongthuc/pdf"

	"github.com/dauquangthanh/epub-converter/internal/model"
)

// PDFParser parses PDF content to Document model.
type PDFParser struct {
	minHeadingFontSize float64
}

// NewPDFParser creates a new PDF parser.
func NewPDFParser() *PDFParser {
	return &PDFParser{
		minHeadingFontSize: 14.0, // Consider text with font size >= 14 as potential heading
	}
}

// Parse converts PDF content to a Document.
func (p *PDFParser) Parse(content []byte, basePath string) (*model.Document, error) {
	doc := model.NewDocument()

	// Create a temporary file to read PDF
	tmpFile, err := os.CreateTemp("", "toepub-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		return nil, fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("closing temp file: %w", err)
	}

	// Open and read PDF
	pdfFile, pdfReader, err := pdf.Open(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("opening PDF: %w", err)
	}
	defer pdfFile.Close()

	numPages := pdfReader.NumPage()
	if numPages == 0 {
		return nil, fmt.Errorf("PDF has no pages")
	}

	// Extract text and structure from all pages
	var allText strings.Builder
	var headings []headingInfo

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		// Extract text content
		pageText, pageHeadings := p.extractPageContent(page, pageNum)
		allText.WriteString(pageText)
		headings = append(headings, pageHeadings...)

		// Add page break marker for multi-page documents
		if pageNum < numPages {
			allText.WriteString("\n\n")
		}
	}

	text := strings.TrimSpace(allText.String())
	if text == "" {
		return nil, fmt.Errorf("PDF contains no extractable text (might be image-based)")
	}

	// Try to extract title from first heading or first line
	title := p.extractTitle(text, headings)
	doc.Metadata.Title = title

	// Convert text to XHTML content
	xhtmlContent := p.textToXHTML(text, headings)

	// Create chapter
	chapter := model.Chapter{
		ID:       "chapter-001",
		Title:    title,
		Level:    1,
		Content:  xhtmlContent,
		FileName: "content/chapter-001.xhtml",
		Order:    0,
	}
	doc.AddChapter(chapter)

	// Build TOC from headings
	doc.TOC = *p.buildTOC(headings)

	return doc, nil
}

// SupportedExtensions returns file extensions this parser handles.
func (p *PDFParser) SupportedExtensions() []string {
	return []string{".pdf"}
}

// extractPageContent extracts text and headings from a PDF page.
func (p *PDFParser) extractPageContent(page pdf.Page, pageNum int) (string, []headingInfo) {
	var text strings.Builder
	var headings []headingInfo

	rows, err := page.GetTextByRow()
	if err != nil {
		// Fall back to plain text extraction
		plainText, err := page.GetPlainText(nil)
		if err == nil {
			text.WriteString(plainText)
		}
		return text.String(), headings
	}

	// Sort rows by Y position (top to bottom)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Position > rows[j].Position
	})

	for _, row := range rows {
		var lineText strings.Builder
		var maxFontSize float64

		for _, word := range row.Content {
			lineText.WriteString(word.S)
			lineText.WriteString(" ")
			if word.FontSize > maxFontSize {
				maxFontSize = word.FontSize
			}
		}

		line := strings.TrimSpace(lineText.String())
		if line == "" {
			continue
		}

		// Detect potential headings based on font size
		if maxFontSize >= p.minHeadingFontSize && p.looksLikeHeading(line) {
			level := p.fontSizeToHeadingLevel(maxFontSize)
			id := generateHeadingID(line)
			headings = append(headings, headingInfo{
				Level: level,
				Title: line,
				ID:    id,
			})
			// Mark as heading in text
			text.WriteString(fmt.Sprintf("\n###HEADING_%d### %s\n", level, line))
		} else {
			text.WriteString(line)
			text.WriteString("\n")
		}
	}

	return text.String(), headings
}

// looksLikeHeading checks if text looks like a heading (not too long, not punctuation-heavy).
func (p *PDFParser) looksLikeHeading(text string) bool {
	// Skip if too long
	if len(text) > 200 {
		return false
	}

	// Skip if it ends with typical sentence punctuation
	trimmed := strings.TrimSpace(text)
	if strings.HasSuffix(trimmed, ",") {
		return false
	}

	// Must have at least one letter
	hasLetter := false
	for _, r := range text {
		if unicode.IsLetter(r) {
			hasLetter = true
			break
		}
	}

	return hasLetter
}

// fontSizeToHeadingLevel converts font size to heading level.
func (p *PDFParser) fontSizeToHeadingLevel(fontSize float64) int {
	switch {
	case fontSize >= 24:
		return 1
	case fontSize >= 18:
		return 2
	case fontSize >= 14:
		return 3
	default:
		return 4
	}
}

// extractTitle extracts a title from the text or headings.
func (p *PDFParser) extractTitle(text string, headings []headingInfo) string {
	// Use first H1 heading if available
	for _, h := range headings {
		if h.Level == 1 {
			return h.Title
		}
	}

	// Use first heading
	if len(headings) > 0 {
		return headings[0].Title
	}

	// Use first non-empty line
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && len(line) < 100 {
			return line
		}
	}

	return "Untitled Document"
}

// textToXHTML converts extracted PDF text to XHTML content.
func (p *PDFParser) textToXHTML(text string, headings []headingInfo) string {
	var xhtml strings.Builder

	// Process text line by line
	lines := strings.Split(text, "\n")
	var currentParagraph strings.Builder
	inParagraph := false

	headingRe := regexp.MustCompile(`^###HEADING_(\d+)###\s*(.+)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for heading marker
		if match := headingRe.FindStringSubmatch(line); match != nil {
			// Close current paragraph if open
			if inParagraph {
				xhtml.WriteString("<p>")
				xhtml.WriteString(escapeXML(strings.TrimSpace(currentParagraph.String())))
				xhtml.WriteString("</p>\n")
				currentParagraph.Reset()
				inParagraph = false
			}

			level := match[1]
			title := match[2]
			id := generateHeadingID(title)
			xhtml.WriteString(fmt.Sprintf("<h%s id=\"%s\">%s</h%s>\n", level, id, escapeXML(title), level))
			continue
		}

		// Empty line marks paragraph break
		if line == "" {
			if inParagraph {
				xhtml.WriteString("<p>")
				xhtml.WriteString(escapeXML(strings.TrimSpace(currentParagraph.String())))
				xhtml.WriteString("</p>\n")
				currentParagraph.Reset()
				inParagraph = false
			}
			continue
		}

		// Accumulate text for paragraph
		if inParagraph {
			currentParagraph.WriteString(" ")
		}
		currentParagraph.WriteString(line)
		inParagraph = true
	}

	// Close final paragraph
	if inParagraph {
		xhtml.WriteString("<p>")
		xhtml.WriteString(escapeXML(strings.TrimSpace(currentParagraph.String())))
		xhtml.WriteString("</p>\n")
	}

	return xhtml.String()
}

// escapeXML escapes special XML characters.
func escapeXML(s string) string {
	var buf bytes.Buffer
	for _, r := range s {
		switch r {
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '&':
			buf.WriteString("&amp;")
		case '"':
			buf.WriteString("&quot;")
		case '\'':
			buf.WriteString("&#39;")
		default:
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

// buildTOC creates table of contents from headings.
func (p *PDFParser) buildTOC(headings []headingInfo) *model.TableOfContents {
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

// extractImagesFromPDF extracts images from PDF using pdfcpu.
// Note: Image extraction is a separate optional step.
func (p *PDFParser) extractImagesFromPDF(pdfPath, outputDir string) ([]model.Resource, error) {
	// This would use pdfcpu for image extraction
	// For now, return empty as image extraction is optional per spec
	return nil, nil
}
