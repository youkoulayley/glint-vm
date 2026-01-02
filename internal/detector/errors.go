package detector

import "errors"

// Detector-related errors.
var (
	// ErrNotDirectory is returned when a path is not a directory.
	ErrNotDirectory = errors.New("path is not a directory")
)
