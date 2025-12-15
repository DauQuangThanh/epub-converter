package epub

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dauquangthanh/epub-converter/internal/model"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	assert.NotNil(t, builder)
}

func TestBuilder_Build_MinimalDocument(t *testing.T) {
	builder := NewBuilder()

	doc := model.NewDocument()
	doc.Metadata.Title = "Test Book"
	doc.AddChapter(model.Chapter{
		ID:       "ch1",
		Title:    "Chapter 1",
		Content:  "<p>Test content</p>",
		FileName: "content/chapter-001.xhtml",
	})

	data, err := builder.Build(doc)

	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Verify it's a valid ZIP
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	// Check required files exist
	fileNames := make(map[string]bool)
	for _, f := range reader.File {
		fileNames[f.Name] = true
	}

	assert.True(t, fileNames["mimetype"], "mimetype missing")
	assert.True(t, fileNames["META-INF/container.xml"], "container.xml missing")
	assert.True(t, fileNames["OEBPS/content.opf"], "content.opf missing")
	assert.True(t, fileNames["OEBPS/nav.xhtml"], "nav.xhtml missing")
}

func TestBuilder_Build_MimetypeFirst(t *testing.T) {
	builder := NewBuilder()

	doc := model.NewDocument()
	doc.Metadata.Title = "Test"
	doc.AddChapter(model.Chapter{
		ID:       "ch1",
		Title:    "Test",
		Content:  "<p>Test</p>",
		FileName: "content/chapter-001.xhtml",
	})

	data, err := builder.Build(doc)
	require.NoError(t, err)

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	// First file must be mimetype
	assert.Equal(t, "mimetype", reader.File[0].Name)
}

func TestBuilder_Build_MimetypeUncompressed(t *testing.T) {
	builder := NewBuilder()

	doc := model.NewDocument()
	doc.Metadata.Title = "Test"
	doc.AddChapter(model.Chapter{
		ID:       "ch1",
		Title:    "Test",
		Content:  "<p>Test</p>",
		FileName: "content/chapter-001.xhtml",
	})

	data, err := builder.Build(doc)
	require.NoError(t, err)

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	// Mimetype must be stored uncompressed
	assert.Equal(t, zip.Store, reader.File[0].Method)
}

func TestBuilder_Build_WithCoverImage(t *testing.T) {
	builder := NewBuilder()

	doc := model.NewDocument()
	doc.Metadata.Title = "Book with Cover"
	doc.AddChapter(model.Chapter{
		ID:       "ch1",
		Title:    "Chapter 1",
		Content:  "<p>Content</p>",
		FileName: "content/chapter-001.xhtml",
	})
	doc.AddResource(model.Resource{
		ID:        "cover-image",
		FileName:  "images/cover.jpg",
		MediaType: "image/jpeg",
		Data:      []byte{0xFF, 0xD8, 0xFF, 0xE0}, // JPEG header
		IsCover:   true,
	})

	data, err := builder.Build(doc)
	require.NoError(t, err)

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	// Check cover image is included
	fileNames := make(map[string]bool)
	for _, f := range reader.File {
		fileNames[f.Name] = true
	}

	assert.True(t, fileNames["OEBPS/images/cover.jpg"])
}

func TestBuilder_Build_MultipleChapters(t *testing.T) {
	builder := NewBuilder()

	doc := model.NewDocument()
	doc.Metadata.Title = "Multi-Chapter Book"

	for i := 1; i <= 3; i++ {
		doc.AddChapter(model.Chapter{
			ID:       "ch" + string(rune('0'+i)),
			Title:    "Chapter",
			Content:  "<p>Content</p>",
			FileName: "content/chapter-00" + string(rune('0'+i)) + ".xhtml",
			Order:    i - 1,
		})
	}

	data, err := builder.Build(doc)
	require.NoError(t, err)

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	fileNames := make(map[string]bool)
	for _, f := range reader.File {
		fileNames[f.Name] = true
	}

	// All chapters should be included
	for i := 1; i <= 3; i++ {
		fileName := "OEBPS/content/chapter-00" + string(rune('0'+i)) + ".xhtml"
		assert.True(t, fileNames[fileName], "Missing: "+fileName)
	}
}
