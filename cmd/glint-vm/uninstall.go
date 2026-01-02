package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/youkoulayley/glint-vm/internal/config"
	"github.com/youkoulayley/glint-vm/internal/downloader"
)

// uninstallCommand removes a specific version.
func uninstallCommand(c *cli.Context) error {
	if c.NArg() < 1 {
		return ErrVersionRequired
	}

	version := config.NormalizeVersion(c.Args().First())

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	currentVersion, _ := cfg.GetCurrentVersion()
	if version == currentVersion {
		log.Warn().Msgf("Warning: %s is currently active. You may want to switch to another version first.", version)
	}

	if !cfg.BinaryExists(version) {
		return fmt.Errorf("version %s: %w", version, config.ErrVersionNotInstalled)
	}

	cm, err := downloader.NewCacheManager()
	if err != nil {
		return fmt.Errorf("failed to initialize cache manager: %w", err)
	}

	if err := cm.Remove(version); err != nil {
		return fmt.Errorf("failed to remove version: %w", err)
	}

	log.Info().Msgf("✓ Removed version %s", version)

	return nil
}
