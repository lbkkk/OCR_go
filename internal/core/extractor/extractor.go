// Package extractor defines the contract for parsing digital documents
// (PDF, DOCX, ...) into the shared document model.
package extractor

import (
	"context"

	"github.com/lbkkk/OCR_go/pkg/document"
)

// Extractor parses a single input file into a document.Document.
type Extractor interface {
	// Supports reports whether this extractor can handle the given file.
	Supports(path string) bool
	// Extract parses the file at path into a Document.
	Extract(ctx context.Context, path string) (*document.Document, error)
}
