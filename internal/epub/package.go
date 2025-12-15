package epub

import (
	"bytes"
	"html/template"
	"time"

	"github.com/dauquangthanh/epub-converter/internal/model"
)

// packageTemplate is the template for content.opf
const packageTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="uid">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="uid">{{.Identifier}}</dc:identifier>
    <dc:title>{{.Title}}</dc:title>
    <dc:language>{{.Language}}</dc:language>
{{- range .Authors}}
    <dc:creator>{{.}}</dc:creator>
{{- end}}
{{- if .Description}}
    <dc:description>{{.Description}}</dc:description>
{{- end}}
{{- if .Publisher}}
    <dc:publisher>{{.Publisher}}</dc:publisher>
{{- end}}
{{- if .Rights}}
    <dc:rights>{{.Rights}}</dc:rights>
{{- end}}
    <dc:date>{{.Date}}</dc:date>
    <meta property="dcterms:modified">{{.Modified}}</meta>
  </metadata>
  <manifest>
    <item id="nav" href="nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
    <item id="css" href="styles/default.css" media-type="text/css"/>
{{- range .Chapters}}
    <item id="{{.ID}}" href="{{.FileName}}" media-type="application/xhtml+xml"/>
{{- end}}
{{- range .Resources}}
    <item id="{{.ID}}" href="{{.FileName}}" media-type="{{.MediaType}}"{{if .IsCover}} properties="cover-image"{{end}}/>
{{- end}}
  </manifest>
  <spine>
{{- range .Chapters}}
    <itemref idref="{{.ID}}"/>
{{- end}}
  </spine>
</package>`

// packageData holds data for the package template
type packageData struct {
	Identifier  string
	Title       string
	Language    string
	Authors     []string
	Description string
	Publisher   string
	Rights      string
	Date        string
	Modified    string
	Chapters    []model.Chapter
	Resources   []model.Resource
}

// generatePackageDocument generates the content.opf file content.
func generatePackageDocument(doc *model.Document) (string, error) {
	tmpl, err := template.New("package").Parse(packageTemplate)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	date := doc.Metadata.Date.Format("2006-01-02")

	data := packageData{
		Identifier:  doc.Metadata.Identifier,
		Title:       doc.Metadata.Title,
		Language:    doc.Metadata.Language,
		Authors:     doc.Metadata.Authors,
		Description: doc.Metadata.Description,
		Publisher:   doc.Metadata.Publisher,
		Rights:      doc.Metadata.Rights,
		Date:        date,
		Modified:    now,
		Chapters:    doc.Chapters,
		Resources:   doc.Resources,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
