package main

import "errors"

// Command-level errors.
var (
	// ErrVersionRequired is returned when a version argument is required but not provided.
	ErrVersionRequired = errors.New("version argument required")
)
