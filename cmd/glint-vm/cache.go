package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/youkoulayley/glint-vm/internal/config"
	"github.com/youkoulayley/glint-vm/internal/downloader"
)

const (
	kilobyte = 1024
)

// cacheListCommand lists all cached versions.
func cacheListCommand(_ *cli.Context) error {
	cacheManager, err := downloader.NewCacheManager()
	if err != nil {
		return fmt.Errorf("failed to initialize cache manager: %w", err)
	}

	versions, err := cacheManager.List()
	if err != nil {
		return fmt.Errorf("failed to list versions: %w", err)
	}

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	fmt.Printf("Cache directory: %s\n", cfg.GetVersionsDir())

	if len(versions) == 0 {
		fmt.Println("No cached versions found.")

		return nil
	}

	totalSize, err := cacheManager.GetTotalSize()
	if err != nil {
		return fmt.Errorf("failed to get total size: %w", err)
	}

	fmt.Printf("Cached versions (%d total, %.2f MB):\n", len(versions), float64(totalSize)/kilobyte/kilobyte)

	for _, version := range versions {
		status := "✓"
		extra := ""

		if !version.IsComplete {
			status = "✗"
			extra = " (incomplete)"
		}

		sizeStr := ""
		if version.Size > 0 {
			sizeStr = fmt.Sprintf(" [%.2f MB]", float64(version.Size)/kilobyte/kilobyte)
		}

		fmt.Printf("  %s %s%s%s\n", status, version.Version, sizeStr, extra)
	}

	return nil
}

// cacheCleanCommand removes old cached versions.
func cacheCleanCommand(ctx *cli.Context) error {
	cacheManager, err := downloader.NewCacheManager()
	if err != nil {
		return fmt.Errorf("failed to initialize cache manager: %w", err)
	}

	if ctx.Bool("all") {
		fmt.Println("Removing all cached versions...")

		removed, err := cacheManager.RemoveAll()
		if err != nil {
			fmt.Printf("Warning: Some versions could not be removed: %v\n", err)
		}

		fmt.Printf("\n✓ Removed %d version(s).\n", removed)

		return nil
	}

	keep := ctx.Int("keep")
	fmt.Printf("Removing old versions (keeping %d most recent)...\n", keep)

	removed, err := cacheManager.RemoveOldest(keep)
	if err != nil {
		fmt.Printf("Warning: Some versions could not be removed: %v\n", err)
	}

	if removed == 0 {
		fmt.Println("No versions to remove.")
	} else {
		fmt.Printf("\n✓ Removed %d version(s).\n", removed)
	}

	return nil
}
