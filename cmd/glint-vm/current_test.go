package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
	"github.com/youkoulayley/glint-vm/internal/config"
)

func TestCurrentCommand_NoVersion(t *testing.T) {
	t.Parallel()

	_, cleanup := setupTestEnv(t)
	defer cleanup()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "current",
				Action: currentCommand,
			},
		},
	}

	output := captureOutput(func() {
		_ = app.Run([]string{"glint-vm", "current"})
	})

	if !strings.Contains(output, "No version currently active") {
		t.Errorf("Output should indicate no active version, got: %s", output)
	}
}

func TestCurrentCommand_WithVersion(t *testing.T) {
	t.Parallel()

	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping symlink test on Windows")
	}

	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	cfg := &config.Config{
		CacheDir: filepath.Join(tmpDir, "glint-vm"),
		OS:       "linux",
		Arch:     "amd64",
	}

	version := "v1.55.2"
	if err := cfg.EnsureVersionDir(version); err != nil {
		t.Fatalf("Failed to create version dir: %v", err)
	}

	binaryPath := cfg.GetBinaryPath(version)

	file, err := os.Create(binaryPath) //nolint:gosec // Test file creation
	if err != nil {
		t.Fatalf("Failed to create binary: %v", err)
	}

	_ = file.Close()

	if err := os.Chmod(binaryPath, 0755); err != nil { //nolint:gosec // Test file permissions
		t.Fatalf("Failed to chmod binary: %v", err)
	}

	if err := cfg.SetCurrentVersion(version); err != nil {
		t.Fatalf("Failed to set current version: %v", err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "current",
				Action: currentCommand,
			},
		},
	}

	output := captureOutput(func() {
		_ = app.Run([]string{"glint-vm", "current"})
	})

	if !strings.Contains(output, "v1.55.2") {
		t.Errorf("Output should contain version v1.55.2, got: %s", output)
	}

	if !strings.Contains(output, "Current version:") {
		t.Errorf("Output should indicate current version, got: %s", output)
	}
}
