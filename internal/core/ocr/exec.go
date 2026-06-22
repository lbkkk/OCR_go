package ocr

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

// LookPath resolves an executable: explicit path if set, otherwise PATH lookup.
func LookPath(explicit, defaultName string) (string, error) {
	if explicit != "" {
		if _, err := os.Stat(explicit); err != nil {
			return "", fmt.Errorf("binary %q: %w", explicit, err)
		}
		return explicit, nil
	}
	path, err := exec.LookPath(defaultName)
	if err != nil {
		return "", fmt.Errorf("%q not found in PATH: install it or set an explicit path", defaultName)
	}
	return path, nil
}

// runCommand executes a command with context cancellation.
func runCommand(ctx context.Context, name string, args ...string) (stdout, stderr []byte, err error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return outBuf.Bytes(), errBuf.Bytes(), fmt.Errorf("%s %v: %w\nstderr: %s", name, args, err, errBuf.String())
	}
	return outBuf.Bytes(), errBuf.Bytes(), nil
}
