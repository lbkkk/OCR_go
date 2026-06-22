package docx

import (
	"archive/zip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/lbkkk/OCR_go/pkg/document"
)

const documentXMLPath = "word/document.xml"
const mediaPrefix = "word/media/"

// Extractor parses DOCX files (Office Open XML zip) into document.Document.
type Extractor struct{}

// NewExtractor returns a DOCX extractor.
func NewExtractor() *Extractor {
	return &Extractor{}
}

// Supports reports whether the path is a DOCX file.
func (e *Extractor) Supports(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".docx")
}

// Extract reads paragraphs and embedded images from a DOCX file.
func (e *Extractor) Extract(ctx context.Context, path string) (*document.Document, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open docx zip: %w", err)
	}
	defer zr.Close()

	doc := &document.Document{
		Source: path,
		Metadata: map[string]string{
			"format": "docx",
		},
	}

	for _, f := range zr.File {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		switch {
		case f.Name == documentXMLPath:
			paras, err := readParagraphs(f)
			if err != nil {
				return nil, fmt.Errorf("parse document.xml: %w", err)
			}
			for _, p := range paras {
				if strings.TrimSpace(p) == "" {
					continue
				}
				doc.Blocks = append(doc.Blocks, document.Paragraph{Text: p})
			}
		case strings.HasPrefix(f.Name, mediaPrefix) && !strings.HasSuffix(f.Name, "/"):
			img, err := readEmbeddedImage(f)
			if err != nil {
				return nil, fmt.Errorf("read %s: %w", f.Name, err)
			}
			doc.Blocks = append(doc.Blocks, img)
		}
	}

	doc.Metadata["blocks"] = fmt.Sprintf("%d", len(doc.Blocks))
	return doc, nil
}

func readParagraphs(f *zip.File) ([]string, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return parseParagraphs(rc)
}

// parseParagraphs walks document.xml and collects text from w:p / w:t elements.
// Namespace prefixes are ignored; only local element names matter.
func parseParagraphs(r io.Reader) ([]string, error) {
	dec := xml.NewDecoder(r)
	var paragraphs []string
	var current strings.Builder
	inParagraph := false
	inText := false

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch se := tok.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "p":
				inParagraph = true
				current.Reset()
			case "t":
				inText = true
			}
		case xml.CharData:
			if inParagraph && inText {
				current.Write(se)
			}
		case xml.EndElement:
			switch se.Name.Local {
			case "t":
				inText = false
			case "p":
				if inParagraph {
					paragraphs = append(paragraphs, strings.TrimSpace(current.String()))
					inParagraph = false
					inText = false
				}
			}
		}
	}

	return paragraphs, nil
}

func readEmbeddedImage(f *zip.File) (document.Image, error) {
	rc, err := f.Open()
	if err != nil {
		return document.Image{}, err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return document.Image{}, err
	}

	base := filepath.Base(f.Name)
	format := strings.TrimPrefix(strings.ToLower(filepath.Ext(base)), ".")
	if format == "jpg" {
		format = "jpeg"
	}

	return document.Image{
		Data:   data,
		Format: format,
		Source: f.Name,
	}, nil
}
