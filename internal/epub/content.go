// ------------------------------------------------------------------
// Developed by Dau Quang Thanh - 2025.
// Enterprise AI Solution Architect
//
// Happy Reading!
// ------------------------------------------------------------------

package epub

import (
	"bytes"
	"html"
	"text/template"

	"github.com/dauquangthanh/epub-converter/internal/model"
)

// contentTemplate is the template for XHTML content documents
const contentTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
  <meta charset="UTF-8"/>
  <title>{{.Title}}</title>
  <link rel="stylesheet" type="text/css" href="styles/default.css"/>
</head>
<body epub:type="bodymatter">
{{.Content}}
</body>
</html>`

// contentData holds data for the content template
type contentData struct {
	Title   string
	Content string
}

// generateContentDocument generates an XHTML content document.
func generateContentDocument(chapter *model.Chapter, bookTitle string) (string, error) {
	tmpl, err := template.New("content").Parse(contentTemplate)
	if err != nil {
		return "", err
	}

	title := chapter.Title
	if title == "" {
		title = bookTitle
	}

	// Escape title for XML safety, but content is already HTML
	data := contentData{
		Title:   html.EscapeString(title),
		Content: chapter.Content,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
