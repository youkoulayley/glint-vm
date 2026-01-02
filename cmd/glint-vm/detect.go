package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/youkoulayley/glint-vm/internal/config"
	"github.com/youkoulayley/glint-vm/internal/detector"
	"github.com/youkoulayley/glint-vm/internal/downloader"
	"github.com/youkoulayley/glint-vm/internal/shell"
)

// detectCommand shows the detected version and source.
func detectCommand(ctx *cli.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	versionDetector, err := detector.New("")
	if err != nil {
		return fmt.Errorf("failed to create detector: %w", err)
	}

	result, err := versionDetector.Detect()
	if err != nil {
		return fmt.Errorf("detection failed: %w", err)
	}

	if result == nil {
		if ctx.Bool("quiet") {
			os.Exit(1)
		}

		fmt.Println("❌ No golangci-lint version detected in this project")
		fmt.Println("Searched in:")
		fmt.Println("  • .golangci-lint.version file")
		fmt.Println("  • GitHub Actions workflows (.github/workflows/*.yml)")
		fmt.Println("  • Semaphore CI (.semaphore/semaphore.yml)")
		fmt.Println("  • Makefile")
		fmt.Println("  • CircleCI (.circleci/config.yml)")
		fmt.Println("  • GitLab CI (.gitlab-ci.yml)")
		fmt.Println()
		fmt.Println("Create a .golangci-lint.version file with your desired version:")
		fmt.Println("  echo \"v1.55.2\" > .golangci-lint.version")

		return nil
	}

	version := result.Version

	if ctx.Bool("use") {
		dl, err := downloader.NewDownloader()
		if err != nil {
			return fmt.Errorf("failed to initialize downloader: %w", err)
		}

		if err = dl.Download(version); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		if err = cfg.SetCurrentVersion(version); err != nil {
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

	if ctx.Bool("install") {
		dl, err := downloader.NewDownloader()
		if err != nil {
			return fmt.Errorf("failed to initialize downloader: %w", err)
		}

		if err = dl.Download(version); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		fmt.Printf("✓ Installed golangci-lint %s\n", version)
		fmt.Println()
		fmt.Println("To activate this version, run:")
		fmt.Printf("  glint-vm use %s\n", version)

		return nil
	}

	if ctx.Bool("quiet") {
		fmt.Println(version)

		return nil
	}

	fmt.Printf("Cache directory: %s\n", cfg.CacheDir)
	fmt.Printf("Platform: %s\n", cfg.GetPlatformString())
	fmt.Println()

	fmt.Println("Detecting golangci-lint version...")
	fmt.Println()

	fmt.Printf("✓ Detected version: %s\n", version)
	fmt.Printf("  Source: %s\n", result.Source)
	fmt.Printf("  Source type: %s\n", result.SourceType)

	if result.LineNumber > 0 {
		fmt.Printf("  Line: %d\n", result.LineNumber)
	}

	fmt.Printf("  Pattern: %s\n", result.Pattern)
	fmt.Println()

	if cfg.BinaryExists(version) {
		fmt.Printf("✓ Binary cached at: %s\n", cfg.GetBinaryPath(version))
	} else {
		fmt.Printf("⚠ Binary not cached. Run 'glint-vm install %s' to download it.\n", version)
	}

	return nil
}
