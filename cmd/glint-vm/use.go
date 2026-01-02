package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/youkoulayley/glint-vm/internal/config"
	"github.com/youkoulayley/glint-vm/internal/downloader"
	"github.com/youkoulayley/glint-vm/internal/shell"
)

// useCommand activates a specific version in the current shell.
func useCommand(c *cli.Context) error {
	if c.NArg() < 1 {
		return ErrVersionRequired
	}

	version := config.NormalizeVersion(c.Args().First())

	dl, err := downloader.NewDownloader()
	if err != nil {
		return fmt.Errorf("failed to initialize downloader: %w", err)
	}

	if err := dl.Download(version); err != nil {
		return fmt.Errorf("failed to download golangci-lint: %w", err)
	}

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

	return nil
}
