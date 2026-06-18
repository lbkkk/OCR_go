---
name: goocr-dev
description: Development guide for the go-ocr CLI that converts PDF/Word to Markdown. Use when implementing, extending, or debugging go-ocr packages (detector, extractor, ocr, llm, enricher, markdown, converter) or when adding a new input format or OCR/LLM backend.
disable-model-invocation: true
---

# go-ocr Development Guide

CLI tool converting PDF/Word (.docx) to Markdown. Module: `github.com/lbkkk/OCR_go`.

## Conversion pipeline

```
input -> detector -> extractor | ocr -> enricher (LLM) -> markdown renderer -> .md
```

The `pkg/document` model is the single intermediate representation. Extraction/OCR produces a `Document`; the renderer consumes it.

## Phased plan

1. Phase 0: environment (Go, Tesseract, Poppler, LLM endpoint).
2. Phase 1: scaffold (go.mod, main.go, cobra root, slog). DONE.
3. Phase 2: `pkg/document` model + interfaces (`Extractor`, `OCREngine`, `Renderer`, `ImageDescriber`, `Refiner`).
4. Phase 3: `internal/llm` OpenAI-compatible client (+ optional vision).
5. Phase 4: `internal/detector` (file type + PDF text-layer check).
6. Phase 5: `internal/extractor/pdf` + `internal/extractor/docx`.
7. Phase 6: `internal/ocr` (`pdftoppm` + `tesseract` via os/exec).
8. Phase 7: `internal/enricher` (LLM refine text + describe images, with fallback).
9. Phase 8: `internal/markdown` renderer.
10. Phase 9: `internal/converter` + CLI `convert` command and flags.
11. Phase 10: tests (mock LLM/OCR). Phase 11: build/CI/Docker/docs.

## Extending

- New input format: add an `Extractor` implementation under `internal/extractor/<fmt>`; register it in the converter by file type. Do not add format logic to the converter.
- New OCR or LLM backend: implement the existing interface (`OCREngine` / LLM client); the converter stays unchanged.
- Hybrid OCR: prefer the text layer / Tesseract result; use the LLM to refine and as a fallback when OCR confidence is low.
- Image handling: send embedded images to `ImageDescriber` (LLM vision) for a description. If the model is text-only, fall back to OCR of the image or a placeholder.

## Conventions

- English for all code/comments/CLI text. `gofmt`-clean, wrap errors with `%w`, propagate `context.Context`.
- Secrets (LLM API key) come from env/config, never hardcoded.
- Verify with `go build ./...` and `go test ./...` after changes.
