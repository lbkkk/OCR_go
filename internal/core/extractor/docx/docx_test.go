package docx

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

const minimalDocumentXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:r><w:t>Hello DOCX</w:t></w:r></w:p>
    <w:p><w:r><w:t>Second paragraph</w:t></w:r></w:p>
  </w:body>
</w:document>`

func TestSupports(t *testing.T) {
	e := NewExtractor()
	if !e.Supports("file.docx") || e.Supports("file.pdf") {
		t.Fatal("Supports mismatch")
	}
}

func TestExtractMinimalDOCX(t *testing.T) {
	path := writeMinimalDOCX(t)
	e := NewExtractor()
	doc, err := e.Extract(context.Background(), path)
	if err != nil {
		t.Fatalf("Extract: %v", err)
	}
	if len(doc.Blocks) != 2 {
		t.Fatalf("blocks = %d, want 2", len(doc.Blocks))
	}
}

func writeMinimalDOCX(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.docx")

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	w, err := zw.Create(documentXMLPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte(minimalDocumentXML)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseParagraphs(t *testing.T) {
	paras, err := parseParagraphs(bytes.NewReader([]byte(minimalDocumentXML)))
	if err != nil {
		t.Fatal(err)
	}
	if len(paras) != 2 || paras[0] != "Hello DOCX" {
		t.Fatalf("paragraphs = %#v", paras)
	}
}
