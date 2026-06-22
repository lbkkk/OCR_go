// Package detector classifies input files and, for PDFs, determines whether
// each page has a text layer (digital) or is image-only (scan).
package detector

import "context"

// Kind is the file format at the top level.
type Kind int

const (
	KindUnknown Kind = iota
	KindPDF
	KindDOCX
	KindImage // standalone .png / .jpg / ...
)

// DocClass summarizes a PDF after scanning every page's text layer.
type DocClass int

const (
	ClassUnknown DocClass = iota
	ClassAllDigital // every page has extractable text
	ClassAllScan    // no page has text layer
	ClassMixed      // some digital pages, some scan pages
)

// PageKind classifies a single PDF page.
type PageKind int

const (
	PageUnknown PageKind = iota
	PageDigital
	PageScan
)

// Result is returned by Detect. For PDFs, Pages holds one entry per page
// (1-based index matches slice index + 1).
type Result struct {
	Path  string
	Kind  Kind
	Class DocClass
	Pages []PageKind
}

// Detector inspects a file and returns classification metadata.
type Detector interface {
	Detect(ctx context.Context, path string) (Result, error)
}
