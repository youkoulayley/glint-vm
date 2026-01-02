package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/youkoulayley/glint-vm/internal/config"
	"github.com/youkoulayley/glint-vm/internal/downloader"
	"github.com/youkoulayley/glint-vm/internal/shell"
)

// installCommand pre-downloads a specific version.
func installCommand(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return ErrVersionRequired
	}

	version := config.NormalizeVersion(ctx.Args().First())

	dl, err := downloader.NewDownloader()
	if err != nil {
		return fmt.Errorf("failed to initialize downloader: %w", err)
	}

	if err := dl.Download(version); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	if ctx.Bool("use") {
		cfg, err := config.New()
		if err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}

		if err := cfg.SetCurrentVersion(version); err != nil {
			return fmt.Errorf("failed to set current version: %w", err)
		}

		shellName := shell.DetectShell()

		integrator, err := shell.NewIntegrator(shellName)
		if err != nil {
			return fmt.Errorf("failed to create shell integrator: %w", err)
		}

		output := integrator.GenerateUse(version)
		fmt.Print(output)
	} else {
		fmt.Fprintf(os.Stderr, "✓ Installed golangci-lint %s\n", version)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "To activate this version, run:")
		fmt.Fprintf(os.Stderr, "  glint-vm use %s\n", version)
	}

	return nil
}
