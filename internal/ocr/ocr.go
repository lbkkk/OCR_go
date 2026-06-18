// Package ocr defines the contract for optical character recognition engines
// used on scanned PDFs and image inputs.
package ocr

import "context"

// Engine recognizes text from a single image.
type Engine interface {
	// Recognize runs OCR over the given image bytes and returns the text.
	// lang selects the language(s), e.g. "eng" or "eng+vie".
	Recognize(ctx context.Context, image []byte, lang string) (string, error)
}
