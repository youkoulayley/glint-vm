package main

import (
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestListCommand_NoVersions(t *testing.T) { //nolint:paralleltest // uses t.Setenv via setupTestEnv
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "list",
				Action: listCommand,
			},
		},
	}

	output := captureOutput(func() {
		_ = app.Run([]string{"glint-vm", "list"})
	})

	if !strings.Contains(output, "No versions installed") {
		t.Errorf("Output should indicate no versions, got: %s", output)
	}
}
