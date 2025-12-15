// ------------------------------------------------------------------
// Developed by Dau Quang Thanh - 2025.
// Enterprise AI Solution Architect
//
// Happy Reading!
// ------------------------------------------------------------------

package model

// TableOfContents represents the navigation hierarchy for the EPUB.
type TableOfContents struct {
	Entries []TOCEntry
}

// TOCEntry is a single navigation item in the table of contents.
type TOCEntry struct {
	Title    string     // Display text for the entry
	Href     string     // Link to content (e.g., "chapter-01.xhtml")
	Level    int        // Hierarchy depth (1-6)
	Children []TOCEntry // Nested entries for sub-sections
}

// NewTableOfContents creates a new empty TableOfContents.
func NewTableOfContents() *TableOfContents {
	return &TableOfContents{
		Entries: make([]TOCEntry, 0),
	}
}

// AddEntry appends a top-level entry to the table of contents.
func (t *TableOfContents) AddEntry(entry TOCEntry) {
	t.Entries = append(t.Entries, entry)
}

// Empty returns true if the table of contents has no entries.
func (t *TableOfContents) Empty() bool {
	return len(t.Entries) == 0
}

// FlatEntries returns all entries in a flat list (depth-first order).
func (t *TableOfContents) FlatEntries() []TOCEntry {
	var result []TOCEntry
	for _, entry := range t.Entries {
		result = append(result, flattenEntry(entry)...)
	}
	return result
}

// flattenEntry recursively flattens an entry and its children.
func flattenEntry(entry TOCEntry) []TOCEntry {
	result := []TOCEntry{entry}
	for _, child := range entry.Children {
		result = append(result, flattenEntry(child)...)
	}
	return result
}

// BuildFromHeadings creates a hierarchical TOC from a flat list of headings.
// Each heading should have a Level (1-6) indicating its depth.
func BuildFromHeadings(headings []TOCEntry) *TableOfContents {
	toc := NewTableOfContents()
	if len(headings) == 0 {
		return toc
	}

	// Simple algorithm: nest entries based on level
	var stack []*TOCEntry

	for i := range headings {
		entry := headings[i]
		entry.Children = make([]TOCEntry, 0)

		// Find parent for this entry
		for len(stack) > 0 && stack[len(stack)-1].Level >= entry.Level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			// Top-level entry
			toc.Entries = append(toc.Entries, entry)
			stack = append(stack, &toc.Entries[len(toc.Entries)-1])
		} else {
			// Child of current parent
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, entry)
			stack = append(stack, &parent.Children[len(parent.Children)-1])
		}
	}

	return toc
}
