package detector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const magicReadSize = 512

// FileDetector implements Detector using extension, magic bytes, and PDF page scan.
type FileDetector struct{}

// New returns a ready-to-use FileDetector.
func New() *FileDetector {
	return &FileDetector{}
}

// Detect classifies the file at path.
func (d *FileDetector) Detect(ctx context.Context, path string) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return Result{}, fmt.Errorf("stat %q: %w", path, err)
	}
	if info.IsDir() {
		return Result{}, fmt.Errorf("path %q is a directory", path)
	}

	kind := kindFromExtension(path)
	if kind == KindUnknown {
		kind = kindFromMagic(path)
	}

	res := Result{Path: path, Kind: kind}

	switch kind {
	case KindPDF:
		pages, class, err := classifyPDFPages(ctx, path)
		if err != nil {
			return Result{}, err
		}
		res.Pages = pages
		res.Class = class
	case KindDOCX, KindImage:
		// No per-page scan needed; downstream always uses extractor or OCR directly.
		res.Class = ClassUnknown
	default:
		return Result{}, fmt.Errorf("unsupported file type: %s", filepath.Base(path))
	}

	return res, nil
}

func kindFromExtension(path string) Kind {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".pdf":
		return KindPDF
	case ".docx":
		return KindDOCX
	case ".png", ".jpg", ".jpeg", ".webp", ".tif", ".tiff", ".bmp":
		return KindImage
	default:
		return KindUnknown
	}
}

func kindFromMagic(path string) Kind {
	f, err := os.Open(path)
	if err != nil {
		return KindUnknown
	}
	defer f.Close()

	buf := make([]byte, magicReadSize)
	n, _ := f.Read(buf)
	if n == 0 {
		return KindUnknown
	}
	head := buf[:n]

	if strings.HasPrefix(string(head), "%PDF") {
		return KindPDF
	}
	if n >= 4 && head[0] == 0x50 && head[1] == 0x4B && head[2] == 0x03 && head[3] == 0x04 {
		// ZIP container — could be DOCX; confirm by extension elsewhere or treat as docx if .docx only from ext
		if strings.EqualFold(filepath.Ext(path), ".docx") {
			return KindDOCX
		}
	}
	if n >= 8 && head[0] == 0x89 && head[1] == 0x50 && head[2] == 0x4E && head[3] == 0x47 {
		return KindImage
	}
	if n >= 3 && head[0] == 0xFF && head[1] == 0xD8 && head[2] == 0xFF {
		return KindImage
	}

	return KindUnknown
}
