// Package shell provides shell integration functionality for glint-vm,
// supporting bash and zsh with automatic version switching capabilities.
package shell

import (
	"fmt"
	"strings"
)

const bash = "bash"

// BashShell implements Integrator for bash.
type BashShell struct{}

// Name returns the shell name.
func (b *BashShell) Name() string {
	return bash
}

// GenerateInit generates bash initialization code.
func (b *BashShell) GenerateInit(opts InitOptions) string {
	root := GetGlintVMRoot()
	currentDir := GetCurrentDir()

	var builder strings.Builder

	// Export GLINT_VM_ROOT
	fmt.Fprintf(&builder, "export GLINT_VM_ROOT=%q\n", root)

	// Add current directory to PATH
	fmt.Fprintf(&builder, "export PATH=\"%s:$PATH\"\n", currentDir)

	// Add wrapper function that auto-evals commands that modify the environment
	builder.WriteString(generateWrapperFunction())

	// Add auto-switch hook if requested
	if opts.AutoSwitch {
		builder.WriteString(generateBashAutoSwitch())
	}

	return builder.String()
}

// GenerateUse generates bash commands to activate a specific version.
func (b *BashShell) GenerateUse(version string) string {
	currentDir := GetCurrentDir()

	var builder strings.Builder

	// Export PATH (reinforcing the path in case it's not in init)
	fmt.Fprintf(&builder, "export PATH=\"%s:$PATH\"\n", currentDir)

	// Export GLINT_VM_VERSION
	fmt.Fprintf(&builder, "export GLINT_VM_VERSION=%q\n", version)

	// Echo success message to stderr (so stdout is clean for eval)
	fmt.Fprintf(&builder, "echo \"Switched to golangci-lint %s\" >&2\n", version)

	return builder.String()
}

// generateWrapperFunction returns the shell wrapper function code.
func generateWrapperFunction() string {
	return `
# glint-vm wrapper function - automatically handles environment updates
glint-vm() {
  local command="${1:-}"

  # Commands that modify the shell environment need to be eval'd
  if [[ "$command" == "use" ]] || \
     [[ "$command" == "install" && "$*" =~ "--use" ]] || \
     [[ "$command" == "detect" && "$*" =~ "--use" ]]; then
    # Capture stdout only (stderr goes to terminal for informational messages)
    local output
    output=$(command glint-vm "$@")
    local exit_code=$?

    if [[ $exit_code -eq 0 ]]; then
      eval "$output"
    else
      return $exit_code
    fi
  else
    # For other commands, just pass through
    command glint-vm "$@"
  fi
}
`
}

// generateBashAutoSwitch returns the bash auto-switch hook code.
func generateBashAutoSwitch() string {
	return `
# Auto-switch golangci-lint version based on detected configuration
_glint_vm_auto_switch() {
  # Detect version from all sources (CI configs, Makefiles, etc.)
  local version=$(command glint-vm detect --quiet 2>/dev/null)

  # Only switch if a version was detected and it's different from current
  if [[ -n "$version" && "$version" != "$GLINT_VM_VERSION" ]]; then
    # Check if version is already cached to avoid blocking download
    local cache_dir="${GLINT_VM_ROOT}/versions/${version}"
    local binary_path="${cache_dir}/golangci-lint"

    if [[ -f "$binary_path" ]]; then
      # Version is cached, switch to it
      local switch_output=$(command glint-vm use "$version" 2>&1)
      if [[ $? -eq 0 ]]; then
        eval "$switch_output"
      fi
    else
      # Version not cached - show one-time message
      if [[ "$_GLINT_VM_NOTIFIED_VERSION" != "$version" ]]; then
        echo "glint-vm: detected $version but it's not installed yet" >&2
        echo "glint-vm: run 'glint-vm install --use $version' to install and activate it" >&2
        export _GLINT_VM_NOTIFIED_VERSION="$version"
      fi
    fi
  fi
}

# Hook into PROMPT_COMMAND
if [[ ! "$PROMPT_COMMAND" =~ _glint_vm_auto_switch ]]; then
  PROMPT_COMMAND="_glint_vm_auto_switch${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
fi

# Also run on shell initialization
_glint_vm_auto_switch
`
}
