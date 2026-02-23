package detector

import (
	"regexp"
	"strings"

	"github.com/youkoulayley/glint-vm/internal/config"
)

// VersionPattern represents a compiled regex pattern for matching golangci-lint versions.
type VersionPattern struct {
	// Name is a descriptive name for this pattern.
	Name string
	// Regex is the compiled regular expression.
	Regex *regexp.Regexp
	// GroupIndex is the capture group index that contains the version (default: 1).
	GroupIndex int
}

// patternSet holds all version patterns.
type patternSet struct {
	atVersion          *VersionPattern
	envVersion         *VersionPattern
	yamlVersion        *VersionPattern
	dockerImage        *VersionPattern
	actionVersion      *VersionPattern
	makefileAssign     *VersionPattern
	cliVersion         *VersionPattern
	altEnv             *VersionPattern
	filename           *VersionPattern
	installVersion     *VersionPattern
	shellScriptVersion *VersionPattern
}

func newPattern(name, pattern string) *VersionPattern {
	return &VersionPattern{
		Name:       name,
		Regex:      regexp.MustCompile(pattern),
		GroupIndex: 1,
	}
}

// createPatterns creates and returns all version patterns.
func createPatterns() *patternSet {
	return &patternSet{
		atVersion:   newPattern("at-version", `golangci-lint@(v?\d+\.\d+\.\d+)`),
		envVersion:  newPattern("env-version", `GOLANGCI_LINT_VERSION[=:\s]+['"]?(v?\d+\.\d+\.\d+)['"]?`),
		yamlVersion: newPattern("yaml-version", `version:\s*['"]?(v?\d+\.\d+\.\d+)['"]?`),
		dockerImage: newPattern("docker-image", `golangci/golangci-lint:(v?\d+\.\d+\.\d+)`),
		actionVersion: newPattern(
			"action-version",
			`golangci-lint-action@v\d+.*?version:\s*['"]?(v?\d+\.\d+\.\d+)['"]?`,
		),
		makefileAssign:     newPattern("makefile-assign", `GOLANGCI_LINT_VERSION\s*:?=\s*['"]?(v?\d+\.\d+\.\d+)['"]?`),
		cliVersion:         newPattern("cli-version", `golangci-lint.*?--version\s+(v?\d+\.\d+\.\d+)`),
		altEnv:             newPattern("alt-env", `GOLANGCI_VERSION[=:\s]+['"]?(v?\d+\.\d+\.\d+)['"]?`),
		filename:           newPattern("filename", `golangci-lint-v?(\d+\.\d+\.\d+)`),
		installVersion:     newPattern("install-version", `install-version:\s*['"]?(v?\d+\.\d+\.\d+)['"]?`),
		shellScriptVersion: newPattern("shell-script-version", `install\.sh.*?\s+(v?\d+\.\d+\.\d+)\s*$`),
	}
}

// AllPatterns returns all available version patterns in priority order.
func AllPatterns() []*VersionPattern {
	p := createPatterns()

	return []*VersionPattern{
		p.envVersion,         // Highest priority - explicit env vars
		p.makefileAssign,     // Makefile assignments
		p.installVersion,     // GitHub Actions install-version
		p.actionVersion,      // GitHub Actions with version
		p.dockerImage,        // Docker images
		p.atVersion,          // @version syntax
		p.yamlVersion,        // Generic YAML version
		p.altEnv,             // Alternative env vars
		p.cliVersion,         // CLI flags
		p.shellScriptVersion, // Shell script install commands
		p.filename,           // Filenames/URLs (lowest priority)
	}
}

// Pattern accessor functions for backward compatibility with tests.

// GetAtVersionPattern returns the pattern for matching @version syntax.
func GetAtVersionPattern() *VersionPattern {
	return createPatterns().atVersion
}

// GetEnvVersionPattern returns the pattern for matching env variables.
func GetEnvVersionPattern() *VersionPattern {
	return createPatterns().envVersion
}

// GetDockerImagePattern returns the pattern for matching Docker images.
func GetDockerImagePattern() *VersionPattern {
	return createPatterns().dockerImage
}

// GetMakefileAssignPattern returns the pattern for matching Makefile assignments.
func GetMakefileAssignPattern() *VersionPattern {
	return createPatterns().makefileAssign
}

// GetYAMLVersionPattern returns the pattern for matching YAML versions.
func GetYAMLVersionPattern() *VersionPattern {
	return createPatterns().yamlVersion
}

// GetInstallVersionPattern returns the pattern for matching install-version.
func GetInstallVersionPattern() *VersionPattern {
	return createPatterns().installVersion
}

// ExtractVersion tries to extract a version string from the given text using the pattern.
func (p *VersionPattern) ExtractVersion(text string) string {
	matches := p.Regex.FindStringSubmatch(text)
	if len(matches) > p.GroupIndex {
		version := strings.TrimSpace(matches[p.GroupIndex])

		return config.NormalizeVersion(version)
	}

	return ""
}

// FindVersion searches through text using all patterns and returns the first match.
func FindVersion(text string) (string, string) {
	for _, pattern := range AllPatterns() {
		if version := pattern.ExtractVersion(text); version != "" {
			return version, pattern.Name
		}
	}

	return "", ""
}

// FindAllVersions searches through text and returns all unique versions found.
func FindAllVersions(text string) []VersionMatch {
	var matches []VersionMatch

	seen := make(map[string]bool)

	for _, pattern := range AllPatterns() {
		if version := pattern.ExtractVersion(text); version != "" {
			if !seen[version] {
				matches = append(matches, VersionMatch{
					Version:     version,
					PatternName: pattern.Name,
				})
				seen[version] = true
			}
		}
	}

	return matches
}

// VersionMatch represents a found version with metadata.
type VersionMatch struct {
	Version     string
	PatternName string
}

// ExtractVersionFromLines processes text line by line and returns the first version found
// This is useful for large files where you want to stop at the first match.
func ExtractVersionFromLines(text string) (string, string, int) {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if version, patternName := FindVersion(line); version != "" {
			return version, patternName, i + 1
		}
	}

	return "", "", 0
}

// ValidateVersion checks if a version string looks valid.
func ValidateVersion(version string) bool {
	if version == "" {
		return false
	}
	// Must start with 'v' and have semver format
	versionRegex := regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

	return versionRegex.MatchString(version)
}
