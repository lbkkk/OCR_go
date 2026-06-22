// Command goocr converts PDF/Word documents to Markdown.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/lbkkk/OCR_go/internal/delivery/cli"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := cli.Execute(context.Background()); err != nil {
		logger.Error("goocr failed", "error", err)
		os.Exit(1)
	}
}
