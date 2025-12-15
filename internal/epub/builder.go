// ------------------------------------------------------------------
// Developed by Dau Quang Thanh - 2025.
// Enterprise AI Solution Architect
//
// Happy Reading!
// ------------------------------------------------------------------

// Package epub provides EPUB 3+ package generation functionality.
package epub

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"github.com/dauquangthanh/epub-converter/internal/model"
)

// Builder creates valid EPUB 3+ packages from Document models.
type Builder struct {
	doc *model.Document
}

// NewBuilder creates a new EPUB builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// Build generates an EPUB file from the document and returns the bytes.
func (b *Builder) Build(doc *model.Document) ([]byte, error) {
	b.doc = doc

	// Ensure document has required metadata
	doc.Metadata.EnsureDefaults()

	if !doc.Valid() {
		return nil, fmt.Errorf("invalid document: missing title or chapters")
	}

	// Add colophon page at the end
	b.addColophon(doc)

	var buf bytes.Buffer
	if err := b.writeEPUB(&buf); err != nil {
		return nil, fmt.Errorf("building EPUB: %w", err)
	}

	return buf.Bytes(), nil
}

// WriteToFile generates an EPUB file and writes it to the specified writer.
func (b *Builder) WriteToFile(doc *model.Document, w io.Writer) error {
	data, err := b.Build(doc)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// writeEPUB creates the complete EPUB archive.
func (b *Builder) writeEPUB(w io.Writer) error {
	zw := zip.NewWriter(w)
	defer zw.Close()

	// 1. Write mimetype first (must be uncompressed and first entry)
	if err := b.writeMimetype(zw); err != nil {
		return fmt.Errorf("writing mimetype: %w", err)
	}

	// 2. Write META-INF/container.xml
	if err := b.writeContainer(zw); err != nil {
		return fmt.Errorf("writing container.xml: %w", err)
	}

	// 3. Write OEBPS/content.opf (package document)
	if err := b.writePackageDocument(zw); err != nil {
		return fmt.Errorf("writing content.opf: %w", err)
	}

	// 4. Write OEBPS/nav.xhtml (navigation document)
	if err := b.writeNavDocument(zw); err != nil {
		return fmt.Errorf("writing nav.xhtml: %w", err)
	}

	// 5. Write OEBPS/content/*.xhtml (content documents)
	if err := b.writeContentDocuments(zw); err != nil {
		return fmt.Errorf("writing content documents: %w", err)
	}

	// 6. Write resources (images, stylesheets)
	if err := b.writeResources(zw); err != nil {
		return fmt.Errorf("writing resources: %w", err)
	}

	// 7. Write default stylesheet
	if err := b.writeDefaultStylesheet(zw); err != nil {
		return fmt.Errorf("writing stylesheet: %w", err)
	}

	return nil
}

// writeMimetype writes the mimetype file (must be first, uncompressed).
func (b *Builder) writeMimetype(zw *zip.Writer) error {
	// Create file header with no compression
	header := &zip.FileHeader{
		Name:   "mimetype",
		Method: zip.Store, // No compression
	}
	w, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("application/epub+zip"))
	return err
}

// writeContainer writes META-INF/container.xml.
func (b *Builder) writeContainer(zw *zip.Writer) error {
	w, err := zw.Create("META-INF/container.xml")
	if err != nil {
		return err
	}

	container := `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>`

	_, err = w.Write([]byte(container))
	return err
}

// writePackageDocument writes OEBPS/content.opf.
func (b *Builder) writePackageDocument(zw *zip.Writer) error {
	w, err := zw.Create("OEBPS/content.opf")
	if err != nil {
		return err
	}

	opf, err := generatePackageDocument(b.doc)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(opf))
	return err
}

// writeNavDocument writes OEBPS/nav.xhtml.
func (b *Builder) writeNavDocument(zw *zip.Writer) error {
	w, err := zw.Create("OEBPS/nav.xhtml")
	if err != nil {
		return err
	}

	nav, err := generateNavDocument(b.doc)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(nav))
	return err
}

// writeContentDocuments writes OEBPS/content/*.xhtml files.
func (b *Builder) writeContentDocuments(zw *zip.Writer) error {
	for _, chapter := range b.doc.Chapters {
		path := "OEBPS/" + chapter.FileName
		w, err := zw.Create(path)
		if err != nil {
			return err
		}

		content, err := generateContentDocument(&chapter, b.doc.Metadata.Title)
		if err != nil {
			return err
		}

		if _, err := w.Write([]byte(content)); err != nil {
			return err
		}
	}
	return nil
}

// writeResources writes embedded resources (images, etc.).
func (b *Builder) writeResources(zw *zip.Writer) error {
	for _, resource := range b.doc.Resources {
		path := "OEBPS/" + resource.FileName
		w, err := zw.Create(path)
		if err != nil {
			return err
		}
		if _, err := w.Write(resource.Data); err != nil {
			return err
		}
	}
	return nil
}

// writeDefaultStylesheet writes a basic stylesheet.
func (b *Builder) writeDefaultStylesheet(zw *zip.Writer) error {
	w, err := zw.Create("OEBPS/styles/default.css")
	if err != nil {
		return err
	}

	css := `/* Default EPUB stylesheet */
body {
  font-family: serif;
  line-height: 1.6;
  margin: 1em;
}

h1, h2, h3, h4, h5, h6 {
  font-family: sans-serif;
  line-height: 1.2;
  margin-top: 1.5em;
  margin-bottom: 0.5em;
}

h1 { font-size: 2em; }
h2 { font-size: 1.5em; }
h3 { font-size: 1.25em; }
h4 { font-size: 1.1em; }
h5 { font-size: 1em; }
h6 { font-size: 0.9em; }

p {
  margin: 0.5em 0;
  text-align: justify;
}

pre, code {
  font-family: monospace;
  font-size: 0.9em;
}

pre {
  background-color: #f5f5f5;
  padding: 1em;
  overflow-x: auto;
  border-radius: 4px;
}

code {
  background-color: #f5f5f5;
  padding: 0.1em 0.3em;
  border-radius: 2px;
}

pre code {
  background-color: transparent;
  padding: 0;
}

blockquote {
  margin: 1em 2em;
  padding-left: 1em;
  border-left: 3px solid #ccc;
  font-style: italic;
}

ul, ol {
  margin: 0.5em 0;
  padding-left: 2em;
}

li {
  margin: 0.25em 0;
}

table {
  border-collapse: collapse;
  width: 100%;
  margin: 1em 0;
}

th, td {
  border: 1px solid #ccc;
  padding: 0.5em;
  text-align: left;
}

th {
  background-color: #f5f5f5;
  font-weight: bold;
}

img {
  max-width: 100%;
  height: auto;
}

a {
  color: #0066cc;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

/* Task list styling */
.task-list {
  list-style-type: none;
  padding-left: 0;
}

.task-list-item {
  display: flex;
  align-items: flex-start;
}

.task-list-item input {
  margin-right: 0.5em;
}
`

	_, err = w.Write([]byte(css))
	return err
}

// addColophon adds an attribution page at the end of the book.
func (b *Builder) addColophon(doc *model.Document) {
	colophonContent := `<hr style="margin: 3em 0;"/>
<div style="text-align: center; font-family: monospace; white-space: pre-wrap; padding: 2em 1em; background-color: #f9f9f9; border: 1px solid #ddd; margin: 2em 0;">
------------------------------------------------------------------
Packaged by Epub Converter Application (c) 2025 Dau Quang Thanh.

URL: <a href="https://github.com/DauQuangThanh/epub-converter">https://github.com/DauQuangThanh/epub-converter</a>

Happy Reading!
------------------------------------------------------------------
</div>`

	colophon := model.Chapter{
		ID:       "colophon",
		Title:    "About This EPUB",
		Level:    1,
		Content:  colophonContent,
		FileName: "content/colophon.xhtml",
		Order:    len(doc.Chapters),
	}

	doc.AddChapter(colophon)
}
