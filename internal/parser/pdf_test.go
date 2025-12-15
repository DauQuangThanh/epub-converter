package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPDFParser_Parse_RealPDF(t *testing.T) {
	// Skip if test PDF doesn't exist
	pdfPath := filepath.Join("..", "..", "tests", "fixtures", "pdf", "sample.pdf")
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		t.Skip("Test PDF not available")
	}

	content, err := os.ReadFile(pdfPath)
	require.NoError(t, err)

	p := NewPDFParser()
	doc, err := p.Parse(content, ".")

	require.NoError(t, err)
	assert.NotNil(t, doc)

	// Should have extracted some content
	assert.NotEmpty(t, doc.Chapters)
	assert.NotEmpty(t, doc.Chapters[0].Content)

	// Should have a title
	assert.NotEmpty(t, doc.Metadata.Title)
}

func TestPDFParser_Parse_InvalidPDF(t *testing.T) {
	p := NewPDFParser()

	// Test with non-PDF content
	_, err := p.Parse([]byte("This is not a PDF"), ".")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PDF")
}

func TestPDFParser_Parse_EmptyContent(t *testing.T) {
	p := NewPDFParser()

	_, err := p.Parse([]byte{}, ".")
	assert.Error(t, err)
}

func TestPDFParser_SupportedExtensions(t *testing.T) {
	p := NewPDFParser()
	exts := p.SupportedExtensions()

	assert.Contains(t, exts, ".pdf")
}

func TestPDFParser_looksLikeHeading(t *testing.T) {
	p := NewPDFParser()

	tests := []struct {
		text     string
		expected bool
	}{
		{"Chapter 1", true},
		{"Introduction", true},
		// This string is over 200 chars to trigger the length check
		{"This is a very long sentence that goes on and on and on and contains lots and lots of text and probably should not ever be considered a heading because it is definitely way too long for a typical heading in any reasonable document.", false},
		{"Ends with comma,", false},
		{"123", false}, // no letters
		{"Valid Heading", true},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := p.looksLikeHeading(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPDFParser_fontSizeToHeadingLevel(t *testing.T) {
	p := NewPDFParser()

	tests := []struct {
		fontSize float64
		expected int
	}{
		{24, 1},
		{30, 1},
		{18, 2},
		{20, 2},
		{14, 3},
		{16, 3},
		{12, 4},
		{10, 4},
	}

	for _, tt := range tests {
		result := p.fontSizeToHeadingLevel(tt.fontSize)
		assert.Equal(t, tt.expected, result, "fontSize %.1f", tt.fontSize)
	}
}

func TestPDFParser_extractTitle(t *testing.T) {
	p := NewPDFParser()

	tests := []struct {
		name     string
		text     string
		headings []headingInfo
		expected string
	}{
		{
			name:     "uses H1 heading",
			text:     "Some text",
			headings: []headingInfo{{Level: 1, Title: "Main Title"}},
			expected: "Main Title",
		},
		{
			name:     "uses first heading if no H1",
			text:     "Some text",
			headings: []headingInfo{{Level: 2, Title: "Section Title"}},
			expected: "Section Title",
		},
		{
			name:     "uses first line if no headings",
			text:     "First Line\nSecond Line",
			headings: []headingInfo{},
			expected: "First Line",
		},
		{
			name:     "falls back to default",
			text:     "",
			headings: []headingInfo{},
			expected: "Untitled Document",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.extractTitle(tt.text, tt.headings)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPDFParser_textToXHTML(t *testing.T) {
	p := NewPDFParser()

	text := `###HEADING_1### Main Title
Some paragraph text.

###HEADING_2### Section
More text here.`

	headings := []headingInfo{
		{Level: 1, Title: "Main Title", ID: "main-title"},
		{Level: 2, Title: "Section", ID: "section"},
	}

	result := p.textToXHTML(text, headings)

	// Should contain heading tags
	assert.Contains(t, result, "<h1")
	assert.Contains(t, result, "<h2")
	assert.Contains(t, result, "Main Title")
	assert.Contains(t, result, "Section")

	// Should contain paragraphs
	assert.Contains(t, result, "<p>")
	assert.Contains(t, result, "paragraph text")
}

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "Hello"},
		{"<tag>", "&lt;tag&gt;"},
		{"A & B", "A &amp; B"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"It's", "It&#39;s"},
	}

	for _, tt := range tests {
		result := escapeXML(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}
