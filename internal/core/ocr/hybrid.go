package ocr

import (
	"context"
	"log/slog"
)

// HybridEngine orchestrates Tesseract pass 1 and Qwen vision pass 2 with fallback.
type HybridEngine struct {
	Tess      Engine
	Qwen      DraftRefiner
	Logger    *slog.Logger
	UseHybrid bool // when false, only Tesseract pass 1 runs
}

// NewHybridEngine wires Tesseract and optional Qwen engines.
func NewHybridEngine(tess Engine, qwen DraftRefiner, logger *slog.Logger, useHybrid bool) *HybridEngine {
	if logger == nil {
		logger = slog.Default()
	}
	return &HybridEngine{
		Tess:      tess,
		Qwen:      qwen,
		Logger:    logger,
		UseHybrid: useHybrid,
	}
}

// Recognize runs hybrid OCR: Tesseract draft -> Qwen refine, falling back to draft on Qwen error.
func (h *HybridEngine) Recognize(ctx context.Context, image []byte, lang string) (string, error) {
	if h.Tess == nil {
		return "", errNoTesseract
	}

	draft, err := h.Tess.Recognize(ctx, image, lang)
	if err != nil {
		return "", err
	}

	if !h.UseHybrid || h.Qwen == nil {
		h.Logger.Info("ocr pass complete", "ocr_pass", "tesseract")
		return draft, nil
	}

	refined, err := h.Qwen.RecognizeWithDraft(ctx, image, draft, lang)
	if err != nil {
		h.Logger.Warn("qwen failed, using tesseract draft", "ocr_pass", "qwen_fallback", "err", err)
		return draft, nil
	}

	h.Logger.Info("ocr pass complete", "ocr_pass", "qwen")
	return refined, nil
}
