package shell

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	bashShell = "bash"
	zshShell  = "zsh"
)

// InitOptions contains configuration for shell initialization.
type InitOptions struct {
	AutoSwitch bool // Enable auto-switching on directory change
}

// Integration provides shell-specific integration.
type Integration struct {
	shellType string
}

// NewIntegrator creates a new shell integrator based on the shell name.
func NewIntegrator(shellName string) (*Integration, error) {
	switch shellName {
	case bashShell, zshShell:
		return &Integration{shellType: shellName}, nil
	default:
		return nil, fmt.Errorf("unsupported shell %s: %w", shellName, ErrUnsupportedShell)
	}
}

// Name returns the shell name.
func (i *Integration) Name() string {
	return i.shellType
}

// GenerateInit outputs shell initialization code.
func (i *Integration) GenerateInit(opts InitOptions) string {
	switch i.shellType {
	case bashShell:
		return (&BashShell{}).GenerateInit(opts)
	case zshShell:
		return (&ZshShell{}).GenerateInit(opts)
	default:
		return ""
	}
}

// GenerateUse outputs PATH modification commands.
func (i *Integration) GenerateUse(version string) string {
	switch i.shellType {
	case bashShell:
		return (&BashShell{}).GenerateUse(version)
	case zshShell:
		return (&ZshShell{}).GenerateUse(version)
	default:
		return ""
	}
}

// GetGlintVMRoot returns the root directory for glint-vm cache.
func GetGlintVMRoot() string {
	if root := os.Getenv("GLINT_VM_ROOT"); root != "" {
		return root
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		// Fallback to home directory
		home, _ := os.UserHomeDir()

		return filepath.Join(home, ".cache", "glint-vm")
	}

	return filepath.Join(cacheDir, "glint-vm")
}

// GetCurrentDir returns the directory containing the current version symlink.
func GetCurrentDir() string {
	return filepath.Join(GetGlintVMRoot(), "current")
}

// GetCurrentBinaryPath returns the path to the current golangci-lint binary.
func GetCurrentBinaryPath() string {
	return filepath.Join(GetCurrentDir(), "golangci-lint")
}

// DetectShell attempts to detect the current shell from environment.
func DetectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return bashShell // Default fallback
	}

	// Extract shell name from path (e.g., /bin/bash -> bash)
	return filepath.Base(shell)
}
