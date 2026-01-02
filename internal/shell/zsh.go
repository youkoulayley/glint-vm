package shell

import (
	"fmt"
	"strings"
)

const zsh = "zsh"

// ZshShell implements Integrator for zsh.
type ZshShell struct{}

// Name returns the shell name.
func (z *ZshShell) Name() string {
	return zsh
}

// GenerateInit generates zsh initialization code.
func (z *ZshShell) GenerateInit(opts InitOptions) string {
	root := GetGlintVMRoot()
	currentDir := GetCurrentDir()

	var builder strings.Builder

	// Export GLINT_VM_ROOT
	builder.WriteString(fmt.Sprintf("export GLINT_VM_ROOT=\"%s\"\n", root))

	// Add current directory to PATH
	builder.WriteString(fmt.Sprintf("export PATH=\"%s:$PATH\"\n", currentDir))

	// Add wrapper function that auto-evals commands that modify the environment
	builder.WriteString(generateWrapperFunction())

	// Add auto-switch hook if requested
	if opts.AutoSwitch {
		builder.WriteString(generateZshAutoSwitch())
	}

	return builder.String()
}

// GenerateUse generates zsh commands to activate a specific version.
func (z *ZshShell) GenerateUse(version string) string {
	currentDir := GetCurrentDir()

	var builder strings.Builder

	// Export PATH (reinforcing the path in case it's not in init)
	builder.WriteString(fmt.Sprintf("export PATH=\"%s:$PATH\"\n", currentDir))

	// Export GLINT_VM_VERSION
	builder.WriteString(fmt.Sprintf("export GLINT_VM_VERSION=\"%s\"\n", version))

	// Echo success message to stderr (so stdout is clean for eval)
	builder.WriteString(fmt.Sprintf("echo \"Switched to golangci-lint %s\" >&2\n", version))

	return builder.String()
}

// generateZshAutoSwitch returns the zsh auto-switch hook code.
func generateZshAutoSwitch() string {
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

# Hook into chpwd (zsh directory change hook)
autoload -U add-zsh-hook
add-zsh-hook chpwd _glint_vm_auto_switch

# Also run on shell initialization
_glint_vm_auto_switch
`
}
