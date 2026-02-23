package main

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/youkoulayley/glint-vm/internal/config"
)

// currentCommand shows the currently active version.
func currentCommand(_ *cli.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	version, err := cfg.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if version == "" {
		fmt.Println("No version currently active.")
		fmt.Println()
		fmt.Println("Activate a version with:")
		fmt.Println("  glint-vm use v1.55.2")

		return nil
	}

	fmt.Printf("Current version: %s\n", version)

	binaryPath := cfg.GetBinaryPath(version)
	if cfg.BinaryExists(version) {
		fmt.Printf("Binary path: %s\n", binaryPath)
	}

	return nil
}
