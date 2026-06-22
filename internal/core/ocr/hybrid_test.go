package ocr

import (
	"context"
	"errors"
	"testing"
)

type stubEngine struct {
	out string
	err error
}

func (s stubEngine) Recognize(ctx context.Context, image []byte, lang string) (string, error) {
	return s.out, s.err
}

type stubQwen struct {
	out string
	err error
}

func (s stubQwen) RecognizeWithDraft(ctx context.Context, image []byte, draft, lang string) (string, error) {
	return s.out, s.err
}

func TestHybridEngineTesseractOnly(t *testing.T) {
	h := NewHybridEngine(stubEngine{out: "draft only"}, stubQwen{out: "refined"}, nil, false)
	got, err := h.Recognize(context.Background(), []byte{1}, "eng")
	if err != nil {
		t.Fatal(err)
	}
	if got != "draft only" {
		t.Fatalf("got %q, want draft only", got)
	}
}

func TestHybridEngineQwenSuccess(t *testing.T) {
	h := NewHybridEngine(stubEngine{out: "draft"}, stubQwen{out: "refined"}, nil, true)
	got, err := h.Recognize(context.Background(), []byte{1}, "eng")
	if err != nil {
		t.Fatal(err)
	}
	if got != "refined" {
		t.Fatalf("got %q, want refined", got)
	}
}

func TestHybridEngineQwenFallback(t *testing.T) {
	h := NewHybridEngine(stubEngine{out: "draft"}, stubQwen{err: errors.New("api down")}, nil, true)
	got, err := h.Recognize(context.Background(), []byte{1}, "eng")
	if err != nil {
		t.Fatal(err)
	}
	if got != "draft" {
		t.Fatalf("got %q, want tesseract fallback", got)
	}
}

func TestHybridEngineNilTess(t *testing.T) {
	h := NewHybridEngine(nil, nil, nil, true)
	_, err := h.Recognize(context.Background(), []byte{1}, "eng")
	if err == nil {
		t.Fatal("expected error when tess is nil")
	}
}

func TestQwenRecognizeWithDraftNilClient(t *testing.T) {
	qwen := &QwenEngine{Client: nil}
	_, err := qwen.RecognizeWithDraft(context.Background(), []byte{1}, "draft", "eng")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}
