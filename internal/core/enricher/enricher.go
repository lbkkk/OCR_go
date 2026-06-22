// Package enricher defines the contracts for LLM-based enrichment: refining
// raw text into Markdown-ready text and describing embedded images.
package enricher

import (
	"context"
	"errors"

	"github.com/lbkkk/OCR_go/pkg/document"
)

// ErrVisionUnsupported is returned by an ImageDescriber whose backing model
// cannot process images, allowing callers to fall back to OCR or a placeholder.
var ErrVisionUnsupported = errors.New("enricher: vision not supported by the model")

// Refiner cleans up and structures raw extracted/OCR text.
type Refiner interface {
	// Refine turns raw text into clean, Markdown-ready text (fixing OCR
	// artifacts and recovering structure where possible).
	Refine(ctx context.Context, raw string) (string, error)
}

// ImageDescriber produces a natural-language description of an image.
type ImageDescriber interface {
	// Describe returns a description of the given image. Implementations that
	// lack vision support should return ErrVisionUnsupported.
	Describe(ctx context.Context, img document.Image) (string, error)
}
