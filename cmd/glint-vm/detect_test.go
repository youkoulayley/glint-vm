package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestDetectCommand_NoVersion(t *testing.T) { //nolint:paralleltest // uses t.Setenv via setupTestEnv
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	projectDir := filepath.Join(tmpDir, "test-project")

	err := os.MkdirAll(projectDir, 0750)
	if err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	oldWd, _ := os.Getwd()

	t.Chdir(projectDir)
	t.Chdir(oldWd)

	app := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:   "detect",
				Action: detectCommand,
			},
		},
	}

	output := captureOutput(func() {
		_ = app.Run(context.Background(), []string{"glint-vm", "detect"})
	})

	if !strings.Contains(output, "No golangci-lint version detected") {
		t.Errorf("Output should indicate no version detected, got: %s", output)
	}
}
