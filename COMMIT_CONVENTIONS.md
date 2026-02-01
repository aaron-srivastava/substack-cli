# Commit Conventions

This project enforces [Conventional Commits](https://www.conventionalcommits.org/) to ensure consistent commit messages, enable automatic changelog generation, and support semantic versioning.

## Setup

### Install pre-commit

Install the [pre-commit](https://pre-commit.com/) framework:

**macOS (with Homebrew):**
```bash
brew install pre-commit
```

**Linux:**
```bash
pip install pre-commit
```

**Windows (with Chocolatey):**
```bash
choco install pre-commit
```

Or install from source:
```bash
pip install pre-commit
```

### Install Hooks

After cloning the repository, install the git hooks:

```bash
cd substack-cli
pre-commit install
pre-commit install --hook-type commit-msg
```

This installs:
- **commitlint**: Validates commit message format
- **golangci-lint**: Lints Go code
- **gosec**: Security scanner for Go
- General checks (trailing whitespace, file endings, etc.)

## Commit Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

Each commit must start with one of these types:

| Type | Description | Version Impact |
|------|-------------|-----------------|
| **feat** | A new feature | Minor (1.0.0 → 1.1.0) |
| **fix** | A bug fix | Patch (1.0.0 → 1.0.1) |
| **docs** | Documentation changes | No auto-release |
| **style** | Code style (formatting, missing semicolons, etc.) | Patch |
| **refactor** | Code refactoring without feature changes | Patch |
| **perf** | Performance improvements | Patch |
| **test** | Adding or updating tests | No auto-release |
| **ci** | CI/CD configuration changes | No auto-release |
| **chore** | Other changes (dependencies, etc.) | No auto-release |
| **revert** | Revert a previous commit | Patch |

### Scope (Optional)

Scope specifies what part of the code is affected:

```bash
feat(markdown): Add table support
fix(api): Handle null responses
docs(readme): Update installation
```

### Breaking Changes

Add a `!` before the colon or use `BREAKING CHANGE:` footer to indicate breaking changes (triggers MAJOR version bump):

```bash
feat!: Change authentication method
# or
feat: Redesign API endpoints

BREAKING CHANGE: Old endpoints no longer supported
```

## Examples

### New Feature

```bash
git commit -m "feat: Add scheduled post support

- Implement scheduling API
- Add CLI flag --schedule
- Store scheduled posts in database"
```

**Result**: v1.0.0 → v1.1.0 (MINOR)

### Bug Fix

```bash
git commit -m "fix: Resolve markdown list rendering issue

The list items were not properly indented.
This fix ensures proper spacing."
```

**Result**: v1.0.0 → v1.0.1 (PATCH)

### Breaking Change

```bash
git commit -m "feat!: Replace XML config with YAML

BREAKING CHANGE: XML configuration no longer supported.
Users must migrate to YAML format."
```

**Result**: v1.0.0 → v2.0.0 (MAJOR)

### With Scope

```bash
feat(cli): Add progress bar to uploads
fix(api): Handle connection timeouts
docs(contributing): Add example commits
```

### Multiple Paragraphs

```bash
git commit -m "feat: Add multi-account support

Users can now manage multiple Substack accounts
with the CLI. Accounts are stored encrypted in
~/.config/substack-cli/accounts.json

Changes:
- Add auth login/logout/switch commands
- Implement account storage
- Add account selection prompt

Fixes: #123, #456"
```

## Pre-Commit Hooks

The pre-commit framework automatically validates your commits **before** they're created. If the message or code doesn't pass checks, the commit is rejected:

### Commit Message Validation

```bash
$ git commit -m "updated stuff"

Commit message validation...
⧗   input: updated stuff
✖   subject case must not be start-case, pascal-case, upper-case
✖   type must be one of [feat, fix, ...] [type-empty]

✖   found 2 problems, 0 warnings

commit-msg hook failed
```

**Valid commit:**

```bash
$ git commit -m "fix: Update dependencies"

Commit message validation...
✔   Commit message valid
```

### Code Quality Checks

Pre-commit also runs on code files:

```bash
$ git add main.go && git commit -m "feat: Add new feature"

Trim trailing whitespace.............................................. Passed
Fix end of file fixer..................................................... Passed
Check JSON........................................................Passed
golangci-lint.............................................................Passed
golang-format..............................................................Passed
golang-vet................................................................Passed
gosec.....................................................................Passed
commitlint.................................................................Passed

✓ All hooks passed
```

### Skip Hooks (Not Recommended)

To skip validation:

```bash
git commit --no-verify -m "your message"
```

⚠️ **Warning**: This bypasses all checks. Your commit may fail CI.

## GitHub Actions Validation

All pull requests are validated to ensure commit messages follow the convention. Check the "Checks" section on your PR for validation results.

## Auto-Release on Main

When commits are pushed to the `main` branch, the auto-release workflow:

1. **Analyzes commits** since the last tag
2. **Determines version bump**:
   - Breaking changes (`!`) → MAJOR (1.0.0 → 2.0.0)
   - Features (`feat:`) → MINOR (1.0.0 → 1.1.0)
   - Fixes (`fix:`) → PATCH (1.0.0 → 1.0.1)
3. **Updates CHANGELOG.md**
4. **Creates and pushes tag**
5. **Triggers release workflow** (builds, announces, etc.)

**No manual tagging needed!** Just merge to main with proper commits.

## Changelog

The changelog is automatically generated from commit messages using semantic commit structure:

- **Features** section: All `feat:` commits
- **Bug fixes** section: All `fix:` commits
- **Performance** section: All `perf:` commits
- etc.

See `.goreleaser.yaml` `changelog` section for configuration.

## Tips & Tricks

### Undoing a Commit

If you committed something and need to fix the message:

```bash
git commit --amend
```

This opens an editor to change the commit message. Make sure it still follows conventions!

### Cherry-Pick from Branch

If you have commits on a feature branch:

```bash
git cherry-pick <commit-hash>
```

The commit message will be preserved, but verify it still follows conventions.

### Manual Validation

Validate without committing:

```bash
pre-commit run --all-files
```

Or run a specific hook:

```bash
pre-commit run commitlint --all-files
pre-commit run golangci-lint --all-files
```

## Troubleshooting

### Pre-commit not found

**Problem**: `pre-commit: command not found`

**Solution**:
```bash
pip install pre-commit
pre-commit install
pre-commit install --hook-type commit-msg
```

### Invalid commit message

**Problem**: Commit rejected with type/scope errors

**Solution**: Check the format:
- Type must be lowercase from the list above
- No period at end of subject
- Keep subject under 50 characters when possible
- Use imperative mood ("add" not "added" or "adds")

### Override hook (not recommended)

If you absolutely must bypass the hook:

```bash
git commit --no-verify
```

⚠️ **Warning**: This bypasses validation. Commits pushed will fail CI checks.

## References

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [commitlint docs](https://commitlint.js.org/)
- [husky docs](https://typicode.github.io/husky/)
