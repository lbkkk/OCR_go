// Package cli defines the command-line interface for goocr.
package cli

import (
	"context"

	"github.com/spf13/cobra"
)

// newRootCmd creates the root "goocr" command. Subcommands (e.g. convert)
// are attached to it.
func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "goocr",
		Short: "Convert PDF/Word documents to Markdown",
		Long: "goocr is a command-line tool that converts PDF and Word (.docx) " +
			"documents to Markdown, supporting text extraction and OCR for scanned/image documents.",
		SilenceUsage: true,
	}

	return root
}

// Execute runs the root command with the context passed from main.
func Execute(ctx context.Context) error {
	return newRootCmd().ExecuteContext(ctx)
}
