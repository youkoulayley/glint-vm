// Package main provides the glint-vm CLI tool for managing golangci-lint versions.
package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/youkoulayley/glint-vm/internal/version"
)

const (
	limitRelease = 20
	keepVersions = 3
)

func main() {
	app := &cli.App{
		Name:                   "glint-vm",
		Usage:                  "golangci-lint version manager - like gvm/nvm for golangci-lint",
		UseShortOptionHandling: true,
		Description: `glint-vm manages golangci-lint versions for your shell environment.

   Quick Start:
     # One-time setup (add to ~/.bashrc or ~/.zshrc)
     eval "$(glint-vm init bash)"

     # Detect and activate your project's version
     glint-vm detect --use

     # Or activate a specific version
     glint-vm use v1.55.2

     # Use golangci-lint directly
     golangci-lint run

   Version detection sources (in priority order):
   1. .golangci-lint.version file
   2. GitHub Actions (.github/workflows/*.yml)
   3. Semaphore CI (.semaphore/semaphore.yml)
   4. Makefile
   5. CircleCI (.circleci/config.yml)
   6. GitLab CI (.gitlab-ci.yml)`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version.Get(), version.GetCommit(), version.GetDate()),
		Commands: []*cli.Command{
			{
				Name:      "init",
				Usage:     "Generate shell initialization code",
				ArgsUsage: "<bash|zsh>",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "auto-switch",
						Usage: "Enable automatic version switching on directory change",
					},
				},
				Action: initCommand,
			},
			{
				Name:  "detect",
				Usage: "Show detected golangci-lint version and source",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "quiet",
						Aliases: []string{"q"},
						Usage:   "Output only the version number (for scripting)",
					},
					&cli.BoolFlag{
						Name:    "install",
						Aliases: []string{"i"},
						Usage:   "Download the detected version",
					},
					&cli.BoolFlag{
						Name:    "use",
						Aliases: []string{"u"},
						Usage:   "Download and activate the detected version",
					},
				},
				Action: detectCommand,
			},
			{
				Name:      "install",
				Usage:     "Download a specific golangci-lint version",
				ArgsUsage: "<version>",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "use",
						Aliases: []string{"u"},
						Usage:   "Activate the version after installing",
					},
				},
				Action: installCommand,
			},
			{
				Name:      "use",
				Usage:     "Activate a specific version in current shell",
				ArgsUsage: "<version>",
				Action:    useCommand,
			},
			{
				Name:   "current",
				Usage:  "Show currently active version",
				Action: currentCommand,
			},
			{
				Name:   "list",
				Usage:  "List installed versions",
				Action: listCommand,
			},
			{
				Name:    "list-remote",
				Usage:   "List available versions from GitHub",
				Aliases: []string{"lr"},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "limit",
						Aliases: []string{"l"},
						Usage:   "Limit the number of versions to display",
						Value:   limitRelease,
					},
				},
				Action: listRemoteCommand,
			},
			{
				Name:      "uninstall",
				Usage:     "Remove a specific version",
				ArgsUsage: "<version>",
				Action:    uninstallCommand,
			},
			{
				Name:  "cache",
				Usage: "Manage cached golangci-lint versions",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List all cached versions",
						Action: cacheListCommand,
					},
					{
						Name:  "clean",
						Usage: "Remove old cached versions",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "all",
								Aliases: []string{"a"},
								Usage:   "Remove all cached versions",
							},
							&cli.IntFlag{
								Name:    "keep",
								Aliases: []string{"k"},
								Usage:   "Keep the N most recent versions",
								Value:   keepVersions,
							},
						},
						Action: cacheCleanCommand,
					},
				},
			},
		},
		EnableBashCompletion: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		os.Exit(1)
	}
}
