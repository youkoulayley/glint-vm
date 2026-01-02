# glint-vm

**golangci-lint version manager** - Like gvm/nvm, but for golangci-lint. Manage multiple golangci-lint versions in your shell environment.

## Features

- 🎯 **Auto-detection** - Scans `.golangci-lint.version`, GitHub Actions, Makefile, and other CI configs
- 📦 **Smart caching** - Downloads once, uses everywhere
- 🔄 **Per-shell activation** - Each terminal can use a different version
- 🚀 **Direct execution** - Run `golangci-lint` directly, not through a wrapper
- 🌍 **Cross-platform** - Linux, macOS support (bash and zsh)
- 🛡️ **XDG compliant** - Respects `$XDG_CACHE_HOME`
- ⚡ **Auto-switching** - Optional automatic version switching on directory change
- 🎨 **No eval needed** - Shell wrapper handles environment updates automatically
- ⚙️ **One-command setup** - `detect --use` for instant setup from project configs

## Installation

### From source

```bash
go install github.com/youkoulayley/glint-vm/cmd/glint-vm@latest
```

### Manual build

```bash
git clone https://github.com/youkoulayley/glint-vm.git
cd glint-vm
make build
sudo make install  # Installs to /usr/local/bin
```

## Quick Start

### 1. One-time shell setup

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
# For bash
eval "$(glint-vm init bash)"

# For zsh
eval "$(glint-vm init zsh)"
```

Then reload your shell:
```bash
source ~/.bashrc  # or source ~/.zshrc
```

### 2. Detect and activate your project's version

The easiest way - let glint-vm detect the version from your project:

```bash
cd your-project
glint-vm detect --use
```

This will:
1. Detect the version from your project config (Makefile, GitHub Actions, `.golangci-lint.version`, etc.)
2. Download it if not cached
3. Activate it automatically

**Or specify a version manually:**

```bash
# Install and activate in one command
glint-vm install --use v1.55.2

# Or install first, then activate
glint-vm install v1.55.2
glint-vm use v1.55.2
```

> **Note:** No `eval` needed! The shell wrapper handles environment updates automatically.

### 3. Use golangci-lint directly

```bash
golangci-lint run
golangci-lint run --enable-all
golangci-lint run --fix
```

That's it! The `golangci-lint` command is now available in your PATH.

## Usage

### Shell Integration

glint-vm manages golangci-lint versions through shell environment manipulation:

```bash
# Generate shell initialization code
glint-vm init bash           # For bash
glint-vm init zsh            # For zsh

# With auto-switching enabled
glint-vm init bash --auto-switch
glint-vm init zsh --auto-switch
```

The `init` command outputs shell code that:
- Adds `~/.cache/glint-vm/current` to your `$PATH`
- Sets up environment variables
- Defines helper functions
- Optionally enables automatic version switching

### Version Management

**Activate a version:**
```bash
glint-vm use v1.55.2
```

This command:
- Downloads the version if not cached
- Updates the symlink at `~/.cache/glint-vm/current/golangci-lint`
- Sets `$GLINT_VM_VERSION` environment variable
- Makes `golangci-lint` available in your `$PATH`

> The shell wrapper automatically handles environment updates - no `eval` needed!

**Check current version:**
```bash
glint-vm current
```

**List installed versions:**
```bash
glint-vm list
```

Output:
```
Installed versions:
  v1.64.5
* v1.55.2
  v1.54.0

* = current version (v1.55.2)
```

**Uninstall a version:**
```bash
glint-vm uninstall v1.54.0
```

### Version Detection

See what version will be detected from your project:

```bash
glint-vm detect
```

**Output:**
```
Cache directory: /home/user/.cache/glint-vm
Platform: linux-amd64

Detecting golangci-lint version...

✓ Detected version: v1.55.2
  Source: /path/to/project/.golangci-lint.version
  Source type: version-file
  Line: 1
  Pattern: plain-version

✓ Binary cached at: /home/user/.cache/glint-vm/versions/v1.55.2/golangci-lint
```

### Auto-Switching (Optional)

Enable automatic version switching when you change directories:

```bash
# Add to ~/.bashrc
eval "$(glint-vm init bash --auto-switch)"

# Add to ~/.zshrc
eval "$(glint-vm init zsh --auto-switch)"
```

With auto-switching enabled, glint-vm will:
- Automatically detect the version from **all sources** (`.golangci-lint.version`, GitHub Actions, Makefiles, CI configs)
- Switch instantly if the version is already cached
- Show a helpful message if the version needs to be installed first (non-blocking!)

**Example:**
```bash
cd my-project
# If v1.60.0 is already installed: switches instantly
# Switched to golangci-lint v1.60.0

# If not installed yet: shows helpful message
# glint-vm: detected v1.60.0 but it's not installed yet
# glint-vm: run 'glint-vm install --use v1.60.0' to install and activate it
```

### Cache Management

**List cached versions:**
```bash
glint-vm cache list
```

**Output:**
```
Cache directory: /home/user/.cache/glint-vm/versions

Cached versions (2 total, 48.78 MB):
  ✓ v1.55.2 [24.39 MB]
  ✓ v1.54.0 [24.39 MB]
```

**Clean up cache:**
```bash
# Remove all cached versions
glint-vm cache clean --all

# Keep 3 most recent versions
glint-vm cache clean --keep 3
```

## Common Workflows

### New Project Setup

When you join a new project:

```bash
cd new-project

# One command to get started
glint-vm detect --use
# Detects v1.60.0 from .golangci-lint.version
# Downloads and activates it
# Switched to golangci-lint v1.60.0

# Start linting immediately
golangci-lint run
```

### Working on Multiple Projects

Different projects, different versions - no problem:

```bash
# Terminal 1 - Project A
cd project-a
glint-vm use v1.55.2
golangci-lint run

# Terminal 2 - Project B
cd project-b
glint-vm use v1.60.0
golangci-lint run

# Each terminal maintains its own version!
```

### Auto-Switching Workflow

Enable auto-switching for seamless project hopping:

```bash
# One-time setup
eval "$(glint-vm init bash --auto-switch)"

# Now switching projects is automatic
cd project-a
# Switched to golangci-lint v1.55.2

cd ../project-b
# Switched to golangci-lint v1.60.0

# No manual version management needed!
```

### Quick Version Check

See what's detected without installing:

```bash
# See what version is configured
glint-vm detect

# Just the version number (for scripts)
glint-vm detect --quiet

# Install only if you want it
glint-vm detect --install
```

## Version Detection

glint-vm automatically detects the golangci-lint version from your project configuration files in this priority order:

1. **`.golangci-lint.version`** file (highest priority)
2. **GitHub Actions** - `.github/workflows/*.yml`
3. **Semaphore CI** - `.semaphore/semaphore.yml`
4. **Makefile** - `Makefile`, `makefile`, `GNUmakefile`
5. **CircleCI** - `.circleci/config.yml`
6. **GitLab CI** - `.gitlab-ci.yml`

### Supported Version Patterns

glint-vm recognizes these patterns:

```yaml
# GitHub Actions
- uses: golangci/golangci-lint-action@v3
  with:
    version: v1.55.2

# Makefile
GOLANGCI_LINT_VERSION := v1.55.2

# Docker
image: golangci/golangci-lint:v1.55.2

# Environment variables
GOLANGCI_LINT_VERSION=v1.55.2

# Plain version file
v1.55.2
```

## Configuration

### Cache Directory

By default, glint-vm caches binaries in:
- Linux/macOS: `$XDG_CACHE_HOME/glint-vm` or `~/.cache/glint-vm`

Override with environment variable:

```bash
export GLINT_VM_ROOT=/custom/cache/glint-vm
```

### Directory Structure

```
~/.cache/glint-vm/
├── versions/
│   ├── v1.55.2/golangci-lint
│   └── v1.54.0/golangci-lint
└── current/
    └── golangci-lint → ../versions/v1.55.2/golangci-lint
```

## CI/CD Integration

In CI environments, you can use glint-vm for consistent linter versions:

### GitHub Actions

```yaml
- name: Install glint-vm
  run: go install github.com/youkoulayley/glint-vm/cmd/glint-vm@latest

- name: Setup shell and golangci-lint
  run: |
    eval "$(glint-vm init bash)"
    glint-vm install --use v1.55.2

- name: Run linter
  run: golangci-lint run
```

### GitLab CI

```yaml
lint:
  before_script:
    - go install github.com/youkoulayley/glint-vm/cmd/glint-vm@latest
    - eval "$(glint-vm init bash)"
    - glint-vm install --use v1.55.2
  script:
    - golangci-lint run
```

## Commands Reference

### `glint-vm init <bash|zsh>`
Generate shell initialization code
- `--auto-switch` - Enable automatic version switching

**Example:**
```bash
eval "$(glint-vm init bash)"
```

### `glint-vm use <version>`
Activate a specific version in current shell

**Example:**
```bash
glint-vm use v1.55.2
```

> No `eval` needed - the shell wrapper handles it automatically!

### `glint-vm current`
Show currently active version

### `glint-vm list`
List installed versions (marks current version with `*`)

### `glint-vm list-remote`
List available golangci-lint versions from GitHub
- `--limit N, -l N` - Limit the number of versions to display (default: 20)
- Alias: `lr`

**Examples:**
```bash
# List 20 latest versions
glint-vm list-remote

# List 10 latest versions
glint-vm list-remote --limit 10

# Using alias
glint-vm lr -l 5
```

### `glint-vm install <version>`
Download a specific golangci-lint version
- `--use, -u` - Activate the version after installing

**Examples:**
```bash
# Just download
glint-vm install v1.55.2

# Download and activate in one command
glint-vm install --use v1.55.2
```

### `glint-vm uninstall <version>`
Remove a specific version

### `glint-vm detect`
Show detected version and source file
- `--quiet, -q` - Output only the version number (for scripting)
- `--install, -i` - Download the detected version
- `--use, -u` - Download and activate the detected version

**Examples:**
```bash
# Full output with details
glint-vm detect

# Just the version (for scripts)
glint-vm detect --quiet

# Detect and install in one command
glint-vm detect --install

# Detect, install, and activate in one command
glint-vm detect --use
```

> The `--use` flag is handled automatically by the shell wrapper - no `eval` needed!

### `glint-vm cache list`
List all cached versions with sizes

### `glint-vm cache clean`
Remove cached versions
- `--all` - Remove all versions
- `--keep N` - Keep N most recent versions (default: 3)

### `glint-vm version`
Show glint-vm version information

## Development

### Build from source

```bash
make build        # Build for current platform
make install      # Install to /usr/local/bin
make test         # Run all tests
make clean        # Clean build artifacts
```

### Run tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/detector/...

# With coverage
go test -cover ./...
```

## How It Works

1. **Shell Integration**: `glint-vm init` adds `~/.cache/glint-vm/current` to your `$PATH`
2. **Version Activation**: `glint-vm use` creates a symlink at `current/golangci-lint` pointing to the specific version
3. **Direct Execution**: You run `golangci-lint` directly from your shell - no wrapper involved
4. **Caching**: Binaries are stored in `~/.cache/glint-vm/versions/`

This architecture is similar to how `gvm`, `nvm`, and `rbenv` work - managing your shell environment rather than wrapping commands.

## Why glint-vm?

- **Consistency** - Ensure same linter version across local dev and CI
- **Simplicity** - One command to detect, download, and activate (`glint-vm detect --use`)
- **Flexibility** - Each terminal can use a different version
- **Speed** - Smart caching means download once, use everywhere
- **No eval needed** - Shell wrapper handles environment updates automatically
- **Auto-switching** - Optionally switch versions automatically when changing directories
- **Non-blocking** - Never wait on downloads; install when you're ready
- **Compatibility** - Works with existing project configurations
- **Familiar** - Same workflow as gvm, nvm, rbenv

## Examples

See the [`examples/`](examples/) directory for shell configuration examples:
- [`examples/bashrc.sh`](examples/bashrc.sh) - Bash configuration
- [`examples/zshrc.sh`](examples/zshrc.sh) - Zsh configuration

## License

MIT

## Contributing

Contributions welcome! Please open an issue or PR.

## Similar Tools

- [asdf-golangci-lint](https://github.com/hypnoglow/asdf-golangci-lint) - asdf plugin
- [golangci-lint install script](https://golangci-lint.run/usage/install/) - Official installer

glint-vm provides automatic version detection, per-shell version management, and optional auto-switching, making it ideal for teams and multi-project workflows.
