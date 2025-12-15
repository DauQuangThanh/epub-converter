package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTMLParser_Parse_SimpleDocument(t *testing.T) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Test Document</title>
    <meta name="author" content="Test Author">
    <meta name="description" content="Test description">
    <meta name="language" content="en">
</head>
<body>
    <h1>Hello World</h1>
    <p>This is a test paragraph.</p>
    <h2>Section One</h2>
    <p>More content here.</p>
</body>
</html>`

	p := NewHTMLParser()
	doc, err := p.Parse([]byte(html), ".")

	require.NoError(t, err)
	assert.NotNil(t, doc)

	// Check metadata
	assert.Equal(t, "Test Document", doc.Metadata.Title)
	assert.Contains(t, doc.Metadata.Authors, "Test Author")
	assert.Equal(t, "Test description", doc.Metadata.Description)
	assert.Equal(t, "en", doc.Metadata.Language)

	// Check chapters
	assert.Len(t, doc.Chapters, 1)
	assert.Equal(t, "chapter-001", doc.Chapters[0].ID)
	assert.Contains(t, doc.Chapters[0].Content, "Hello World")
	assert.Contains(t, doc.Chapters[0].Content, "This is a test paragraph")

	// Check TOC
	assert.NotEmpty(t, doc.TOC.Entries)
}

func TestHTMLParser_Parse_ExtractsHeadings(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<body>
    <h1>Main Title</h1>
    <h2>Section 1</h2>
    <h3>Subsection 1.1</h3>
    <h2>Section 2</h2>
</body>
</html>`

	p := NewHTMLParser()
	doc, err := p.Parse([]byte(html), ".")

	require.NoError(t, err)
	assert.NotEmpty(t, doc.TOC.Entries)

	// Verify TOC structure
	entries := doc.TOC.Entries
	assert.Equal(t, "Main Title", entries[0].Title)
	assert.Equal(t, 1, entries[0].Level)
}

func TestHTMLParser_Parse_ConvertsToXHTML(t *testing.T) {
	html := `<html>
<body>
    <p>Text with<br>line break</p>
    <hr>
    <img src="test.png" alt="test">
</body>
</html>`

	p := NewHTMLParser()
	doc, err := p.Parse([]byte(html), ".")

	require.NoError(t, err)

	content := doc.Chapters[0].Content
	// Verify void elements are self-closed
	assert.Contains(t, content, "<br />")
	assert.Contains(t, content, "<hr />")
	assert.Contains(t, content, "/>")
}

func TestHTMLParser_Parse_StripsJavaScript(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <script>alert('test');</script>
</head>
<body>
    <h1>Test</h1>
    <script type="text/javascript">
        console.log("should be removed");
    </script>
    <button onclick="doSomething()">Click</button>
    <a href="#" onmouseover="highlight()">Link</a>
</body>
</html>`

	p := NewHTMLParser()
	doc, err := p.Parse([]byte(html), ".")

	require.NoError(t, err)

	content := doc.Chapters[0].Content

	// Verify scripts are removed
	assert.NotContains(t, content, "<script")
	assert.NotContains(t, content, "alert")
	assert.NotContains(t, content, "console.log")

	// Verify event handlers are removed
	assert.NotContains(t, content, "onclick")
	assert.NotContains(t, content, "onmouseover")
}

func TestHTMLParser_Parse_ExtractsCSS(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: serif; }
        h1 { color: blue; }
    </style>
</head>
<body>
    <h1>Test</h1>
</body>
</html>`

	p := NewHTMLParser()
	doc, err := p.Parse([]byte(html), ".")

	require.NoError(t, err)

	// Should have inline CSS resource
	var cssResource *struct {
		found bool
		data  string
	}
	for _, r := range doc.Resources {
		if r.MediaType == "text/css" {
			cssResource = &struct {
				found bool
				data  string
			}{true, string(r.Data)}
			break
		}
	}

	assert.NotNil(t, cssResource)
	assert.True(t, cssResource.found)
	assert.Contains(t, cssResource.data, "font-family")
	assert.Contains(t, cssResource.data, "color: blue")
}

func TestHTMLParser_Parse_ExtractsLocalImageRefs(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<body>
    <h1>Images</h1>
    <img src="local.png" alt="Local">
    <img src="images/photo.jpg" alt="Photo">
    <img src="https://example.com/remote.png" alt="Remote">
    <img src="data:image/png;base64,ABC123" alt="Data URI">
</body>
</html>`

	p := NewHTMLParser()
	doc, err := p.Parse([]byte(html), ".")

	require.NoError(t, err)

	// Count image resources (only local images should be tracked)
	var imageCount int
	for _, r := range doc.Resources {
		if strings.HasPrefix(r.MediaType, "image/") {
			imageCount++
		}
	}

	// Should have 2 local images (local.png and photo.jpg)
	assert.Equal(t, 2, imageCount)
}

func TestHTMLParser_Parse_RewritesImagePaths(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<body>
    <img src="images/photo.jpg" alt="Photo">
    <img src="https://example.com/remote.png" alt="Remote">
</body>
</html>`

	p := NewHTMLParser()
	doc, err := p.Parse([]byte(html), ".")

	require.NoError(t, err)

	content := doc.Chapters[0].Content

	// Local image paths should be rewritten
	assert.Contains(t, content, "../images/photo.jpg")

	// Remote URLs should be preserved
	assert.Contains(t, content, "https://example.com/remote.png")
}

func TestHTMLParser_Parse_NoBody(t *testing.T) {
	// HTML without body tag
	html := `<h1>Title</h1><p>Content</p>`

	p := NewHTMLParser()
	doc, err := p.Parse([]byte(html), ".")

	require.NoError(t, err)
	assert.NotNil(t, doc)
	// Should still produce output using the entire content
	assert.NotEmpty(t, doc.Chapters)
}

func TestHTMLParser_Parse_UppercaseTags(t *testing.T) {
	html := `<HTML>
<BODY>
    <H1>Test</H1>
    <BR>
    <HR>
    <IMG src="test.png">
</BODY>
</HTML>`

	p := NewHTMLParser()
	doc, err := p.Parse([]byte(html), ".")

	require.NoError(t, err)

	content := doc.Chapters[0].Content
	// golang.org/x/net/html normalizes tags to lowercase
	assert.Contains(t, content, "<h1>")
	assert.Contains(t, content, "<br />")
	assert.Contains(t, content, "<hr />")
}

func TestHTMLParser_SupportedExtensions(t *testing.T) {
	p := NewHTMLParser()
	exts := p.SupportedExtensions()

	assert.Contains(t, exts, ".html")
	assert.Contains(t, exts, ".htm")
}

func TestHTMLParser_Parse_MultipleAuthors(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Multi-Author</title>
    <meta name="author" content="Author One">
    <meta name="author" content="Author Two">
</head>
<body>
    <h1>Test</h1>
</body>
</html>`

	p := NewHTMLParser()
	doc, err := p.Parse([]byte(html), ".")

	require.NoError(t, err)
	assert.Len(t, doc.Metadata.Authors, 2)
	assert.Contains(t, doc.Metadata.Authors, "Author One")
	assert.Contains(t, doc.Metadata.Authors, "Author Two")
}

func TestHTMLParser_Parse_HeadingIDs(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<body>
    <h1 id="custom-id">Custom ID Heading</h1>
    <h2>Auto Generated ID</h2>
</body>
</html>`

	p := NewHTMLParser()
	doc, err := p.Parse([]byte(html), ".")

	require.NoError(t, err)

	// Check that custom ID is preserved in TOC
	entries := doc.TOC.Entries
	assert.NotEmpty(t, entries)
	assert.Contains(t, entries[0].Href, "#custom-id")
}
