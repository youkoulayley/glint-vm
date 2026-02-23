package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/youkoulayley/glint-vm/internal/shell"
)

// initCommand generates shell initialization code.
func initCommand(_ context.Context, cmd *cli.Command) error {
	// Get shell name from arguments (first non-flag argument)
	shellName := cmd.Args().First()
	if shellName == "" {
		shellName = shell.DetectShell()
	}

	integrator, err := shell.NewIntegrator(shellName)
	if err != nil {
		return fmt.Errorf("failed to create shell integrator: %w", err)
	}

	opts := shell.InitOptions{
		AutoSwitch: cmd.Bool("auto-switch"),
	}

	output := integrator.GenerateInit(opts)
	fmt.Print(output)

	return nil
}
