package config

import "errors"

// Config-related errors.
var (
	// ErrVersionNotInstalled is returned when a requested version is not installed.
	ErrVersionNotInstalled = errors.New("version is not installed")
)
