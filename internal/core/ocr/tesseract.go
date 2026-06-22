package ocr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TesseractEngine runs OCR via the tesseract CLI.
type TesseractEngine struct {
	Bin string // empty => lookup "tesseract" in PATH or TESSERACT_PATH
}

// NewTesseractEngine creates a TesseractEngine. bin may be empty to use PATH.
func NewTesseractEngine(bin string) (*TesseractEngine, error) {
	resolved, err := LookPath(bin, "tesseract")
	if err != nil {
		return nil, err
	}
	return &TesseractEngine{Bin: resolved}, nil
}

// Recognize writes image bytes to a temp file and runs tesseract stdout mode.
func (e *TesseractEngine) Recognize(ctx context.Context, image []byte, lang string) (string, error) {
	if len(image) == 0 {
		return "", fmt.Errorf("tesseract: empty image")
	}
	if lang == "" {
		lang = "eng"
	}

	dir, err := os.MkdirTemp("", "goocr-tess-*")
	if err != nil {
		return "", fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	imgPath := filepath.Join(dir, "page.png")
	if err := os.WriteFile(imgPath, image, 0o600); err != nil {
		return "", fmt.Errorf("write image: %w", err)
	}

	// tesseract image stdout -l eng
	stdout, _, err := runCommand(ctx, e.Bin, imgPath, "stdout", "-l", lang)
	if err != nil {
		return "", fmt.Errorf("tesseract: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}
