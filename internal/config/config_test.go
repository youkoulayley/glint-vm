package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const testVersion = "v1.55.2"

func TestNew(t *testing.T) {
	t.Parallel()

	cfg, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if cfg.CacheDir == "" {
		t.Error("CacheDir should not be empty")
	}

	if cfg.OS != runtime.GOOS {
		t.Errorf("OS = %s, want %s", cfg.OS, runtime.GOOS)
	}

	if cfg.Arch != runtime.GOARCH {
		t.Errorf("Arch = %s, want %s", cfg.Arch, runtime.GOARCH)
	}
}

func TestGetCacheDir(t *testing.T) {
	tests := []struct {
		name        string
		xdgCache    string
		wantContain string
	}{
		{
			name:        "with XDG_CACHE_HOME",
			xdgCache:    "/tmp/custom-cache",
			wantContain: filepath.Join("/tmp/custom-cache", AppName),
		},
		{
			name:        "without XDG_CACHE_HOME",
			xdgCache:    "",
			wantContain: AppName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Save and restore original env
			originalXDG := os.Getenv("XDG_CACHE_HOME")

			t.Setenv("XDG_CACHE_HOME", originalXDG)

			if test.xdgCache != "" {
				t.Setenv("XDG_CACHE_HOME", test.xdgCache)
			}

			got, err := getCacheDir()
			if err != nil {
				t.Fatalf("getCacheDir() failed: %v", err)
			}

			if got != test.wantContain {
				if test.xdgCache == "" {
					// When XDG is not set, just check it contains AppName
					if !filepath.IsAbs(got) {
						t.Errorf("getCacheDir() should return absolute path, got %s", got)
					}

					if filepath.Base(got) != AppName {
						t.Errorf("getCacheDir() should end with %s, got %s", AppName, got)
					}
				} else {
					t.Errorf("getCacheDir() = %s, want %s", got, test.wantContain)
				}
			}
		})
	}
}

func TestGetVersionsDir(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		CacheDir: "/test/cache",
		OS:       "linux",
		Arch:     "amd64",
	}

	got := cfg.GetVersionsDir()
	want := filepath.Join("/test/cache", VersionsDir)

	if got != want {
		t.Errorf("GetVersionsDir() = %s, want %s", got, want)
	}
}

func TestGetVersionDir(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		CacheDir: "/test/cache",
		OS:       "linux",
		Arch:     "amd64",
	}

	version := testVersion
	got := cfg.GetVersionDir(version)
	want := filepath.Join("/test/cache", VersionsDir, version)

	if got != want {
		t.Errorf("GetVersionDir(%s) = %s, want %s", version, got, want)
	}
}

func TestGetBinaryPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		os       string
		version  string
		wantName string
	}{
		{
			name:     "linux",
			os:       "linux",
			version:  testVersion,
			wantName: "golangci-lint",
		},
		{
			name:     "darwin",
			os:       "darwin",
			version:  testVersion,
			wantName: "golangci-lint",
		},
		{
			name:     "windows",
			os:       "windows",
			version:  testVersion,
			wantName: "golangci-lint.exe",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cfg := &Config{
				CacheDir: "/test/cache",
				OS:       test.os,
				Arch:     "amd64",
			}

			got := cfg.GetBinaryPath(test.version)
			expectedPath := filepath.Join("/test/cache", VersionsDir, test.version, test.wantName)

			if got != expectedPath {
				t.Errorf("GetBinaryPath(%s) = %s, want %s", test.version, got, expectedPath)
			}
		})
	}
}

func TestGetPlatformString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		os   string
		arch string
		want string
	}{
		{
			name: "linux-amd64",
			os:   "linux",
			arch: "amd64",
			want: "linux-amd64",
		},
		{
			name: "darwin-arm64",
			os:   "darwin",
			arch: "arm64",
			want: "darwin-arm64",
		},
		{
			name: "windows-amd64",
			os:   "windows",
			arch: "amd64",
			want: "windows-amd64",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cfg := &Config{
				CacheDir: "/test/cache",
				OS:       test.os,
				Arch:     test.arch,
			}

			got := cfg.GetPlatformString()
			if got != test.want {
				t.Errorf("GetPlatformString() = %s, want %s", got, test.want)
			}
		})
	}
}

func TestNormalizeVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "with v prefix",
			input: testVersion,
			want:  testVersion,
		},
		{
			name:  "without v prefix",
			input: "1.55.2",
			want:  testVersion,
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NormalizeVersion(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeVersion(%s) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestEnsureVersionDir(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	cfg := &Config{
		CacheDir: tmpDir,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}

	version := testVersion

	err := cfg.EnsureVersionDir(version)
	if err != nil {
		t.Fatalf("EnsureVersionDir() failed: %v", err)
	}

	// Check if directory exists
	versionDir := cfg.GetVersionDir(version)

	info, err := os.Stat(versionDir)
	if err != nil {
		t.Fatalf("version directory not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("version path is not a directory")
	}

	// Check permissions (on Unix systems)
	if runtime.GOOS != windows {
		perm := info.Mode().Perm()
		if perm != 0700 {
			t.Errorf("directory permissions = %o, want 0700", perm)
		}
	}
}

func TestBinaryExists(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	cfg := &Config{
		CacheDir: tmpDir,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}

	version := testVersion

	// Initially, binary should not exist
	if cfg.BinaryExists(version) {
		t.Error("BinaryExists() = true, want false (binary doesn't exist yet)")
	}

	// Create the version directory and binary file
	err := cfg.EnsureVersionDir(version)
	if err != nil {
		t.Fatalf("EnsureVersionDir() failed: %v", err)
	}

	binaryPath := cfg.GetBinaryPath(version)
	// Create an empty file
	file, err := os.Create(binaryPath)
	if err != nil {
		t.Fatalf("failed to create binary file: %v", err)
	}

	_ = file.Close()

	// Make it executable (on Unix systems)
	if runtime.GOOS != windows {
		err = os.Chmod(binaryPath, 0755)
		if err != nil {
			t.Fatalf("failed to chmod binary: %v", err)
		}
	}

	// Now binary should exist
	if !cfg.BinaryExists(version) {
		t.Error("BinaryExists() = false, want true (binary exists)")
	}
}

func TestGetCurrentDir(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		CacheDir: "/test/cache",
		OS:       "linux",
		Arch:     "amd64",
	}

	got := cfg.GetCurrentDir()
	want := filepath.Join("/test/cache", "current")

	if got != want {
		t.Errorf("GetCurrentDir() = %s, want %s", got, want)
	}
}

func TestGetCurrentBinaryPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		os       string
		wantName string
	}{
		{
			name:     "linux",
			os:       "linux",
			wantName: "golangci-lint",
		},
		{
			name:     "darwin",
			os:       "darwin",
			wantName: "golangci-lint",
		},
		{
			name:     windows,
			os:       windows,
			wantName: "golangci-lint.exe",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cfg := &Config{
				CacheDir: "/test/cache",
				OS:       test.os,
				Arch:     "amd64",
			}

			got := cfg.GetCurrentBinaryPath()
			want := filepath.Join("/test/cache", "current", test.wantName)

			if got != want {
				t.Errorf("GetCurrentBinaryPath() = %s, want %s", got, want)
			}
		})
	}
}

// Helper functions for test setup and verification.

func setupTestBinary(t *testing.T, cfg *Config, version string) {
	t.Helper()

	if err := cfg.EnsureVersionDir(version); err != nil {
		t.Fatalf("EnsureVersionDir() failed: %v", err)
	}

	binaryPath := cfg.GetBinaryPath(version)

	file, err := os.Create(binaryPath)
	if err != nil {
		t.Fatalf("failed to create binary: %v", err)
	}

	_ = file.Close()

	if err := os.Chmod(binaryPath, 0o755); err != nil {
		t.Fatalf("failed to chmod binary: %v", err)
	}
}

func verifySymlink(t *testing.T, cfg *Config, expectedVersion string) {
	t.Helper()

	currentBinaryPath := cfg.GetCurrentBinaryPath()

	info, err := os.Lstat(currentBinaryPath)
	if err != nil {
		t.Fatalf("symlink not created: %v", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("current binary is not a symlink")
	}

	target, err := os.Readlink(currentBinaryPath)
	if err != nil {
		t.Fatalf("failed to read symlink: %v", err)
	}

	expectedTarget := cfg.GetBinaryPath(expectedVersion)
	if target != expectedTarget {
		t.Errorf("symlink target = %s, want %s", target, expectedTarget)
	}
}

func checkError(t *testing.T, err error, wantErr bool, wantErrSubstr string) bool {
	t.Helper()

	if wantErr {
		if err == nil {
			t.Error("expected error but got nil")

			return false
		}

		if wantErrSubstr != "" && !contains(err.Error(), wantErrSubstr) {
			t.Errorf("error = %q, want substring %q", err.Error(), wantErrSubstr)
		}

		return false
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)

		return false
	}

	return true
}

func TestSetCurrentVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		version       string
		setupBinary   bool
		wantErr       bool
		wantErrSubstr string
	}{
		{
			name:        "success - creates symlink",
			version:     testVersion,
			setupBinary: true,
			wantErr:     false,
		},
		{
			name:          "error - version not installed",
			version:       testVersion,
			setupBinary:   false,
			wantErr:       true,
			wantErrSubstr: "not installed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// Skip symlink tests on Windows
			if runtime.GOOS == windows {
				t.Skip("Skipping symlink test on Windows")
			}

			tmpDir := t.TempDir()
			cfg := &Config{
				CacheDir: tmpDir,
				OS:       runtime.GOOS,
				Arch:     runtime.GOARCH,
			}

			// Setup binary if needed
			if test.setupBinary {
				setupTestBinary(t, cfg, test.version)
			}

			// Test SetCurrentVersion
			err := cfg.SetCurrentVersion(test.version)

			// Check error expectations
			if !checkError(t, err, test.wantErr, test.wantErrSubstr) {
				return
			}

			// Verify symlink was created correctly
			verifySymlink(t, cfg, test.version)
		})
	}
}

func TestSetCurrentVersion_UpdatesExisting(t *testing.T) {
	t.Parallel()

	// Skip symlink tests on Windows
	if runtime.GOOS == windows {
		t.Skip("Skipping symlink test on Windows")
	}

	tmpDir := t.TempDir()
	cfg := &Config{
		CacheDir: tmpDir,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}

	// Setup two versions
	versions := []string{testVersion, "v1.54.0"}
	for _, version := range versions {
		setupTestBinary(t, cfg, version)
	}

	// Set and verify first version
	setAndVerifyVersion(t, cfg, versions[0])

	// Update and verify second version
	setAndVerifyVersion(t, cfg, versions[1])
}

func setAndVerifyVersion(t *testing.T, cfg *Config, version string) {
	t.Helper()

	if err := cfg.SetCurrentVersion(version); err != nil {
		t.Fatalf("SetCurrentVersion(%s) failed: %v", version, err)
	}

	currentVersion, err := cfg.GetCurrentVersion()
	if err != nil {
		t.Fatalf("GetCurrentVersion() failed: %v", err)
	}

	if currentVersion != version {
		t.Errorf("current version = %s, want %s", currentVersion, version)
	}
}

func TestGetCurrentVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupSymlink bool
		version      string
		wantVersion  string
		wantErr      bool
	}{
		{
			name:         "success - reads symlink",
			setupSymlink: true,
			version:      testVersion,
			wantVersion:  testVersion,
			wantErr:      false,
		},
		{
			name:         "no current version set",
			setupSymlink: false,
			wantVersion:  "",
			wantErr:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// Skip symlink tests on Windows
			if runtime.GOOS == windows {
				t.Skip("Skipping symlink test on Windows")
			}

			tmpDir := t.TempDir()
			cfg := &Config{
				CacheDir: tmpDir,
				OS:       runtime.GOOS,
				Arch:     runtime.GOARCH,
			}

			// Setup symlink if needed
			if test.setupSymlink {
				setupTestBinary(t, cfg, test.version)

				if err := cfg.SetCurrentVersion(test.version); err != nil {
					t.Fatalf("SetCurrentVersion() failed: %v", err)
				}
			}

			// Test GetCurrentVersion
			got, err := cfg.GetCurrentVersion()

			// Check error expectations
			if test.wantErr {
				if err == nil {
					t.Error("GetCurrentVersion() expected error but got nil")
				}

				return
			}

			if err != nil {
				t.Errorf("GetCurrentVersion() unexpected error: %v", err)

				return
			}

			if got != test.wantVersion {
				t.Errorf("GetCurrentVersion() = %s, want %s", got, test.wantVersion)
			}
		})
	}
}

// Helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}

	return s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
