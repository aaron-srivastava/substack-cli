# Development Setup Guide

This guide covers setting up the substack-cli repository for development with all hooks and automation.

## Prerequisites

- **Go** 1.25.6 or later
- **Python** 3.7 or later (for pre-commit hooks)
- **Git** 2.13 or later

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/aaronsrivastava/substack-cli.git
cd substack-cli
```

### 2. Install Dependencies

```bash
# Install Go dependencies
go mod download
go mod tidy

# Install pre-commit framework
pip install pre-commit

# Install git hooks
pre-commit install
pre-commit install --hook-type commit-msg
```

This installs:
- **pre-commit**: Git hooks framework
- **commitlint**: Commit message validator
- **golangci-lint**: Go code linter
- **gosec**: Go security scanner
- General checks (trailing whitespace, etc.)

### 3. Verify Installation

```bash
# Check Go
go version

# Check Python/pre-commit
python --version
pre-commit --version

# Verify hooks installed
ls -la .git/hooks/
```

Should show:
```
.git/hooks/
├── pre-commit
├── commit-msg
└── ...
```

## Development Workflow

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feat/my-feature
   ```

2. **Make your changes** and test:
   ```bash
   go test ./...
   go build -o substack .
   ```

3. **Commit with conventional message**:
   ```bash
   git commit -m "feat: Add new feature"
   ```

   The pre-commit hook will validate the message automatically!

4. **Push and create PR**:
   ```bash
   git push origin feat/my-feature
   ```

### Commit Message Format

All commits must follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer]
```

**Valid examples:**
```bash
git commit -m "feat: Add post scheduling"
git commit -m "fix: Resolve markdown rendering"
git commit -m "docs: Update contributing guide"
git commit -m "perf: Optimize API calls"
```

See [COMMIT_CONVENTIONS.md](COMMIT_CONVENTIONS.md) for detailed guidelines.

## Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...

# Check linting
golangci-lint run ./...

# Build binary
go build -o substack .
```

## Git Hooks

### Pre-Commit Hook

Automatically validates commit messages before committing:

```bash
$ git commit -m "updated stuff"
# ✗ Error: Invalid commit format

$ git commit -m "feat: Added stuff"
# ✓ Commit valid
```

**The hook will:**
- Check type is valid (feat, fix, docs, etc.)
- Verify subject is lowercase and under 100 chars
- Ensure no period at end of subject
- Validate format matches Conventional Commits

**Skip validation (not recommended):**
```bash
git commit --no-verify -m "your message"
```

### Automatic Releases

When commits are pushed to `main`:

1. Workflow analyzes commits
2. Determines version bump (major/minor/patch)
3. Updates CHANGELOG.md
4. Creates and pushes git tag
5. Triggers full release (build, announce, etc.)

**No manual tagging needed!**

See [AUTO_RELEASE.md](AUTO_RELEASE.md) for details.

## Repository Structure

```
substack-cli/
├── cmd/                       # CLI commands
│   ├── auth.go
│   ├── config.go
│   ├── draft.go
│   ├── post.go
│   └── root.go
├── internal/                  # Internal packages
│   ├── api/                   # Substack API client
│   ├── auth/                  # Authentication & credential storage
│   ├── markdown/              # Markdown to ProseMirror converter
│   └── model/                 # Data types
├── .github/workflows/         # GitHub Actions workflows
│   ├── test.yml              # Tests and linting
│   ├── release.yml           # Manual release workflow
│   └── auto-release.yml      # Automatic releases on main
├── .husky/                   # Git hooks (managed by husky)
│   └── pre-commit            # Commit message validation
├── homebrew-tap/             # Homebrew formula directory
├── scripts/                  # Installation scripts
├── main.go                   # CLI entry point
├── go.mod                    # Go dependencies
├── package.json              # Node dependencies (for hooks)
├── commitlint.config.js      # Commit validation config
├── .goreleaser.yaml          # Release configuration
├── COMMIT_CONVENTIONS.md     # Commit guidelines
├── AUTO_RELEASE.md           # Auto-release documentation
├── CONTRIBUTING.md           # Contributor guidelines
└── README.md                 # Project overview
```

## Common Tasks

### Build Binary

```bash
go build -o substack .
./substack --version
```

### Run Tests

```bash
go test -v ./...
```

### Check Code Quality

```bash
# Linting
golangci-lint run ./...

# Go vet
go vet ./...

# Go fmt
go fmt ./...
```

### Update Dependencies

```bash
# Go
go get -u ./...
go mod tidy

# Node (if needed)
npm update
```

### Test Commit Validation

```bash
# Valid commit
git commit -m "feat: Add feature"  # ✓ Works

# Invalid commit
git commit -m "added feature"      # ✗ Fails pre-commit hook

# Check without committing
pre-commit run --all-files
```

### Make a Release

**Automatic** (on merge to main):
```bash
git merge feat/my-feature --ff-only
# Auto-release workflow runs automatically!
```

**Manual** (rare):
```bash
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
# GoReleaser workflow triggers
```

## Troubleshooting

### "pre-commit: command not found"

**Cause**: pre-commit not installed

**Fix**:
```bash
pip install pre-commit
pre-commit install
pre-commit install --hook-type commit-msg
```

### Pre-commit hook fails with "type-empty"

**Cause**: Commit message doesn't follow format

**Fix**: Use proper format:
```bash
git commit -m "feat: Add feature description"
```

### "golangci-lint: command not found"

**Cause**: Linter not installed

**Fix**:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Workflow shows "should-release: false"

**Cause**: Commits don't match conventional format

**Fix**: Push commits with proper types:
```bash
git commit -m "feat: New feature"
git push origin main
```

### Can't push because tag exists

**Cause**: Tag already created

**Fix**:
```bash
git tag -d v1.2.0
git push origin :refs/tags/v1.2.0
# Try again
```

## IDE Setup

### VSCode

Install extensions:
- **Go** (golang.go)

Create `.vscode/settings.json`:
```json
{
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go",
    "editor.formatOnSave": true
  }
}
```

### GoLand / IntelliJ IDEA

- Enable gofmt on save
- Run `go mod download` before opening project
- Mark `internal/` as Sources Root

## Environment Variables

For local development:

```bash
# Substack authentication (see README.md for setup)
export PUBLICATION_URL="https://your-publication.substack.com"
export USER_ID="your-user-id"
export SID="your-connect.sid"
export SUBSTACK_SID="your-substack.sid"
export SUBSTACK_LLI="your-substack.lli"
```

Or store in `.env` (don't commit):
```bash
# .env (add to .gitignore)
export PUBLICATION_URL="..."
source .env
```

## CI/CD

All workflows run automatically:

1. **test.yml** - On every push/PR
   - Runs tests
   - Linting
   - Commit validation
   - Builds binary

2. **auto-release.yml** - On push to main
   - Analyzes commits
   - Creates version tag
   - Updates changelog

3. **release.yml** - On version tag push
   - Builds all binaries
   - Creates GitHub release
   - Updates Homebrew
   - Announces on Bluesky

## Additional Resources

- [COMMIT_CONVENTIONS.md](COMMIT_CONVENTIONS.md) - Detailed commit guidelines
- [AUTO_RELEASE.md](AUTO_RELEASE.md) - Automatic release workflow
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines
- [RELEASING.md](RELEASING.md) - Manual release process
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Getting Help

- Check existing [Issues](https://github.com/aaronsrivastava/substack-cli/issues)
- Read documentation files (CONTRIBUTING.md, AUTO_RELEASE.md, etc.)
- Review commit history for examples: `git log --oneline -20`

---

Happy coding! 🚀
