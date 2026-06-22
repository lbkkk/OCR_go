package detector

import (
	"context"
	"testing"
)

func TestKindFromExtension(t *testing.T) {
	tests := []struct {
		path string
		want Kind
	}{
		{"report.pdf", KindPDF},
		{"doc.DOCX", KindDOCX},
		{"photo.png", KindImage},
		{"unknown.xyz", KindUnknown},
	}
	for _, tt := range tests {
		if got := kindFromExtension(tt.path); got != tt.want {
			t.Errorf("kindFromExtension(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestDocClassFromPages(t *testing.T) {
	tests := []struct {
		name  string
		pages []PageKind
		want  DocClass
	}{
		{"all digital", []PageKind{PageDigital, PageDigital}, ClassAllDigital},
		{"all scan", []PageKind{PageScan, PageScan}, ClassAllScan},
		{"mixed", []PageKind{PageDigital, PageScan, PageDigital}, ClassMixed},
		{"empty", nil, ClassAllScan},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := docClassFromPages(tt.pages); got != tt.want {
				t.Errorf("docClassFromPages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountSignificantRunes(t *testing.T) {
	if got := countSignificantRunes("  hello  \n world "); got != 10 {
		t.Errorf("count = %d, want 10", got)
	}
}

func TestDetectUnsupported(t *testing.T) {
	d := New()
	_, err := d.Detect(context.Background(), "file.unknown")
	if err == nil {
		t.Fatal("expected error for unsupported extension")
	}
}
