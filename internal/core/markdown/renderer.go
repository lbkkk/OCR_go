// Package markdown defines the contract for rendering a document model to
// Markdown output.
package markdown

import "github.com/lbkkk/OCR_go/pkg/document"

// Renderer converts a document.Document into Markdown bytes.
type Renderer interface {
	Render(doc *document.Document) ([]byte, error)
}
