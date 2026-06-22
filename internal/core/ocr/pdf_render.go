package ocr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultRenderDPI = 200

// PDFRenderer converts PDF pages to PNG images via pdftoppm (Poppler).
type PDFRenderer struct {
	Bin string // empty => lookup "pdftoppm" in PATH
	DPI int
}

// NewPDFRenderer creates a renderer. bin may be empty to use PATH.
func NewPDFRenderer(bin string, dpi int) (*PDFRenderer, error) {
	resolved, err := LookPath(bin, "pdftoppm")
	if err != nil {
		return nil, err
	}
	if dpi <= 0 {
		dpi = defaultRenderDPI
	}
	return &PDFRenderer{Bin: resolved, DPI: dpi}, nil
}

// RenderPage renders one 1-based PDF page to PNG bytes.
func (r *PDFRenderer) RenderPage(ctx context.Context, pdfPath string, pageNum int) ([]byte, error) {
	if pageNum < 1 {
		return nil, fmt.Errorf("pdftoppm: page number must be >= 1, got %d", pageNum)
	}

	dir, err := os.MkdirTemp("", "goocr-pdf-*")
	if err != nil {
		return nil, fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	outPrefix := filepath.Join(dir, "page")
	args := []string{
		"-singlefile",
		"-png",
		"-f", fmt.Sprintf("%d", pageNum),
		"-l", fmt.Sprintf("%d", pageNum),
		"-r", fmt.Sprintf("%d", r.DPI),
		pdfPath,
		outPrefix,
	}

	if _, _, err := runCommand(ctx, r.Bin, args...); err != nil {
		return nil, fmt.Errorf("pdftoppm page %d: %w", pageNum, err)
	}

	pngPath := outPrefix + ".png"
	data, err := os.ReadFile(pngPath)
	if err != nil {
		return nil, fmt.Errorf("read rendered png: %w", err)
	}
	return data, nil
}

// RenderAllPages renders every page in the PDF (1..N) to PNG bytes.
func (r *PDFRenderer) RenderAllPages(ctx context.Context, pdfPath string, numPages int) ([][]byte, error) {
	images := make([][]byte, 0, numPages)
	for i := 1; i <= numPages; i++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		img, err := r.RenderPage(ctx, pdfPath, i)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

// ReadImageFile loads a standalone image file for OCR.
func ReadImageFile(path string) ([]byte, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".tif", ".tiff", ".bmp":
	default:
		return nil, fmt.Errorf("unsupported image extension: %s", ext)
	}
	return os.ReadFile(path)
}
