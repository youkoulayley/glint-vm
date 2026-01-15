// Package config provides configuration management for glint-vm,
// including cache directory management, platform detection, and version handling.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	// AppName is the application name used for cache directories.
	AppName = "glint-vm"
	// VersionsDir is the subdirectory name for storing versions.
	VersionsDir                     = "versions"
	directoryPermission os.FileMode = 0o700
	windows                         = "windows"
)

// Config holds the configuration for glint-vm.
type Config struct {
	// CacheDir is the base cache directory for glint-vm
	CacheDir string
	// OS is the operating system (linux, darwin, windows)
	OS string
	// Arch is the architecture (amd64, arm64, etc.)
	Arch string
}

// New creates a new Config with detected values.
func New() (*Config, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache directory: %w", err)
	}

	return &Config{
		CacheDir: cacheDir,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}, nil
}

// getCacheDir returns the cache directory following XDG Base Directory specification
// Priority:
// 1. $XDG_CACHE_HOME/glint-vm (if XDG_CACHE_HOME is set)
// 2. $HOME/.cache/glint-vm (fallback).
func getCacheDir() (string, error) {
	var baseDir string

	// Check XDG_CACHE_HOME first (XDG Base Directory compliant)
	if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
		baseDir = xdgCache
	} else {
		// Fallback to $HOME/.cache
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}

		baseDir = filepath.Join(homeDir, ".cache")
	}

	return filepath.Join(baseDir, AppName), nil
}

// GetVersionsDir returns the directory where all versions are cached.
func (c *Config) GetVersionsDir() string {
	return filepath.Join(c.CacheDir, VersionsDir)
}

// GetVersionDir returns the directory for a specific version.
func (c *Config) GetVersionDir(version string) string {
	return filepath.Join(c.GetVersionsDir(), version)
}

// GetBinaryPath returns the full path to the golangci-lint binary for a specific version.
func (c *Config) GetBinaryPath(version string) string {
	binaryName := "golangci-lint"
	if c.OS == windows {
		binaryName += ".exe"
	}

	return filepath.Join(c.GetVersionDir(version), binaryName)
}

// EnsureVersionDir creates the version directory if it doesn't exist
// Sets permissions to 0700 (user only) for security.
func (c *Config) EnsureVersionDir(version string) error {
	versionDir := c.GetVersionDir(version)

	err := os.MkdirAll(versionDir, directoryPermission)
	if err != nil {
		return fmt.Errorf("failed to create version directory: %w", err)
	}

	return nil
}

// BinaryExists checks if the binary exists for a specific version.
func (c *Config) BinaryExists(version string) bool {
	binaryPath := c.GetBinaryPath(version)

	info, err := os.Stat(binaryPath)
	if err != nil {
		return false
	}
	// Check if it's a regular file and executable
	return info.Mode().IsRegular() && isExecutable(info)
}

// isExecutable checks if a file is executable.
func isExecutable(info os.FileInfo) bool {
	// On Unix systems, check if any execute bit is set
	if runtime.GOOS != windows {
		return info.Mode().Perm()&0o111 != 0
	}
	// On Windows, any regular file can be "executed" if it has the right extension
	return true
}

// GetPlatformString returns the platform string in the format expected by golangci-lint releases
// Examples: "linux-amd64", "darwin-arm64", "windows-amd64".
func (c *Config) GetPlatformString() string {
	return fmt.Sprintf("%s-%s", c.OS, c.Arch)
}

// GetCurrentDir returns the directory containing the current version symlink.
func (c *Config) GetCurrentDir() string {
	return filepath.Join(c.CacheDir, "current")
}

// GetCurrentBinaryPath returns the path to the current golangci-lint binary symlink.
func (c *Config) GetCurrentBinaryPath() string {
	binaryName := "golangci-lint"
	if c.OS == windows {
		binaryName += ".exe"
	}

	return filepath.Join(c.GetCurrentDir(), binaryName)
}

// SetCurrentVersion manages the symlink to point to a specific version
// It creates the current directory if needed, removes old symlink, and creates new one.
func (c *Config) SetCurrentVersion(version string) error {
	// Ensure the version binary exists
	if !c.BinaryExists(version) {
		return fmt.Errorf("version %s: %w", version, ErrVersionNotInstalled)
	}

	// Create current directory if it doesn't exist
	currentDir := c.GetCurrentDir()

	err := os.MkdirAll(currentDir, directoryPermission)
	if err != nil {
		return fmt.Errorf("failed to create current directory: %w", err)
	}

	// Get paths
	currentBinaryPath := c.GetCurrentBinaryPath()
	targetBinaryPath := c.GetBinaryPath(version)

	// Remove old symlink if it exists (ignore error if it doesn't exist)
	_ = os.Remove(currentBinaryPath)

	// Create new symlink
	err = os.Symlink(targetBinaryPath, currentBinaryPath)
	if err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// GetCurrentVersion reads the symlink to determine the current version
// Returns empty string and nil error if no current version is set.
func (c *Config) GetCurrentVersion() (string, error) {
	currentBinaryPath := c.GetCurrentBinaryPath()

	// Read the symlink
	target, err := os.Readlink(currentBinaryPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No current version set
			return "", nil
		}

		return "", fmt.Errorf("failed to read current version symlink: %w", err)
	}

	// Extract version from the target path
	// Target path format: /path/to/cache/versions/v1.55.2/golangci-lint
	versionDir := filepath.Dir(target)
	version := filepath.Base(versionDir)

	return version, nil
}

// NormalizeVersion ensures version strings start with 'v'.
func NormalizeVersion(version string) string {
	if version == "" {
		return ""
	}

	if version[0] != 'v' {
		return "v" + version
	}

	return version
}
