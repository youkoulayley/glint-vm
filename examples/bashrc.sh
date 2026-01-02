# glint-vm initialization for bash
#
# Add this to your ~/.bashrc to enable glint-vm

# Basic initialization (recommended for most users)
eval "$(glint-vm init bash)"

# Alternative: Enable automatic version switching
# Uncomment the line below to automatically switch golangci-lint versions
# when you cd into directories with .golangci-lint.version files
# eval "$(glint-vm init --auto-switch bash)"
