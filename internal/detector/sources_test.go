package detector

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	testVersion = "v1.55.2"
	versionFile = "version-file"
)

func TestVersionFileDetector(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	t.Run("version file exists with valid version", func(t *testing.T) {
		t.Parallel()
		// Create version file
		versionFilePath := filepath.Join(tmpDir, ".golangci-lint.version")

		err := os.WriteFile(versionFilePath, []byte("v1.55.2\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		detector := &VersionFileDetector{}

		result, err := detector.Detect(tmpDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if result == nil {
			t.Fatal("Detect() returned nil result")
		}

		if result.Version != testVersion {
			t.Errorf("Version = %q, want %q", result.Version, testVersion)
		}

		if result.SourceType != versionFile {
			t.Errorf("SourceType = %q, want %q", result.SourceType, versionFile)
		}
	})

	t.Run("version file does not exist", func(t *testing.T) {
		t.Parallel()
		emptyDir := t.TempDir()
		detector := &VersionFileDetector{}

		result, err := detector.Detect(emptyDir)
		if err != nil {
			t.Errorf("Detect() should not error when file doesn't exist: %v", err)
		}

		if result != nil {
			t.Errorf("Detect() should return nil when file doesn't exist")
		}
	})
}

func TestGitHubActionsDetector(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	t.Run("detects version from workflow file", func(t *testing.T) {
		t.Parallel()
		// Create .github/workflows directory
		workflowDir := filepath.Join(tmpDir, ".github", "workflows")

		err := os.MkdirAll(workflowDir, 0755) //nolint:gosec // Test directory
		if err != nil {
			t.Fatalf("Failed to create workflows dir: %v", err)
		}

		// Create workflow file with golangci-lint version
		workflowContent := `name: Lint
on: [push]
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: v1.55.2
`
		workflowFile := filepath.Join(workflowDir, "lint.yml")

		err = os.WriteFile(workflowFile, []byte(workflowContent), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create workflow file: %v", err)
		}

		detector := &GitHubActionsDetector{}

		result, err := detector.Detect(tmpDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if result == nil {
			t.Fatal("Detect() returned nil result")
		}

		if result.Version != testVersion {
			t.Errorf("Version = %q, want %q", result.Version, testVersion)
		}

		if result.SourceType != "github-actions" {
			t.Errorf("SourceType = %q, want %q", result.SourceType, "github-actions")
		}

		if result.LineNumber == 0 {
			t.Error("LineNumber should not be 0")
		}
	})

	t.Run("workflows directory does not exist", func(t *testing.T) {
		t.Parallel()
		emptyDir := t.TempDir()
		detector := &GitHubActionsDetector{}

		result, err := detector.Detect(emptyDir)
		if err != nil {
			t.Errorf("Detect() should not error when directory doesn't exist: %v", err)
		}

		if result != nil {
			t.Errorf("Detect() should return nil when directory doesn't exist")
		}
	})

	t.Run("multiple workflow files - returns first match", func(t *testing.T) {
		t.Parallel()

		workflowDir := filepath.Join(tmpDir, ".github", "workflows")
		_ = os.RemoveAll(workflowDir)

		err := os.MkdirAll(workflowDir, 0755) //nolint:gosec // Test directory
		if err != nil {
			t.Fatalf("Failed to create workflows dir: %v", err)
		}

		// Create first workflow with version
		workflow1 := filepath.Join(workflowDir, "ci.yml")

		err = os.WriteFile(workflow1, []byte("version: v1.54.0\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create workflow: %v", err)
		}

		// Create second workflow
		workflow2 := filepath.Join(workflowDir, "lint.yml")

		err = os.WriteFile(workflow2, []byte("version: v1.55.0\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create workflow: %v", err)
		}

		detector := &GitHubActionsDetector{}

		result, err := detector.Detect(tmpDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if result == nil {
			t.Fatal("Detect() returned nil result")
		}
		// Should return a version (order depends on directory listing)
		if result.Version == "" {
			t.Error("Version should not be empty")
		}
	})
}

func TestMakefileDetector(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	t.Run("detects version from Makefile", func(t *testing.T) {
		t.Parallel()

		makefileContent := `GOLANGCI_LINT_VERSION := v1.55.2

.PHONY: lint
lint:
	golangci-lint run
`
		makefilePath := filepath.Join(tmpDir, "Makefile")

		err := os.WriteFile(makefilePath, []byte(makefileContent), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create Makefile: %v", err)
		}

		detector := &MakefileDetector{}

		result, err := detector.Detect(tmpDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if result == nil {
			t.Fatal("Detect() returned nil result")
		}

		if result.Version != testVersion {
			t.Errorf("Version = %q, want %q", result.Version, testVersion)
		}

		if result.SourceType != "makefile" {
			t.Errorf("SourceType = %q, want %q", result.SourceType, "makefile")
		}
	})

	t.Run("tries lowercase makefile", func(t *testing.T) {
		t.Parallel()
		tmpDir2 := t.TempDir()
		makefileContent := `GOLANGCI_LINT_VERSION = v1.54.0
`
		makefilePath := filepath.Join(tmpDir2, "makefile")

		err := os.WriteFile(makefilePath, []byte(makefileContent), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create makefile: %v", err)
		}

		detector := &MakefileDetector{}

		result, err := detector.Detect(tmpDir2)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if result == nil {
			t.Fatal("Detect() returned nil result")
		}

		if result.Version != "v1.54.0" {
			t.Errorf("Version = %q, want %q", result.Version, "v1.54.0")
		}
	})

	t.Run("no Makefile exists", func(t *testing.T) {
		t.Parallel()
		emptyDir := t.TempDir()
		detector := &MakefileDetector{}

		result, err := detector.Detect(emptyDir)
		if err != nil {
			t.Errorf("Detect() should not error when file doesn't exist: %v", err)
		}

		if result != nil {
			t.Errorf("Detect() should return nil when file doesn't exist")
		}
	})
}

func TestSemaphoreDetector(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	t.Run("detects version from semaphore config", func(t *testing.T) {
		t.Parallel()

		semaphoreDir := filepath.Join(tmpDir, ".semaphore")

		err := os.MkdirAll(semaphoreDir, 0755) //nolint:gosec // Test directory
		if err != nil {
			t.Fatalf("Failed to create .semaphore dir: %v", err)
		}

		semaphoreContent := `version: v1.0
name: Lint
agent:
  machine:
    type: e1-standard-2
blocks:
  - name: Lint
    task:
      jobs:
        - name: golangci-lint
          commands:
            - checkout
            - curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | ` +
			`sh -s -- -b $(go env GOPATH)/bin v1.55.2
`
		configPath := filepath.Join(semaphoreDir, "semaphore.yml")

		err = os.WriteFile(configPath, []byte(semaphoreContent), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create semaphore config: %v", err)
		}

		detector := &SemaphoreDetector{}

		result, err := detector.Detect(tmpDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if result == nil {
			t.Fatal("Detect() returned nil result")
		}

		if result.Version != testVersion {
			t.Errorf("Version = %q, want %q", result.Version, testVersion)
		}

		if result.SourceType != "semaphore-ci" {
			t.Errorf("SourceType = %q, want %q", result.SourceType, "semaphore-ci")
		}
	})

	t.Run("semaphore config does not exist", func(t *testing.T) {
		t.Parallel()
		emptyDir := t.TempDir()
		detector := &SemaphoreDetector{}

		result, err := detector.Detect(emptyDir)
		if err != nil {
			t.Errorf("Detect() should not error when file doesn't exist: %v", err)
		}

		if result != nil {
			t.Errorf("Detect() should return nil when file doesn't exist")
		}
	})
}

func TestCircleCIDetector(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	t.Run("detects version from circleci config", func(t *testing.T) {
		t.Parallel()

		circleCIDir := filepath.Join(tmpDir, ".circleci")

		err := os.MkdirAll(circleCIDir, 0755) //nolint:gosec // Test directory
		if err != nil {
			t.Fatalf("Failed to create .circleci dir: %v", err)
		}

		circleContent := `version: 2.1
jobs:
  lint:
    docker:
      - image: golangci/golangci-lint:v1.55.2
    steps:
      - checkout
      - run: golangci-lint run
`
		configPath := filepath.Join(circleCIDir, "config.yml")

		err = os.WriteFile(configPath, []byte(circleContent), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create circleci config: %v", err)
		}

		detector := &CircleCIDetector{}

		result, err := detector.Detect(tmpDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if result == nil {
			t.Fatal("Detect() returned nil result")
		}

		if result.Version != testVersion {
			t.Errorf("Version = %q, want %q", result.Version, testVersion)
		}

		if result.SourceType != "circleci" {
			t.Errorf("SourceType = %q, want %q", result.SourceType, "circleci")
		}
	})
}

func TestGitLabCIDetector(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	t.Run("detects version from gitlab-ci config", func(t *testing.T) {
		t.Parallel()

		gitlabContent := `stages:
  - lint

lint:
  image: golangci/golangci-lint:v1.55.2
  script:
    - golangci-lint run
`
		configPath := filepath.Join(tmpDir, ".gitlab-ci.yml")

		err := os.WriteFile(configPath, []byte(gitlabContent), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create gitlab-ci config: %v", err)
		}

		detector := &GitLabCIDetector{}

		result, err := detector.Detect(tmpDir)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if result == nil {
			t.Fatal("Detect() returned nil result")
		}

		if result.Version != testVersion {
			t.Errorf("Version = %q, want %q", result.Version, testVersion)
		}

		if result.SourceType != "gitlab-ci" {
			t.Errorf("SourceType = %q, want %q", result.SourceType, "gitlab-ci")
		}
	})
}

func TestDetectVersion(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	t.Run("detects from highest priority source", func(t *testing.T) {
		t.Parallel()
		// Create version file (highest priority)
		versionFile := filepath.Join(tmpDir, ".golangci-lint.version")

		err := os.WriteFile(versionFile, []byte("v1.60.0\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create version file: %v", err)
		}

		// Also create a Makefile with different version (lower priority)
		makefilePath := filepath.Join(tmpDir, "Makefile")

		err = os.WriteFile(makefilePath, []byte("GOLANGCI_LINT_VERSION := v1.55.2\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create Makefile: %v", err)
		}

		result, err := DetectVersion(tmpDir)
		if err != nil {
			t.Fatalf("DetectVersion() error = %v", err)
		}

		if result == nil {
			t.Fatal("DetectVersion() returned nil")
		}
		// Should prefer version file
		if result.Version != "v1.60.0" {
			t.Errorf("Version = %q, want %q (should use highest priority)", result.Version, "v1.60.0")
		}

		if result.SourceType != versionFile {
			t.Errorf("SourceType = %q, want %q", result.SourceType, versionFile)
		}
	})

	t.Run("returns nil when no version found", func(t *testing.T) {
		t.Parallel()
		emptyDir := t.TempDir()

		result, err := DetectVersion(emptyDir)
		if err != nil {
			t.Errorf("DetectVersion() should not error when nothing found: %v", err)
		}

		if result != nil {
			t.Errorf("DetectVersion() should return nil when nothing found")
		}
	})
}

func TestDetectVersionFromAll(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	t.Run("detects from multiple sources", func(t *testing.T) {
		t.Parallel()
		// Create version file
		versionFile := filepath.Join(tmpDir, ".golangci-lint.version")

		err := os.WriteFile(versionFile, []byte("v1.60.0\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create version file: %v", err)
		}

		// Create Makefile
		makefilePath := filepath.Join(tmpDir, "Makefile")

		err = os.WriteFile(makefilePath, []byte("GOLANGCI_LINT_VERSION := v1.55.2\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create Makefile: %v", err)
		}

		results, err := DetectVersionFromAll(tmpDir)
		if err != nil {
			t.Fatalf("DetectVersionFromAll() error = %v", err)
		}

		if len(results) < 2 {
			t.Errorf("DetectVersionFromAll() returned %d results, want at least 2", len(results))
		}
	})

	t.Run("returns empty slice when nothing found", func(t *testing.T) {
		t.Parallel()
		emptyDir := t.TempDir()

		results, err := DetectVersionFromAll(emptyDir)
		if err != nil {
			t.Errorf("DetectVersionFromAll() should not error: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("DetectVersionFromAll() returned %d results, want 0", len(results))
		}
	})
}
