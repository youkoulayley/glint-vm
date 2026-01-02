package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/youkoulayley/glint-vm/internal/shell"
)

// initCommand generates shell initialization code.
func initCommand(ctx *cli.Context) error {
	// Get shell name from arguments (first non-flag argument)
	shellName := ctx.Args().First()
	if shellName == "" {
		shellName = shell.DetectShell()
	}

	integrator, err := shell.NewIntegrator(shellName)
	if err != nil {
		return fmt.Errorf("failed to create shell integrator: %w", err)
	}

	opts := shell.InitOptions{
		AutoSwitch: ctx.Bool("auto-switch"),
	}

	output := integrator.GenerateInit(opts)
	fmt.Print(output)

	return nil
}
