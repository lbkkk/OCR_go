package ocr

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lbkkk/OCR_go/internal/adapter/llm"
)

const defaultQwenSystemPrompt = "You are a document OCR correction assistant. Fix Tesseract errors using the page image. Do not invent content."

const defaultUserPromptPrefix = "Correct the OCR draft below. Output clean Markdown preserving structure (headings, tables, lists). OCR draft:\n"

// DraftRefiner improves a Tesseract draft using vision (pass 2).
type DraftRefiner interface {
	RecognizeWithDraft(ctx context.Context, image []byte, draft, lang string) (string, error)
}

// QwenEngine refines OCR drafts using a vision-capable LLM (Qwen via OpenAI-compatible API).
type QwenEngine struct {
	Client       *llm.Client
	SystemPrompt string
	UserPrefix   string
	Timeout      time.Duration
	ImageFormat  string // e.g. "png"
}

// NewQwenEngine returns a QwenEngine with sensible defaults.
func NewQwenEngine(client *llm.Client, opts ...QwenOption) *QwenEngine {
	e := &QwenEngine{
		Client:       client,
		SystemPrompt: defaultQwenSystemPrompt,
		UserPrefix:   defaultUserPromptPrefix,
		Timeout:      60 * time.Second,
		ImageFormat:  "png",
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// QwenOption customizes QwenEngine.
type QwenOption func(*QwenEngine)

// WithQwenTimeout sets the LLM call timeout for pass 2.
func WithQwenTimeout(d time.Duration) QwenOption {
	return func(e *QwenEngine) { e.Timeout = d }
}

// WithImageFormat sets the MIME subtype sent to the vision API (default "png").
func WithImageFormat(format string) QwenOption {
	return func(e *QwenEngine) { e.ImageFormat = format }
}

// RecognizeWithDraft sends image + Tesseract draft to the vision model.
func (e *QwenEngine) RecognizeWithDraft(ctx context.Context, image []byte, draft, lang string) (string, error) {
	if e.Client == nil {
		return "", fmt.Errorf("qwen: nil llm client")
	}
	_ = lang // reserved for future prompt tuning

	ctx, cancel := context.WithTimeout(ctx, e.Timeout)
	defer cancel()

	format := e.ImageFormat
	if format == "" {
		format = "png"
	}

	messages := []llm.Message{
		llm.System(e.SystemPrompt),
		llm.UserWithImages(e.UserPrefix+draft, []llm.Image{{Data: image, Format: format}}),
	}

	out, err := e.Client.Complete(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("qwen complete: %w", err)
	}
	return strings.TrimSpace(out), nil
}

// Recognize satisfies Engine but requires a draft; use HybridEngine instead.
func (e *QwenEngine) Recognize(ctx context.Context, image []byte, lang string) (string, error) {
	return "", fmt.Errorf("qwen: use RecognizeWithDraft or HybridEngine")
}
