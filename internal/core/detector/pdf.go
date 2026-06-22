package detector

import (
	"context"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/ledongthuc/pdf"
)

// minTextRunes is the minimum number of non-space runes on a page to treat it as digital.
const minTextRunes = 8

// classifyPDFPages opens the PDF and labels each page as PageDigital or PageScan.
func classifyPDFPages(ctx context.Context, path string) ([]PageKind, DocClass, error) {
	if err := ctx.Err(); err != nil {
		return nil, ClassUnknown, err
	}

	f, r, err := pdf.Open(path)
	if err != nil {
		return nil, ClassUnknown, fmt.Errorf("open pdf %q: %w", path, err)
	}
	defer f.Close()

	return classifyReaderPages(ctx, r)
}

func classifyReaderPages(ctx context.Context, r *pdf.Reader) ([]PageKind, DocClass, error) {
	numPages := r.NumPage()
	if numPages == 0 {
		return nil, ClassAllScan, nil
	}

	pages := make([]PageKind, numPages)
	for i := 1; i <= numPages; i++ {
		if err := ctx.Err(); err != nil {
			return nil, ClassUnknown, err
		}
		pages[i-1] = classifyPage(r, i)
	}

	return pages, docClassFromPages(pages), nil
}

func classifyPage(r *pdf.Reader, pageNum int) PageKind {
	page := r.Page(pageNum)
	if page.V.IsNull() {
		return PageScan
	}

	text, err := page.GetPlainText(nil)
	if err != nil {
		return PageScan
	}
	if countSignificantRunes(text) >= minTextRunes {
		return PageDigital
	}
	return PageScan
}

func countSignificantRunes(s string) int {
	n := 0
	for _, r := range s {
		if !unicode.IsSpace(r) {
			n++
		}
	}
	return n
}

func docClassFromPages(pages []PageKind) DocClass {
	if len(pages) == 0 {
		return ClassAllScan
	}

	hasDigital := false
	hasScan := false
	for _, p := range pages {
		switch p {
		case PageDigital:
			hasDigital = true
		case PageScan:
			hasScan = true
		}
	}

	switch {
	case hasDigital && hasScan:
		return ClassMixed
	case hasDigital:
		return ClassAllDigital
	case hasScan:
		return ClassAllScan
	default:
		return ClassAllScan
	}
}

// PageText returns plain text for a single PDF page (1-based page number).
func PageText(r *pdf.Reader, pageNum int) (string, error) {
	page := r.Page(pageNum)
	if page.V.IsNull() {
		return "", nil
	}
	text, err := page.GetPlainText(nil)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

// OpenPDF opens a PDF file and returns the OS file handle and reader.
// Caller must close the returned file.
func OpenPDF(path string) (*os.File, *pdf.Reader, error) {
	return pdf.Open(path)
}
