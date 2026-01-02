package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("creates detector with explicit directory", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()

		detector, err := New(tmpDir)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		if detector.baseDir != tmpDir {
			t.Errorf("baseDir = %q, want %q", detector.baseDir, tmpDir)
		}
	})

	t.Run("creates detector with empty directory (uses cwd)", func(t *testing.T) {
		t.Parallel()

		detector, err := New("")
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		if detector.baseDir == "" {
			t.Error("baseDir should not be empty when using cwd")
		}
	})

	t.Run("returns error for non-existent directory", func(t *testing.T) {
		t.Parallel()

		_, err := New("/nonexistent/directory/that/does/not/exist")
		if err == nil {
			t.Error("New() should return error for non-existent directory")
		}
	})

	t.Run("returns error when path is a file not directory", func(t *testing.T) {
		t.Parallel()
		tmpFile := filepath.Join(t.TempDir(), "file.txt")

		err := os.WriteFile(tmpFile, []byte("test"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		_, err = New(tmpFile)
		if err == nil {
			t.Error("New() should return error when path is a file")
		}
	})
}

func TestDetect(t *testing.T) {
	t.Parallel()

	t.Run("detects version from configured directory", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()

		// Create a version file
		versionFile := filepath.Join(tmpDir, ".golangci-lint.version")

		err := os.WriteFile(versionFile, []byte("v1.55.2\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create version file: %v", err)
		}

		detector, err := New(tmpDir)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		result, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if result.Version != testVersion {
			t.Errorf("Version = %q, want %q", result.Version, testVersion)
		}
	})

	t.Run("returns errInvalidVersion when no version found", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()

		detector, err := New(tmpDir)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		result, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}

		if result != nil {
			t.Error("Detect() should return nil when no version found")
		}
	})
}

func TestDetectAll(t *testing.T) {
	t.Parallel()

	t.Run("detects from multiple sources", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()

		// Create version file
		versionFile := filepath.Join(tmpDir, ".golangci-lint.version")

		err := os.WriteFile(versionFile, []byte("v1.60.0\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create version file: %v", err)
		}

		// Create Makefile
		makefile := filepath.Join(tmpDir, "Makefile")

		err = os.WriteFile(makefile, []byte("GOLANGCI_LINT_VERSION := v1.55.2\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create Makefile: %v", err)
		}

		detector, err := New(tmpDir)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		results, err := detector.DetectAll()
		if err != nil {
			t.Fatalf("DetectAll() error = %v", err)
		}

		if len(results) < 2 {
			t.Errorf("DetectAll() returned %d results, want at least 2", len(results))
		}
	})
}

func TestDetectWithFallback(t *testing.T) {
	t.Parallel()

	t.Run("returns detected version when found", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()

		versionFile := filepath.Join(tmpDir, ".golangci-lint.version")

		err := os.WriteFile(versionFile, []byte("v1.55.2\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create version file: %v", err)
		}

		detector, err := New(tmpDir)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		result, err := detector.DetectWithFallback()
		if err != nil {
			t.Fatalf("DetectWithFallback() error = %v", err)
		}

		if result.Version != testVersion {
			t.Errorf("Version = %q, want %q", result.Version, testVersion)
		}
	})

	t.Run("returns nil when no version found (no fallback)", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()

		detector, err := New(tmpDir)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		result, err := detector.DetectWithFallback()
		if err != nil {
			t.Fatalf("DetectWithFallback() error = %v", err)
		}

		if result != nil {
			t.Error("DetectWithFallback() should return nil when nothing found")
		}
	})
}

func TestGetBaseDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	detector, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got := detector.GetBaseDir()
	if got != tmpDir {
		t.Errorf("GetBaseDir() = %q, want %q", got, tmpDir)
	}
}

func TestQuickDetect(t *testing.T) { //nolint:paralleltest // uses t.Chdir
	t.Run("detects from current directory", func(t *testing.T) {
		// Save current directory
		originalDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get current directory: %v", err)
		}

		t.Chdir(originalDir)
		tmpDir := t.TempDir()
		versionFile := filepath.Join(tmpDir, ".golangci-lint.version")

		err = os.WriteFile(versionFile, []byte("v1.55.2\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create version file: %v", err)
		}

		t.Chdir(tmpDir)

		result, err := QuickDetect()
		if err != nil {
			t.Fatalf("QuickDetect() error = %v", err)
		}

		if result.Version != testVersion {
			t.Errorf("Version = %q, want %q", result.Version, testVersion)
		}
	})
}

func TestQuickDetectFrom(t *testing.T) {
	t.Parallel()

	t.Run("detects from specific directory", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()

		versionFile := filepath.Join(tmpDir, ".golangci-lint.version")

		err := os.WriteFile(versionFile, []byte("v1.55.2\n"), 0644) //nolint:gosec // Test file
		if err != nil {
			t.Fatalf("Failed to create version file: %v", err)
		}

		result, err := QuickDetectFrom(tmpDir)
		if err != nil {
			t.Fatalf("QuickDetectFrom() error = %v", err)
		}

		if result.Version != testVersion {
			t.Errorf("Version = %q, want %q", result.Version, testVersion)
		}
	})

	t.Run("returns error for invalid directory", func(t *testing.T) {
		t.Parallel()

		_, err := QuickDetectFrom("/nonexistent/directory")
		if err == nil {
			t.Error("QuickDetectFrom() should return error for invalid directory")
		}
	})
}
