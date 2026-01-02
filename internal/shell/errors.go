package shell

import "errors"

// Shell-related errors.
var (
	// ErrUnsupportedShell is returned when an unsupported shell is requested.
	ErrUnsupportedShell = errors.New("unsupported shell (supported: bash, zsh)")
)
