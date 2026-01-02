package downloader

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/youkoulayley/glint-vm/internal/config"
)

// CachedVersion represents a cached golangci-lint version.
type CachedVersion struct {
	Version    string
	Path       string
	BinaryPath string
	Size       int64
	ModTime    time.Time
	IsComplete bool // Whether the binary exists and is executable
}

// CacheManager manages cached golangci-lint versions.
type CacheManager struct {
	config *config.Config
}

// NewCacheManager creates a new cache manager.
func NewCacheManager() (*CacheManager, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	return &CacheManager{
		config: cfg,
	}, nil
}

// List returns all cached versions.
func (cm *CacheManager) List() ([]*CachedVersion, error) {
	versionsDir := cm.config.GetVersionsDir()

	// Check if versions directory exists
	if _, err := os.Stat(versionsDir); os.IsNotExist(err) {
		return []*CachedVersion{}, nil
	}

	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read versions directory: %w", err)
	}

	versions := make([]*CachedVersion, 0, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		version := entry.Name()
		versionPath := cm.config.GetVersionDir(version)
		binaryPath := cm.config.GetBinaryPath(version)

		cached := &CachedVersion{
			Version:    version,
			Path:       versionPath,
			BinaryPath: binaryPath,
			IsComplete: cm.config.BinaryExists(version),
		}

		// Get binary info if it exists
		if info, err := os.Stat(binaryPath); err == nil {
			cached.Size = info.Size()
			cached.ModTime = info.ModTime()
		} else {
			// Get directory info as fallback
			if info, err := os.Stat(versionPath); err == nil {
				cached.ModTime = info.ModTime()
			}
		}

		versions = append(versions, cached)
	}

	// Sort by semantic version (latest first)
	sort.Slice(versions, func(i, j int) bool {
		vi, erri := semver.NewVersion(versions[i].Version)
		vj, errj := semver.NewVersion(versions[j].Version)

		// If both versions are valid semver, compare them
		if erri == nil && errj == nil {
			return vi.GreaterThan(vj)
		}

		// Fallback: invalid versions go to the end, or compare by ModTime
		if erri != nil && errj != nil {
			return versions[i].ModTime.After(versions[j].ModTime)
		}

		return erri == nil // Valid version comes first
	})

	return versions, nil
}

// Remove removes a specific cached version.
func (cm *CacheManager) Remove(version string) error {
	version = config.NormalizeVersion(version)

	versionDir := cm.config.GetVersionDir(version)
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not cached: %w", version, err)
	}

	err := os.RemoveAll(versionDir)
	if err != nil {
		return fmt.Errorf("failed to remove version %s: %w", version, err)
	}

	return nil
}

// RemoveAll removes all cached versions.
func (cm *CacheManager) RemoveAll() (int, error) {
	versions, err := cm.List()
	if err != nil {
		return 0, err
	}

	removed := 0

	var lastErr error

	for _, v := range versions {
		err := cm.Remove(v.Version)
		if err != nil {
			lastErr = err
		} else {
			removed++
		}
	}

	return removed, lastErr
}

// RemoveOldest removes all but the N latest versions (by semantic version).
func (cm *CacheManager) RemoveOldest(keep int) (int, error) {
	if keep < 0 {
		return 0, ErrInvalidKeepValue
	}

	versions, err := cm.List()
	if err != nil {
		return 0, err
	}

	if len(versions) <= keep {
		return 0, nil // Nothing to remove
	}

	// versions is already sorted by semantic version (latest first)
	toRemove := versions[keep:]

	removed := 0

	var lastErr error

	for _, v := range toRemove {
		err := cm.Remove(v.Version)
		if err != nil {
			lastErr = err
		} else {
			removed++
		}
	}

	return removed, lastErr
}

// RemoveIncomplete removes incomplete/corrupted versions (directories without valid binaries).
func (cm *CacheManager) RemoveIncomplete() (int, error) {
	versions, err := cm.List()
	if err != nil {
		return 0, err
	}

	removed := 0

	var lastErr error

	for _, v := range versions {
		if !v.IsComplete {
			err := cm.Remove(v.Version)
			if err != nil {
				lastErr = err
			} else {
				removed++
			}
		}
	}

	return removed, lastErr
}

// GetTotalSize returns the total size of all cached binaries in bytes.
func (cm *CacheManager) GetTotalSize() (int64, error) {
	versions, err := cm.List()
	if err != nil {
		return 0, err
	}

	var total int64
	for _, v := range versions {
		total += v.Size
	}

	return total, nil
}

// EnsureVersionDir ensures the directory for a version exists with proper permissions.
func (cm *CacheManager) EnsureVersionDir(version string) error {
	return fmt.Errorf("ensure version dir: %w", cm.config.EnsureVersionDir(version))
}

// GetVersionDir returns the directory path for a version.
func (cm *CacheManager) GetVersionDir(version string) string {
	return cm.config.GetVersionDir(version)
}

// GetBinaryPath returns the binary path for a version.
func (cm *CacheManager) GetBinaryPath(version string) string {
	return cm.config.GetBinaryPath(version)
}

// IsCached checks if a version is cached and complete.
func (cm *CacheManager) IsCached(version string) bool {
	return cm.config.BinaryExists(version)
}
