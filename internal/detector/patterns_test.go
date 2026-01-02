package detector

import (
	"testing"
)

func TestAtVersionPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "with v prefix",
			input: "golangci-lint@v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "without v prefix",
			input: "golangci-lint@1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "in longer text",
			input: "run golangci-lint@v1.54.1 for linting",
			want:  "v1.54.1",
		},
		{
			name:  "no match",
			input: "golangci-lint without version",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetAtVersionPattern().ExtractVersion(tt.input)
			if got != tt.want {
				t.Errorf("ExtractVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEnvVersionPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "equals sign",
			input: "GOLANGCI_LINT_VERSION=v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "colon (YAML style)",
			input: "GOLANGCI_LINT_VERSION: v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "with spaces",
			input: "GOLANGCI_LINT_VERSION = v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "without v prefix",
			input: "GOLANGCI_LINT_VERSION=1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "in script",
			input: "export GOLANGCI_LINT_VERSION=v1.54.0",
			want:  "v1.54.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetEnvVersionPattern().ExtractVersion(tt.input)
			if got != tt.want {
				t.Errorf("ExtractVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDockerImagePattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "docker image",
			input: "golangci/golangci-lint:v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "without v prefix",
			input: "golangci/golangci-lint:1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "in docker run command",
			input: "docker run golangci/golangci-lint:v1.54.1",
			want:  "v1.54.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetDockerImagePattern().ExtractVersion(tt.input)
			if got != tt.want {
				t.Errorf("ExtractVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMakefileAssignPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple assignment",
			input: "GOLANGCI_LINT_VERSION = v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "colon equals (Makefile style)",
			input: "GOLANGCI_LINT_VERSION := v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "without spaces",
			input: "GOLANGCI_LINT_VERSION=v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "without v prefix",
			input: "GOLANGCI_LINT_VERSION := 1.55.2",
			want:  "v1.55.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetMakefileAssignPattern().ExtractVersion(tt.input)
			if got != tt.want {
				t.Errorf("ExtractVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestYAMLVersionPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "yaml version",
			input: "version: v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "with extra spaces",
			input: "version:     v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "multiple spaces",
			input: "version:   v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "without v prefix",
			input: "version: 1.55.2",
			want:  "v1.55.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetYAMLVersionPattern().ExtractVersion(tt.input)
			if got != tt.want {
				t.Errorf("ExtractVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInstallVersionPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "install version",
			input: "install-version: v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "with spaces",
			input: "install-version:  v1.55.2",
			want:  "v1.55.2",
		},
		{
			name:  "without v prefix",
			input: "install-version: 1.55.2",
			want:  "v1.55.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetInstallVersionPattern().ExtractVersion(tt.input)
			if got != tt.want {
				t.Errorf("ExtractVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFindVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		input           string
		wantVersion     string
		wantPatternName string
	}{
		{
			name:            "env version (highest priority)",
			input:           "GOLANGCI_LINT_VERSION=v1.55.2",
			wantVersion:     "v1.55.2",
			wantPatternName: "env-version",
		},
		{
			name:            "docker image",
			input:           "golangci/golangci-lint:v1.54.1",
			wantVersion:     "v1.54.1",
			wantPatternName: "docker-image",
		},
		{
			name:            "at version",
			input:           "golangci-lint@v1.53.0",
			wantVersion:     "v1.53.0",
			wantPatternName: "at-version",
		},
		{
			name:            "no match",
			input:           "some random text",
			wantVersion:     "",
			wantPatternName: "",
		},
		{
			name:            "multiple patterns - prefers env",
			input:           "GOLANGCI_LINT_VERSION=v1.55.2 golangci-lint@v1.54.0",
			wantVersion:     "v1.55.2",
			wantPatternName: "env-version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotVersion, gotPattern := FindVersion(test.input)
			if gotVersion != test.wantVersion {
				t.Errorf("FindVersion() version = %q, want %q", gotVersion, test.wantVersion)
			}

			if gotPattern != test.wantPatternName {
				t.Errorf("FindVersion() pattern = %q, want %q", gotPattern, test.wantPatternName)
			}
		})
	}
}

func TestFindAllVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantCount int
		wantFirst string
	}{
		{
			name:      "single version",
			input:     "GOLANGCI_LINT_VERSION=v1.55.2",
			wantCount: 1,
			wantFirst: "v1.55.2",
		},
		{
			name:      "multiple different versions",
			input:     "GOLANGCI_LINT_VERSION=v1.55.2 golangci-lint@v1.54.0",
			wantCount: 2,
			wantFirst: "v1.55.2",
		},
		{
			name:      "duplicate versions",
			input:     "GOLANGCI_LINT_VERSION=v1.55.2 version: v1.55.2",
			wantCount: 1,
			wantFirst: "v1.55.2",
		},
		{
			name:      "no versions",
			input:     "some random text",
			wantCount: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			matches := FindAllVersions(test.input)
			if len(matches) != test.wantCount {
				t.Errorf("FindAllVersions() returned %d matches, want %d", len(matches), test.wantCount)
			}

			if test.wantCount > 0 && matches[0].Version != test.wantFirst {
				t.Errorf("FindAllVersions() first version = %q, want %q", matches[0].Version, test.wantFirst)
			}
		})
	}
}

func TestExtractVersionFromLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          string
		wantVersion    string
		wantPattern    string
		wantLineNumber int
	}{
		{
			name: "version on first line",
			input: `GOLANGCI_LINT_VERSION=v1.55.2
other content`,
			wantVersion:    "v1.55.2",
			wantPattern:    "env-version",
			wantLineNumber: 1,
		},
		{
			name: "version on third line",
			input: `first line
second line
version: v1.54.0
fourth line`,
			wantVersion:    "v1.54.0",
			wantPattern:    "yaml-version",
			wantLineNumber: 3,
		},
		{
			name: "no version",
			input: `first line
second line
third line`,
			wantVersion:    "",
			wantPattern:    "",
			wantLineNumber: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotVersion, gotPattern, gotLine := ExtractVersionFromLines(test.input)
			if gotVersion != test.wantVersion {
				t.Errorf("ExtractVersionFromLines() version = %q, want %q", gotVersion, test.wantVersion)
			}

			if gotPattern != test.wantPattern {
				t.Errorf("ExtractVersionFromLines() pattern = %q, want %q", gotPattern, test.wantPattern)
			}

			if gotLine != test.wantLineNumber {
				t.Errorf("ExtractVersionFromLines() line = %d, want %d", gotLine, test.wantLineNumber)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{
			name:    "valid version with v",
			version: "v1.55.2",
			want:    true,
		},
		{
			name:    "valid version different",
			version: "v2.0.1",
			want:    true,
		},
		{
			name:    "invalid - no v prefix",
			version: "1.55.2",
			want:    false,
		},
		{
			name:    "invalid - empty",
			version: "",
			want:    false,
		},
		{
			name:    "invalid - wrong format",
			version: "v1.55",
			want:    false,
		},
		{
			name:    "invalid - extra text",
			version: "v1.55.2-alpha",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ValidateVersion(tt.version)
			if got != tt.want {
				t.Errorf("ValidateVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestRealWorldExamples(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		wantVersion string
	}{
		{
			name: "GitHub Actions workflow",
			input: `      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: v1.55.2`,
			wantVersion: "v1.55.2",
		},
		{
			name: "Makefile",
			input: `GOLANGCI_LINT_VERSION := v1.54.2

lint:
	golangci-lint run`,
			wantVersion: "v1.54.2",
		},
		{
			name: "Docker in CI",
			input: `lint:
  docker:
    - image: golangci/golangci-lint:v1.53.3`,
			wantVersion: "v1.53.3",
		},
		{
			name: "Semaphore CI",
			input: `  - name: Lint
    commands:
      - curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | ` +
				`sh -s -- -b $(go env GOPATH)/bin v1.52.0`,
			wantVersion: "v1.52.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			version, _ := FindVersion(tt.input)
			if version != tt.wantVersion {
				t.Errorf("FindVersion() = %q, want %q", version, tt.wantVersion)
			}
		})
	}
}
