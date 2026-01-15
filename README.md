# glint-vm

A version manager for golangci-lint. Manage multiple golangci-lint versions and switch between them easily.

## Installation

```bash
go install github.com/youkoulayley/glint-vm/cmd/glint-vm@latest
```

Or build from source:

```bash
git clone https://github.com/youkoulayley/glint-vm.git
cd glint-vm
make build
sudo make install
```

## Shell Configuration

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

## Usage

**Detect and use the version from your project:**
```bash
glint-vm detect --use
```

**Install a specific version:**
```bash
glint-vm install v1.55.2
glint-vm install --use v1.55.2  # Install and activate
```

**Switch to a version:**
```bash
glint-vm use v1.55.2
```

**List installed versions:**
```bash
glint-vm list
```

**Show current version:**
```bash
glint-vm current
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

## License

MIT
