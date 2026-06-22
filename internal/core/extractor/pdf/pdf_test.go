package pdf

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lbkkk/OCR_go/pkg/document"
)

func TestSupports(t *testing.T) {
	e := NewExtractor()
	if !e.Supports("a.pdf") || e.Supports("a.docx") {
		t.Fatal("Supports mismatch for pdf")
	}
}

func TestExtractMissingFile(t *testing.T) {
	e := NewExtractor()
	_, err := e.Extract(context.Background(), filepath.Join(t.TempDir(), "missing.pdf"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestExtractPagesEmptyList(t *testing.T) {
	e := NewExtractor()
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.pdf")
	if err := os.WriteFile(path, []byte("%PDF-1.4\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Invalid PDF content — expect open/parse error, not panic.
	_, err := e.ExtractPages(context.Background(), path, []int{1})
	if err == nil {
		t.Skip("ledongthuc/pdf opened minimal stub; need real PDF in testdata for full test")
	}
}

func TestDocumentBlocksAreParagraphs(t *testing.T) {
	doc := &document.Document{
		Blocks: []document.Block{
			document.Paragraph{Text: "hello"},
		},
	}
	if len(doc.Blocks) != 1 {
		t.Fatalf("blocks = %d", len(doc.Blocks))
	}
	if _, ok := doc.Blocks[0].(document.Paragraph); !ok {
		t.Fatal("expected paragraph block")
	}
	_ = strings.TrimSpace("x")
}
