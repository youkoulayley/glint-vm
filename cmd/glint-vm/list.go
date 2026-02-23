package main

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/youkoulayley/glint-vm/internal/config"
	"github.com/youkoulayley/glint-vm/internal/downloader"
)

// listCommand lists all installed versions.
func listCommand(_ *cli.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	cm, err := downloader.NewCacheManager()
	if err != nil {
		return fmt.Errorf("failed to initialize cache manager: %w", err)
	}

	versions, err := cm.List()
	if err != nil {
		return fmt.Errorf("failed to list versions: %w", err)
	}

	if len(versions) == 0 {
		fmt.Println("No versions installed.")
		fmt.Println()
		fmt.Println("Install a version with:")
		fmt.Println("  glint-vm install v1.55.2")

		return nil
	}

	currentVersion, _ := cfg.GetCurrentVersion()

	fmt.Println("Installed versions:")

	for _, version := range versions {
		marker := " "
		if version.Version == currentVersion {
			marker = "*"
		}

		status := ""
		if !version.IsComplete {
			status = " (incomplete)"
		}

		fmt.Printf("%s %s%s\n", marker, version.Version, status)
	}

	if currentVersion != "" {
		fmt.Println()
		fmt.Printf("* = current version (%s)\n", currentVersion)
	}

	return nil
}
