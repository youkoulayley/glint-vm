// Package detector provides functionality to detect golangci-lint versions
// from various configuration sources like version files, CI configs, and Makefiles.
package detector

import (
	"fmt"
	"os"
)

// VersionDetector is the main orchestrator for detecting golangci-lint versions.
type VersionDetector struct {
	baseDir string
}

// New creates a new VersionDetector for the given directory
// If baseDir is empty, uses current working directory.
func New(baseDir string) (*VersionDetector, error) {
	if baseDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}

		baseDir = cwd
	}

	// Verify directory exists
	info, err := os.Stat(baseDir)
	if err != nil {
		return nil, fmt.Errorf("directory does not exist: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s: %w", baseDir, ErrNotDirectory)
	}

	return &VersionDetector{
		baseDir: baseDir,
	}, nil
}

// Detect attempts to detect golangci-lint version from the configured directory
// Returns the first successful detection result based on priority order.
func (d *VersionDetector) Detect() (*DetectionResult, error) {
	result, err := DetectVersion(d.baseDir)
	if err != nil {
		return nil, fmt.Errorf("detection failed: %w", err)
	}

	return result, nil
}

// DetectAll attempts detection using all sources and returns all results.
func (d *VersionDetector) DetectAll() ([]*DetectionResult, error) {
	results, err := DetectVersionFromAll(d.baseDir)
	if err != nil {
		return nil, fmt.Errorf("detection failed: %w", err)
	}

	return results, nil
}

// DetectWithFallback attempts to detect version with a fallback strategy
// If no version is found in config files, returns nil (no fallback to latest).
func (d *VersionDetector) DetectWithFallback() (*DetectionResult, error) {
	result, err := d.Detect()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetBaseDir returns the base directory being used for detection.
func (d *VersionDetector) GetBaseDir() string {
	return d.baseDir
}

// QuickDetect is a convenience function for detecting version from current directory.
func QuickDetect() (*DetectionResult, error) {
	detector, err := New("")
	if err != nil {
		return nil, err
	}

	return detector.Detect()
}

// QuickDetectFrom is a convenience function for detecting version from a specific directory.
func QuickDetectFrom(dir string) (*DetectionResult, error) {
	detector, err := New(dir)
	if err != nil {
		return nil, err
	}

	return detector.Detect()
}
