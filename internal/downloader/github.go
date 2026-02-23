package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	clientTimeout = 30 * time.Second

	githubAPIURL = "https://api.github.com/repos/golangci/golangci-lint/releases"
)

// GitHubRelease represents a GitHub release.
type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
}

// FetchAvailableVersions fetches available golangci-lint versions from GitHub.
func FetchAvailableVersions(limit int) ([]GitHubRelease, error) {
	url := fmt.Sprintf("%s?per_page=%d", githubAPIURL, limit)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{
		Timeout: clientTimeout,
	}

	//nolint:gosec // URL is constructed from hardcoded GitHub API URL constant
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("%w: status %d: %s", ErrGitHubAPI, resp.StatusCode, string(body))
	}

	var releases []GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Filter out drafts and prereleases
	var stableReleases []GitHubRelease

	for _, release := range releases {
		if !release.Draft && !release.Prerelease {
			stableReleases = append(stableReleases, release)
		}
	}

	return stableReleases, nil
}
