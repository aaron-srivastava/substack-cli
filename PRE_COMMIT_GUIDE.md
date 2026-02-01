# Pre-Commit Framework Guide

This project uses [pre-commit](https://pre-commit.com/) for managing git hooks. Pre-commit is a language-agnostic framework that ensures code quality and commit message standards.

## What is Pre-Commit?

Pre-commit is a framework that:
- Manages git hooks automatically
- Runs checks before commits
- Prevents bad code from being committed
- Works with any language/project
- Has thousands of hooks available

No Node.js needed for Go projects!

## Installation

### Install pre-commit Framework

Choose your installation method:

**macOS (Homebrew):**
```bash
brew install pre-commit
```

**Linux (pip):**
```bash
pip install pre-commit
```

**Windows (Chocolatey):**
```bash
choco install pre-commit
```

**Any OS (pip):**
```bash
pip install pre-commit
```

### Install Git Hooks

After cloning the repository:

```bash
cd substack-cli
pre-commit install
pre-commit install --hook-type commit-msg
```

This installs hooks for:
- `pre-commit`: Runs before commits
- `commit-msg`: Validates commit message format

## What Hooks Are Installed?

Configuration in `.pre-commit-config.yaml`:

### Commit Message
- **commitlint**: Validates conventional commit format

### Go Code Quality
- **golangci-lint**: Comprehensive Go linter (auto-fixes issues)
- **golang-format**: Formats Go code
- **golang-vet**: Go vet checks
- **gosec**: Security scanner for Go code

### General Checks
- **trailing-whitespace**: Removes trailing whitespace
- **end-of-file-fixer**: Fixes missing newlines
- **check-yaml**: Validates YAML syntax
- **check-json**: Validates JSON syntax
- **check-merge-conflict**: Detects merge conflicts
- **check-added-large-files**: Warns about large files (>1MB)

## Usage

### Normal Workflow

```bash
# Make changes
echo "new code" >> main.go

# Stage changes
git add main.go

# Commit (hooks run automatically)
git commit -m "feat: Add new feature"

# Hooks run and validate:
Trim trailing whitespace.............................................. Passed
Fix end of file fixer..................................................... Passed
golangci-lint.............................................................Passed
golang-format..............................................................Passed
golang-vet................................................................Passed
gosec.....................................................................Passed
commitlint.................................................................Passed

# Commit successful! ✓
```

### Commit Message Validation

Pre-commit validates messages follow [Conventional Commits](https://www.conventionalcommits.org/):

**Valid messages:**
```bash
git commit -m "feat: Add feature"     # ✓
git commit -m "fix: Bug fix"          # ✓
git commit -m "docs: Update readme"   # ✓
git commit -m "perf: Optimize code"   # ✓
git commit -m "test: Add tests"       # ✓
```

**Invalid messages:**
```bash
git commit -m "updated stuff"         # ✗ No type
git commit -m "Added feature"         # ✗ Invalid type
git commit -m "FEAT: Add feature"     # ✗ Type not lowercase
git commit -m "feat: Add feature."    # ✗ Period at end
```

### Code Quality Checks

Pre-commit automatically fixes many issues:

```bash
# Has trailing whitespace
$ echo "code   " >> main.go
$ git add main.go
$ git commit -m "fix: Clean code"

Trim trailing whitespace.............................................. Fixed

# File was auto-fixed! Run git diff to see changes
$ git diff main.go
# Whitespace removed

# Now commit again
$ git add main.go
$ git commit -m "fix: Clean code"
# Success!
```

## Manual Pre-Commit Checks

Run checks without committing:

```bash
# Check all files
pre-commit run --all-files

# Check staged files only
pre-commit run

# Run specific hook
pre-commit run golangci-lint --all-files
pre-commit run commitlint --all-files

# Run on specific files
pre-commit run --files main.go cmd/post.go
```

Example output:

```bash
$ pre-commit run --all-files
Trim trailing whitespace.............................................. Passed
Fix end of file fixer..................................................... Passed
Check JSON..................................................................Passed
Check YAML..................................................................Passed
golangci-lint.............................................................Passed
golang-format..............................................................Passed
golang-vet................................................................Passed
gosec......................................................................Passed
commitlint.................................................................Skipped (commit-msg only)
```

## Skipping Hooks

### Skip All Hooks (Not Recommended)

```bash
git commit --no-verify -m "your message"
```

⚠️ **Warning**: Bypasses all validation. Your commit may fail CI checks.

### Skip Specific Hook

```bash
# Skip gosec (security check)
SKIP=gosec git commit -m "feat: Add feature"

# Skip multiple hooks
SKIP=golangci-lint,gosec git commit -m "feat: Add feature"
```

## Troubleshooting

### "pre-commit: command not found"

**Cause**: pre-commit not installed

**Fix:**
```bash
pip install pre-commit
pre-commit install
pre-commit install --hook-type commit-msg
```

### Hooks not running on commit

**Cause**: Hooks not installed

**Fix:**
```bash
pre-commit install
pre-commit install --hook-type commit-msg
```

### Hook failed but I want to commit anyway

**Option 1**: Fix the issue (recommended)
```bash
# Check what failed
pre-commit run --all-files

# Fix issues (many auto-fix)
# Re-stage and commit
git add .
git commit -m "feat: Add feature"
```

**Option 2**: Force commit (not recommended)
```bash
git commit --no-verify -m "feat: Add feature"
```

### "commitlint: command not found"

This shouldn't happen - commitlint runs via pre-commit. But if you need it:

```bash
pre-commit run commitlint --all-files
```

### Hooks running slowly

First time running hooks takes longer (downloading/installing). Subsequent runs are cached and faster.

To speed up:
```bash
# Skip some hooks
SKIP=gosec pre-commit run --all-files

# Or use autoupdate to refresh hook versions
pre-commit autoupdate
```

### Hooks modifying files I don't want changed

Pre-commit can auto-fix issues. Review changes:

```bash
git diff    # See what changed

# If acceptable, stage and commit
git add .
git commit -m "feat: Add feature"

# If not acceptable, restore and fix manually
git checkout -- .  # Discard changes
# Make manual fixes
git add .
git commit --no-verify -m "feat: Add feature"
```

## Configuration

Configuration file: `.pre-commit-config.yaml`

Key sections:

```yaml
repos:
  - repo: https://github.com/compilerla/pre-commit-commitlint
    rev: v19.0.0              # Version of hook
    hooks:
      - id: commitlint
        stages: [commit-msg]  # When to run (commit-msg, commit, etc.)
```

### Updating Hooks

Update all hooks to latest versions:

```bash
pre-commit autoupdate
```

Update specific hook:

```bash
pre-commit autoupdate --repo https://github.com/golangci/golangci-lint
```

## GitHub Actions Integration

Pre-commit runs on every pull request via GitHub Actions (see `.github/workflows/test.yml`).

If your local pre-commit passes but CI fails:
1. Ensure you installed hooks: `pre-commit install`
2. Run locally: `pre-commit run --all-files`
3. Fix issues and commit

## Best Practices

1. **Always install hooks after cloning**
   ```bash
   pre-commit install
   pre-commit install --hook-type commit-msg
   ```

2. **Let pre-commit auto-fix issues** when possible
   ```bash
   git add .
   git commit -m "feat: Add feature"
   # Auto-fixes happen, commit succeeds
   ```

3. **Check hook output** to understand failures
   ```bash
   pre-commit run --all-files  # See what failed and why
   ```

4. **Don't ignore hook failures**
   Using `--no-verify` bypasses safety checks

5. **Keep commits clean**
   Use proper conventional commit format:
   ```bash
   git commit -m "feat: Descriptive message"
   ```

## Advanced Usage

### Run hooks only on changed files

```bash
pre-commit run  # Checks staged files
```

### Run specific hook

```bash
pre-commit run golangci-lint
```

### Install hooks from different branch

```bash
pre-commit install --hook-type commit-msg
pre-commit run -a  # Run all on all files
```

### View installed hooks

```bash
pre-commit run --list
```

### Clean cached data

```bash
pre-commit clean
```

## Debugging

Enable verbose output:

```bash
pre-commit run -v --all-files
```

See what each hook is doing:

```bash
pre-commit run --show-diff-on-failure --all-files
```

## See Also

- [Pre-commit documentation](https://pre-commit.com/)
- [COMMIT_CONVENTIONS.md](COMMIT_CONVENTIONS.md) - Commit message format
- [AUTO_RELEASE.md](AUTO_RELEASE.md) - Automatic releases
- [SETUP.md](SETUP.md) - Development setup
- [.pre-commit-config.yaml](.pre-commit-config.yaml) - Hook configuration
