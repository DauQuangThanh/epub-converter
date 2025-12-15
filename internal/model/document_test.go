package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDocument(t *testing.T) {
	doc := NewDocument()

	assert.NotNil(t, doc)
	assert.Empty(t, doc.Chapters)
	assert.Empty(t, doc.Resources)
	assert.NotNil(t, doc.TOC.Entries)
}

func TestDocument_AddChapter(t *testing.T) {
	doc := NewDocument()

	chapter := Chapter{
		ID:      "ch1",
		Title:   "Chapter 1",
		Content: "<p>Content</p>",
	}

	doc.AddChapter(chapter)

	assert.Len(t, doc.Chapters, 1)
	assert.Equal(t, "ch1", doc.Chapters[0].ID)
}

func TestDocument_AddResource(t *testing.T) {
	doc := NewDocument()

	resource := Resource{
		ID:        "img1",
		FileName:  "images/test.png",
		MediaType: "image/png",
		Data:      []byte{0x89, 0x50, 0x4E, 0x47},
	}

	doc.AddResource(resource)

	assert.Len(t, doc.Resources, 1)
	assert.Equal(t, "img1", doc.Resources[0].ID)
}

func TestMetadata_Merge(t *testing.T) {
	base := &Metadata{
		Title:    "Original Title",
		Authors:  []string{"Author 1"},
		Language: "en",
	}

	override := &Metadata{
		Title:   "New Title",
		Authors: []string{"Author 2"},
	}

	base.Merge(override)

	assert.Equal(t, "New Title", base.Title)
	assert.Contains(t, base.Authors, "Author 2")
}

func TestMetadata_Merge_EmptyOverride(t *testing.T) {
	base := &Metadata{
		Title:    "Original",
		Language: "en",
	}

	override := &Metadata{}
	base.Merge(override)

	assert.Equal(t, "Original", base.Title)
	assert.Equal(t, "en", base.Language)
}

func TestBuildFromHeadings(t *testing.T) {
	entries := []TOCEntry{
		{Title: "Chapter 1", Href: "ch1.xhtml", Level: 1},
		{Title: "Section 1.1", Href: "ch1.xhtml#s1", Level: 2},
		{Title: "Section 1.2", Href: "ch1.xhtml#s2", Level: 2},
		{Title: "Chapter 2", Href: "ch2.xhtml", Level: 1},
	}

	toc := BuildFromHeadings(entries)

	assert.NotNil(t, toc)
	assert.NotEmpty(t, toc.Entries)
}

func TestResource_IsCoverImage(t *testing.T) {
	resource := Resource{
		ID:        "cover",
		FileName:  "images/cover.jpg",
		MediaType: "image/jpeg",
		IsCover:   true,
	}

	assert.True(t, resource.IsCover)
}

func TestChapter_Properties(t *testing.T) {
	chapter := Chapter{
		ID:       "ch-001",
		Title:    "Introduction",
		Level:    1,
		Content:  "<h1>Introduction</h1><p>Content</p>",
		FileName: "content/chapter-001.xhtml",
		Order:    0,
	}

	assert.Equal(t, "ch-001", chapter.ID)
	assert.Equal(t, "Introduction", chapter.Title)
	assert.Equal(t, 1, chapter.Level)
	assert.Contains(t, chapter.Content, "Introduction")
	assert.Contains(t, chapter.FileName, "chapter-001")
	assert.Equal(t, 0, chapter.Order)
}
