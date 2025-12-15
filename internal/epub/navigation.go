package epub

import (
	"bytes"
	"html/template"

	"github.com/dauquangthanh/epub-converter/internal/model"
)

// navTemplate is the template for nav.xhtml
const navTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="{{.Language}}" lang="{{.Language}}">
<head>
  <meta charset="UTF-8"/>
  <title>{{.Title}}</title>
  <link rel="stylesheet" type="text/css" href="styles/default.css"/>
</head>
<body>
  <nav epub:type="toc" id="toc">
    <h1>Table of Contents</h1>
{{.TOCList}}
  </nav>
  <nav epub:type="landmarks" id="landmarks" hidden="">
    <h2>Landmarks</h2>
    <ol>
      <li><a epub:type="toc" href="nav.xhtml">Table of Contents</a></li>
{{- if .HasContent}}
      <li><a epub:type="bodymatter" href="{{.FirstChapterHref}}">Start of Content</a></li>
{{- end}}
    </ol>
  </nav>
</body>
</html>`

// navData holds data for the navigation template
type navData struct {
	Language         string
	Title            string
	TOCList          template.HTML
	HasContent       bool
	FirstChapterHref string
}

// generateNavDocument generates the nav.xhtml file content.
func generateNavDocument(doc *model.Document) (string, error) {
	tmpl, err := template.New("nav").Parse(navTemplate)
	if err != nil {
		return "", err
	}

	tocList := renderTOCList(doc.TOC.Entries)

	var firstChapter string
	if len(doc.Chapters) > 0 {
		firstChapter = doc.Chapters[0].FileName
	}

	data := navData{
		Language:         doc.Metadata.Language,
		Title:            doc.Metadata.Title,
		TOCList:          template.HTML(tocList),
		HasContent:       len(doc.Chapters) > 0,
		FirstChapterHref: firstChapter,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// renderTOCList renders the TOC entries as nested ordered lists.
func renderTOCList(entries []model.TOCEntry) string {
	if len(entries) == 0 {
		return "    <ol></ol>"
	}

	var buf bytes.Buffer
	buf.WriteString("    <ol>\n")
	for _, entry := range entries {
		renderTOCEntry(&buf, entry, 3)
	}
	buf.WriteString("    </ol>")
	return buf.String()
}

// renderTOCEntry renders a single TOC entry with its children.
func renderTOCEntry(buf *bytes.Buffer, entry model.TOCEntry, indent int) {
	indentStr := spaces(indent)

	// Escape HTML in title
	escapedTitle := template.HTMLEscapeString(entry.Title)

	buf.WriteString(indentStr)
	buf.WriteString("<li>\n")
	buf.WriteString(indentStr)
	buf.WriteString("  <a href=\"")
	buf.WriteString(entry.Href)
	buf.WriteString("\">")
	buf.WriteString(escapedTitle)
	buf.WriteString("</a>\n")

	if len(entry.Children) > 0 {
		buf.WriteString(indentStr)
		buf.WriteString("  <ol>\n")
		for _, child := range entry.Children {
			renderTOCEntry(buf, child, indent+2)
		}
		buf.WriteString(indentStr)
		buf.WriteString("  </ol>\n")
	}

	buf.WriteString(indentStr)
	buf.WriteString("</li>\n")
}

// spaces returns a string of n spaces for indentation.
func spaces(n int) string {
	s := make([]byte, n*2)
	for i := range s {
		s[i] = ' '
	}
	return string(s)
}
