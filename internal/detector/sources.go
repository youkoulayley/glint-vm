package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/youkoulayley/glint-vm/internal/config"
)

// DetectionResult represents the result of version detection.
type DetectionResult struct {
	Version    string
	Source     string // File path or source name
	SourceType string // Type of source (e.g., "github-actions", "makefile")
	LineNumber int    // Line number where version was found (0 if not applicable)
	Pattern    string // Pattern name that matched
}

// Detector is an interface for version detection sources.
type Detector interface {
	// Name returns the detector name
	Name() string
	// Detect attempts to find a golangci-lint version
	Detect(baseDir string) (*DetectionResult, error)
}

// VersionFileDetector detects version from .golangci-lint.version file.
type VersionFileDetector struct{}

// Name returns the identifier for this detector.
func (d *VersionFileDetector) Name() string {
	return "version-file"
}

// Detect searches for version in .golangci-lint.version file.
func (d *VersionFileDetector) Detect(baseDir string) (*DetectionResult, error) {
	filePath := filepath.Join(baseDir, ".golangci-lint.version")

	content, err := os.ReadFile(filePath) //nolint:gosec // Path is constructed internally
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil //nolint:nilnil // File doesn't exist, not an error
		}

		return nil, fmt.Errorf("failed to read version file: %w", err)
	}

	// For version files, try pattern matching first
	version, pattern, lineNum := ExtractVersionFromLines(string(content))

	// If no pattern matches, try to parse as plain version string
	if version == "" {
		trimmed := strings.TrimSpace(string(content))
		version = config.NormalizeVersion(trimmed)

		// Validate it's a proper version
		if !ValidateVersion(version) {
			return nil, nil //nolint:nilnil // Invalid version format, not an error
		}

		pattern = "plain-version"
		lineNum = 1
	}

	return &DetectionResult{
		Version:    version,
		Source:     filePath,
		SourceType: d.Name(),
		LineNumber: lineNum,
		Pattern:    pattern,
	}, nil
}

// GitHubActionsDetector detects version from GitHub Actions workflows.
type GitHubActionsDetector struct{}

// Name returns the identifier for this detector.
func (d *GitHubActionsDetector) Name() string {
	return "github-actions"
}

// Detect searches for version in GitHub Actions workflow files.
func (d *GitHubActionsDetector) Detect(baseDir string) (*DetectionResult, error) {
	workflowsDir := filepath.Join(baseDir, ".github", "workflows")

	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil //nolint:nilnil // Directory doesn't exist, not an error
		}

		return nil, fmt.Errorf("failed to read workflows directory: %w", err)
	}

	// Scan all YAML files in workflows directory
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if it's a YAML file
		ext := filepath.Ext(entry.Name())
		if ext != ".yml" && ext != ".yaml" {
			continue
		}

		filePath := filepath.Join(workflowsDir, entry.Name())

		content, err := os.ReadFile(filePath) //nolint:gosec // Path is constructed internally
		if err != nil {
			continue // Skip files we can't read
		}

		version, pattern, lineNum := ExtractVersionFromLines(string(content))
		if version != "" {
			return &DetectionResult{
				Version:    version,
				Source:     filePath,
				SourceType: d.Name(),
				LineNumber: lineNum,
				Pattern:    pattern,
			}, nil
		}
	}

	return nil, nil //nolint:nilnil // No version found in workflows, not an error
}

// SemaphoreDetector detects version from Semaphore CI config.
type SemaphoreDetector struct{}

// Name returns the identifier for this detector.
func (d *SemaphoreDetector) Name() string {
	return "semaphore-ci"
}

// Detect searches for version in Semaphore CI configuration.
func (d *SemaphoreDetector) Detect(baseDir string) (*DetectionResult, error) {
	filePath := filepath.Join(baseDir, ".semaphore", "semaphore.yml")

	content, err := os.ReadFile(filePath) //nolint:gosec // Path is constructed internally
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil //nolint:nilnil // File doesn't exist, not an error
		}

		return nil, fmt.Errorf("failed to read semaphore config: %w", err)
	}

	version, pattern, lineNum := ExtractVersionFromLines(string(content))
	if version == "" {
		return nil, nil //nolint:nilnil // No version found in config, not an error
	}

	return &DetectionResult{
		Version:    version,
		Source:     filePath,
		SourceType: d.Name(),
		LineNumber: lineNum,
		Pattern:    pattern,
	}, nil
}

// MakefileDetector detects version from Makefile.
type MakefileDetector struct{}

// Name returns the identifier for this detector.
func (d *MakefileDetector) Name() string {
	return "makefile"
}

// Detect searches for version in Makefile or makefile.
func (d *MakefileDetector) Detect(baseDir string) (*DetectionResult, error) {
	// Try common Makefile names
	makefileNames := []string{"Makefile", "makefile", "GNUmakefile"}

	for _, name := range makefileNames {
		filePath := filepath.Join(baseDir, name)

		content, err := os.ReadFile(filePath) //nolint:gosec // Path is constructed internally
		if err != nil {
			if os.IsNotExist(err) {
				continue // Try next filename
			}

			return nil, fmt.Errorf("failed to read makefile: %w", err)
		}

		version, pattern, lineNum := ExtractVersionFromLines(string(content))
		if version != "" {
			return &DetectionResult{
				Version:    version,
				Source:     filePath,
				SourceType: d.Name(),
				LineNumber: lineNum,
				Pattern:    pattern,
			}, nil
		}
	}

	return nil, nil //nolint:nilnil // No Makefile found, not an error
}

// CircleCIDetector detects version from CircleCI config.
type CircleCIDetector struct{}

// Name returns the identifier for this detector.
func (d *CircleCIDetector) Name() string {
	return "circleci"
}

// Detect searches for version in CircleCI configuration.
func (d *CircleCIDetector) Detect(baseDir string) (*DetectionResult, error) {
	filePath := filepath.Join(baseDir, ".circleci", "config.yml")

	content, err := os.ReadFile(filePath) //nolint:gosec // Path is constructed internally
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil //nolint:nilnil // File doesn't exist, not an error
		}

		return nil, fmt.Errorf("failed to read circleci config: %w", err)
	}

	version, pattern, lineNum := ExtractVersionFromLines(string(content))
	if version == "" {
		return nil, nil //nolint:nilnil // No version found in config, not an error
	}

	return &DetectionResult{
		Version:    version,
		Source:     filePath,
		SourceType: d.Name(),
		LineNumber: lineNum,
		Pattern:    pattern,
	}, nil
}

// GitLabCIDetector detects version from GitLab CI config.
type GitLabCIDetector struct{}

// Name returns the identifier for this detector.
func (d *GitLabCIDetector) Name() string {
	return "gitlab-ci"
}

// Detect searches for version in GitLab CI configuration.
func (d *GitLabCIDetector) Detect(baseDir string) (*DetectionResult, error) {
	filePath := filepath.Join(baseDir, ".gitlab-ci.yml")

	content, err := os.ReadFile(filePath) //nolint:gosec // Path is constructed internally
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil //nolint:nilnil // File doesn't exist, not an error
		}

		return nil, fmt.Errorf("failed to read gitlab-ci config: %w", err)
	}

	version, pattern, lineNum := ExtractVersionFromLines(string(content))
	if version == "" {
		return nil, nil //nolint:nilnil // No version found in config, not an error
	}

	return &DetectionResult{
		Version:    version,
		Source:     filePath,
		SourceType: d.Name(),
		LineNumber: lineNum,
		Pattern:    pattern,
	}, nil
}

// AllDetectors returns all available detectors in priority order.
func AllDetectors() []Detector {
	return []Detector{
		&VersionFileDetector{},   // Highest priority - explicit version file
		&GitHubActionsDetector{}, // GitHub Actions
		&SemaphoreDetector{},     // Semaphore CI
		&MakefileDetector{},      // Makefile
		&CircleCIDetector{},      // CircleCI
		&GitLabCIDetector{},      // GitLab CI
	}
}

// DetectVersion attempts to detect golangci-lint version from the given directory
// Returns the first successful detection result.
func DetectVersion(baseDir string) (*DetectionResult, error) {
	for _, detector := range AllDetectors() {
		result, err := detector.Detect(baseDir)
		if err != nil {
			continue
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, nil //nolint:nilnil // No version detected from any source, not an error
}

// DetectVersionFromAll attempts detection using all sources and returns all results.
func DetectVersionFromAll(baseDir string) ([]*DetectionResult, error) {
	var results []*DetectionResult

	for _, detector := range AllDetectors() {
		result, err := detector.Detect(baseDir)
		if err != nil {
			continue
		}

		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}
