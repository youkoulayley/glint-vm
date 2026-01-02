package main

import (
	"bytes"
	"os"
	"testing"
)

// setupTestEnv creates a temporary test environment.
func setupTestEnv(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir := t.TempDir()

	oldXDG := os.Getenv("XDG_CACHE_HOME")

	t.Setenv("XDG_CACHE_HOME", tmpDir)

	cleanup := func() {
		if oldXDG != "" {
			t.Setenv("XDG_CACHE_HOME", oldXDG)
		}
	}

	return tmpDir, cleanup
}

func captureOutput(f func()) string {
	old := os.Stdout
	pipe, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer

	_, _ = buf.ReadFrom(pipe)

	return buf.String()
}
