package main

import (
	"context"
	"strings"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestListCommand_NoVersions(t *testing.T) { //nolint:paralleltest // uses t.Setenv via setupTestEnv
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	app := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:   "list",
				Action: listCommand,
			},
		},
	}

	output := captureOutput(func() {
		_ = app.Run(context.Background(), []string{"glint-vm", "list"})
	})

	if !strings.Contains(output, "No versions installed") {
		t.Errorf("Output should indicate no versions, got: %s", output)
	}
}
