// Package document defines the intermediate representation shared between
// extraction/OCR and Markdown rendering.
package document

// Document is the parsed, source-agnostic representation of an input file.
type Document struct {
	// Title is an optional document title (e.g. from metadata or first heading).
	Title string
	// Source is the original file path the document was produced from.
	Source string
	// Blocks holds the ordered content of the document.
	Blocks []Block
	// Metadata carries optional key/value information (author, page count, ...).
	Metadata map[string]string
}

// Block is a single unit of content. It is a sealed interface: only the
// concrete types declared in this package implement it.
type Block interface {
	isBlock()
}

// Heading is a section title with a level from 1 (top) to 6.
type Heading struct {
	Level int
	Text  string
}

func (Heading) isBlock() {}

// Paragraph is a block of running text.
type Paragraph struct {
	Text string
}

func (Paragraph) isBlock() {}

// ListItem is a single entry in a List, with an optional nesting level.
type ListItem struct {
	Text  string
	Level int
}

// List is an ordered or unordered list.
type List struct {
	Ordered bool
	Items   []ListItem
}

func (List) isBlock() {}

// Table is a simple grid with an optional header row.
type Table struct {
	Header []string
	Rows   [][]string
}

func (Table) isBlock() {}

// Image is an embedded image together with any derived text.
type Image struct {
	// Data holds the raw image bytes; may be empty if only Source is known.
	Data []byte
	// Format is the image encoding, e.g. "png" or "jpeg".
	Format string
	// AltText is short alternative text (e.g. OCR'd text inside the image).
	AltText string
	// Description is a natural-language description, typically from an LLM.
	Description string
	// Source is the original reference or path of the image.
	Source string
}

func (Image) isBlock() {}

// CodeBlock is preformatted/code text with an optional language hint.
type CodeBlock struct {
	Language string
	Text     string
}

func (CodeBlock) isBlock() {}
