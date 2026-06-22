package ocr

import (
	"context"
	"fmt"

	"github.com/lbkkk/OCR_go/internal/core/detector"
	pdfextract "github.com/lbkkk/OCR_go/internal/core/extractor/pdf"
	"github.com/lbkkk/OCR_go/pkg/document"
)

// PDFProcessor routes PDF pages through extractor (digital) or hybrid OCR (scan).
type PDFProcessor struct {
	Extractor *pdfextract.Extractor
	Renderer  *PDFRenderer
	Engine    Engine
}

// NewPDFProcessor wires digital extraction and scan OCR dependencies.
func NewPDFProcessor(extractor *pdfextract.Extractor, renderer *PDFRenderer, engine Engine) *PDFProcessor {
	return &PDFProcessor{
		Extractor: extractor,
		Renderer:  renderer,
		Engine:    engine,
	}
}

// Process handles AllDigital, AllScan, and Mixed PDFs using detector.Result.
func (p *PDFProcessor) Process(ctx context.Context, path string, det detector.Result, lang string) (*document.Document, error) {
	if p.Extractor == nil || p.Renderer == nil || p.Engine == nil {
		return nil, fmt.Errorf("pdf processor: missing dependency")
	}
	if det.Kind != detector.KindPDF {
		return nil, fmt.Errorf("pdf processor: expected PDF, got kind %v", det.Kind)
	}

	switch det.Class {
	case detector.ClassAllDigital:
		return p.Extractor.Extract(ctx, path)
	case detector.ClassAllScan:
		return p.processAllPagesOCR(ctx, path, det, lang)
	case detector.ClassMixed:
		return p.processMixed(ctx, path, det, lang)
	default:
		return nil, fmt.Errorf("pdf processor: unknown class %v", det.Class)
	}
}

func (p *PDFProcessor) processAllPagesOCR(ctx context.Context, path string, det detector.Result, lang string) (*document.Document, error) {
	doc := &document.Document{
		Source: path,
		Metadata: map[string]string{
			"format": "pdf",
			"class":  "all_scan",
		},
	}

	for pageNum := 1; pageNum <= len(det.Pages); pageNum++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		text, err := p.ocrPage(ctx, path, pageNum, lang)
		if err != nil {
			return nil, fmt.Errorf("page %d ocr: %w", pageNum, err)
		}
		if text == "" {
			continue
		}
		doc.Blocks = append(doc.Blocks, document.Paragraph{Text: text})
	}
	return doc, nil
}

func (p *PDFProcessor) processMixed(ctx context.Context, path string, det detector.Result, lang string) (*document.Document, error) {
	doc := &document.Document{
		Source: path,
		Metadata: map[string]string{
			"format": "pdf",
			"class":  "mixed",
		},
	}

	f, r, err := detector.OpenPDF(path)
	if err != nil {
		return nil, fmt.Errorf("open pdf: %w", err)
	}
	defer f.Close()

	for pageNum, kind := range det.Pages {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		oneBased := pageNum + 1

		var text string
		switch kind {
		case detector.PageDigital:
			text, err = pdfextract.ExtractPageText(ctx, r, oneBased)
			if err != nil {
				return nil, fmt.Errorf("page %d extract: %w", oneBased, err)
			}
		case detector.PageScan:
			text, err = p.ocrPage(ctx, path, oneBased, lang)
			if err != nil {
				return nil, fmt.Errorf("page %d ocr: %w", oneBased, err)
			}
		default:
			continue
		}
		if text == "" {
			continue
		}
		doc.Blocks = append(doc.Blocks, document.Paragraph{Text: text})
	}
	return doc, nil
}

func (p *PDFProcessor) ocrPage(ctx context.Context, pdfPath string, pageNum int, lang string) (string, error) {
	img, err := p.Renderer.RenderPage(ctx, pdfPath, pageNum)
	if err != nil {
		return "", err
	}
	return p.Engine.Recognize(ctx, img, lang)
}

// RecognizeImageFile OCRs a standalone image file.
func RecognizeImageFile(ctx context.Context, engine Engine, path, lang string) (string, error) {
	data, err := ReadImageFile(path)
	if err != nil {
		return "", err
	}
	return engine.Recognize(ctx, data, lang)
}
