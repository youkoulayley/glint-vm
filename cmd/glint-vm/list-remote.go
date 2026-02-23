package main

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/youkoulayley/glint-vm/internal/config"
	"github.com/youkoulayley/glint-vm/internal/downloader"
)

// listRemoteCommand lists available versions from GitHub.
func listRemoteCommand(c *cli.Context) error {
	limit := c.Int("limit")

	fmt.Println("Fetching available golangci-lint versions from GitHub...")
	fmt.Println()

	releases, err := downloader.FetchAvailableVersions(limit)
	if err != nil {
		return fmt.Errorf("failed to fetch versions: %w", err)
	}

	if len(releases) == 0 {
		fmt.Println("No stable releases found.")

		return nil
	}

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	currentVersion, err := cfg.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	fmt.Printf("Available versions (%d latest stable releases):\n", len(releases))

	for _, release := range releases {
		marker := " "
		status := ""

		if release.TagName == currentVersion {
			marker = "*"
			status = " (current)"
		} else if cfg.BinaryExists(release.TagName) {
			status = " (installed)"
		}

		fmt.Printf("%s %s%s\n", marker, release.TagName, status)
	}

	return nil
}
