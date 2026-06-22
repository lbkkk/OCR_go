package pdf

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"

	"github.com/lbkkk/OCR_go/internal/core/detector"
	"github.com/lbkkk/OCR_go/pkg/document"
)

// Extractor reads text from digital PDF pages.
type Extractor struct{}

// NewExtractor returns a PDF text extractor.
func NewExtractor() *Extractor {
	return &Extractor{}
}

// Supports reports whether the path is a PDF file.
func (e *Extractor) Supports(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".pdf")
}

// Extract reads text from every page (for ClassAllDigital fast path).
func (e *Extractor) Extract(ctx context.Context, path string) (*document.Document, error) {
	f, r, err := detector.OpenPDF(path)
	if err != nil {
		return nil, fmt.Errorf("open pdf: %w", err)
	}
	defer f.Close()

	numPages := r.NumPage()
	pages := make([]int, numPages)
	for i := range pages {
		pages[i] = i + 1
	}
	return e.extractPages(ctx, path, r, pages)
}

// ExtractPages reads text from specific 1-based page numbers (Mixed PDF fast path per page).
func (e *Extractor) ExtractPages(ctx context.Context, path string, pageNums []int) (*document.Document, error) {
	f, r, err := detector.OpenPDF(path)
	if err != nil {
		return nil, fmt.Errorf("open pdf: %w", err)
	}
	defer f.Close()

	return e.extractPages(ctx, path, r, pageNums)
}

// ExtractPageText returns plain text for one 1-based page using an already-open reader.
func ExtractPageText(ctx context.Context, r *pdf.Reader, pageNum int) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return detector.PageText(r, pageNum)
}

func (e *Extractor) extractPages(ctx context.Context, path string, r *pdf.Reader, pageNums []int) (*document.Document, error) {
	doc := &document.Document{
		Source: path,
		Metadata: map[string]string{
			"format": "pdf",
		},
	}

	for _, pageNum := range pageNums {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		text, err := detector.PageText(r, pageNum)
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", pageNum, err)
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		doc.Blocks = append(doc.Blocks, document.Paragraph{Text: text})
	}

	doc.Metadata["pages_extracted"] = fmt.Sprintf("%d", len(pageNums))
	return doc, nil
}
