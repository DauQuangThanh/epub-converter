// ------------------------------------------------------------------
// Developed by Dau Quang Thanh - 2025.
// Enterprise AI Solution Architect
//
// Happy Reading!
// ------------------------------------------------------------------

package epub

import (
	"github.com/dauquangthanh/epub-converter/internal/model"
)

// MergeMetadata combines source metadata with CLI overrides.
// CLI overrides take precedence over source metadata.
func MergeMetadata(source, cli *model.Metadata) *model.Metadata {
	result := model.NewMetadata()

	// Start with source values
	if source != nil {
		result.Title = source.Title
		result.Authors = append(result.Authors, source.Authors...)
		result.Language = source.Language
		result.Identifier = source.Identifier
		result.Description = source.Description
		result.Publisher = source.Publisher
		result.Date = source.Date
		result.Rights = source.Rights
		result.CoverImage = source.CoverImage
	}

	// Override with CLI values if provided
	if cli != nil {
		result.Merge(cli)
	}

	// Ensure defaults are set
	result.EnsureDefaults()

	return result
}

// ValidateMetadata checks that required metadata fields are present.
// Returns nil if valid, otherwise returns an error describing the issue.
func ValidateMetadata(meta *model.Metadata) error {
	if !meta.Valid() {
		return ErrMissingTitle
	}
	return nil
}
