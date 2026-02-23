// Package downloader provides functionality for downloading and caching golangci-lint binaries.
package downloader

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/youkoulayley/glint-vm/internal/config"
)

const (
	gitHubReleasesURL                = "https://github.com/golangci/golangci-lint/releases/download"
	downloadTimeout                  = 10 * time.Minute
	maxExtractSize                   = 500 * 1024 * 1024
	executablePermission os.FileMode = 0o755
)

// Downloader handles downloading golangci-lint binaries.
type Downloader struct {
	config       *config.Config
	cacheManager *CacheManager
	httpClient   *http.Client
}

// NewDownloader creates a new downloader.
func NewDownloader() (*Downloader, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	cacheManager, err := NewCacheManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache manager: %w", err)
	}

	return &Downloader{
		config:       cfg,
		cacheManager: cacheManager,
		httpClient: &http.Client{
			Timeout: downloadTimeout,
		},
	}, nil
}

// Download downloads and installs a specific version of golangci-lint.
func (d *Downloader) Download(version string) error {
	version = config.NormalizeVersion(version)

	// Check if already cached
	if d.cacheManager.IsCached(version) {
		return nil // Already downloaded
	}

	// Construct download URL
	archiveURL := d.getDownloadURL(version)
	checksumURL := archiveURL + ".sha256"

	fmt.Fprintf(os.Stderr, "Downloading golangci-lint %s...\n", version)
	fmt.Fprintf(os.Stderr, "URL: %s\n", archiveURL)

	// Create version directory
	err := d.cacheManager.EnsureVersionDir(version)
	if err != nil {
		return fmt.Errorf("failed to create version directory: %w", err)
	}

	versionDir := d.cacheManager.GetVersionDir(version)
	archivePath := filepath.Join(versionDir, "archive.tar.gz")

	// Download archive
	err = d.downloadFile(archiveURL, archivePath)
	if err != nil {
		return fmt.Errorf("failed to download archive: %w", err)
	}

	// Download and verify checksum
	err = d.verifyChecksum(archivePath, checksumURL)
	if err != nil {
		// Checksum verification failed, clean up
		_ = os.RemoveAll(versionDir)

		return fmt.Errorf("checksum verification failed: %w", err)
	}

	// Extract archive
	err = d.extractArchive(archivePath, versionDir)
	if err != nil {
		_ = os.RemoveAll(versionDir)

		return fmt.Errorf("failed to extract archive: %w", err)
	}

	// Remove archive file
	_ = os.Remove(archivePath)

	// Verify binary exists and is executable
	if !d.cacheManager.IsCached(version) {
		_ = os.RemoveAll(versionDir)

		return ErrBinaryNotFound
	}

	fmt.Fprintf(os.Stderr, "✓ Successfully installed golangci-lint %s\n", version)
	fmt.Fprintf(os.Stderr, "  Location: %s\n", d.cacheManager.GetBinaryPath(version))

	return nil
}

// getDownloadURL constructs the download URL for a version.
func (d *Downloader) getDownloadURL(version string) string {
	platform := d.config.GetPlatformString()
	// Example: https://github.com/golangci/golangci-lint/releases/download/v1.55.2/golangci-lint-1.55.2-linux-amd64.tar.gz
	versionWithoutV := strings.TrimPrefix(version, "v")
	filename := fmt.Sprintf("golangci-lint-%s-%s.tar.gz", versionWithoutV, platform)

	return fmt.Sprintf("%s/%s/%s", gitHubReleasesURL, version, filename)
}

// downloadFile downloads a file from URL to destination.
func (d *Downloader) downloadFile(url, dest string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	//nolint:gosec // URL is constructed from hardcoded GitHub releases URL and validated version
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: HTTP %d: %s", ErrHTTPRequest, resp.StatusCode, resp.Status)
	}

	out, err := os.Create(dest) //nolint:gosec // Path is internally controlled
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() { _ = out.Close() }()

	// Copy with progress (simple version without progress bar for now)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// verifyChecksum downloads the checksum file and verifies the archive.
func (d *Downloader) verifyChecksum(archivePath, checksumURL string) error {
	// Download checksum file
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, checksumURL, http.NoBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: Could not create checksum request, skipping verification")

		return nil
	}

	//nolint:gosec // URL is constructed from hardcoded GitHub releases URL and validated version
	resp, err := d.httpClient.Do(req)
	if err != nil {
		// Checksum file might not exist for all versions, skip verification
		fmt.Fprintln(os.Stderr, "Warning: Could not download checksum file, skipping verification")

		return nil
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		// Checksum file doesn't exist, skip verification
		fmt.Fprintln(os.Stderr, "Warning: Checksum file not available, skipping verification")

		return nil
	}

	// Read expected checksum
	expectedChecksum, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read checksum: %w", err)
	}

	expectedHash := strings.TrimSpace(strings.Fields(string(expectedChecksum))[0])

	// Calculate actual checksum
	file, err := os.Open(archivePath) //nolint:gosec // Path is internally controlled
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}

	defer func() { _ = file.Close() }()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	actualHash := hex.EncodeToString(hash.Sum(nil))

	if actualHash != expectedHash {
		return fmt.Errorf("%w: expected %s, got %s", ErrChecksumMismatch, expectedHash, actualHash)
	}

	fmt.Fprintln(os.Stderr, "✓ Checksum verified")

	return nil
}

// extractArchive extracts a tar.gz archive to destination directory.
func (d *Downloader) extractArchive(archivePath, destDir string) error {
	fmt.Fprintln(os.Stderr, "Extracting archive...")

	file, err := os.Open(archivePath) //nolint:gosec // Path is internally controlled
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}

	defer func() { _ = file.Close() }()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}

	defer func() { _ = gzr.Close() }()

	reader := tar.NewReader(gzr)

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Only extract the golangci-lint binary
		if !strings.HasSuffix(header.Name, "golangci-lint") && !strings.HasSuffix(header.Name, "golangci-lint.exe") {
			continue
		}

		if header.Typeflag == tar.TypeReg {
			target := filepath.Join(destDir, filepath.Base(header.Name))

			if err := d.extractFile(reader, target, header.Mode); err != nil {
				return err
			}

			//nolint:gosec // Writing to stderr, not a web context - XSS not applicable
			fmt.Fprintf(os.Stderr, "✓ Extracted binary to %s\n", target)

			return nil // Found and extracted the binary
		}
	}

	return ErrBinaryNotFound
}

func (d *Downloader) extractFile(reader io.Reader, target string, mode int64) error {
	//nolint:gosec // Path is internally controlled
	openFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(mode))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	// Limit extraction size to prevent decompression bombs
	limitedReader := io.LimitReader(reader, maxExtractSize)
	if _, err := io.Copy(openFile, limitedReader); err != nil {
		_ = openFile.Close()

		return fmt.Errorf("failed to extract file: %w", err)
	}

	_ = openFile.Close()

	// Make executable on Unix systems
	if d.config.OS != "windows" {
		//nolint:gosec // Path is validated using filepath.Base() to prevent traversal
		if err := os.Chmod(target, executablePermission); err != nil {
			return fmt.Errorf("failed to set executable permissions: %w", err)
		}
	}

	return nil
}
