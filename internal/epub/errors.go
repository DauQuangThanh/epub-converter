package epub

import "errors"

// EPUB generation errors
var (
	ErrMissingTitle    = errors.New("missing required title metadata")
	ErrNoChapters      = errors.New("document has no chapters")
	ErrInvalidDocument = errors.New("invalid document")
)
