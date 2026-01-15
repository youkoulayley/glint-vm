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

	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer

	_, _ = buf.ReadFrom(r)

	return buf.String()
}
