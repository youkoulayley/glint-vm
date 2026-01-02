package downloader

import "errors"

// Downloader-related errors.
var (
	// ErrBinaryNotFound is returned when the golangci-lint binary is not found after extraction.
	ErrBinaryNotFound = errors.New("golangci-lint binary not found after extraction")

	// ErrChecksumMismatch is returned when downloaded file checksum doesn't match expected.
	ErrChecksumMismatch = errors.New("checksum mismatch")

	// ErrHTTPRequest is returned when an HTTP request fails.
	ErrHTTPRequest = errors.New("HTTP request failed")

	// ErrGitHubAPI is returned when GitHub API returns an error.
	ErrGitHubAPI = errors.New("GitHub API error")

	// ErrInvalidKeepValue is returned when keep parameter is invalid.
	ErrInvalidKeepValue = errors.New("keep must be >= 0")
)
