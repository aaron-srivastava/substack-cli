# Quick Release Guide

This is a quick reference for making releases. See `RELEASE.md` for detailed setup instructions.

## Prerequisites

- [ ] GitHub account and repository access
- [ ] GitHub token with `public_repo` scope (create at https://github.com/settings/tokens)
- [ ] GoReleaser installed (`go install github.com/goreleaser/goreleaser@latest`)
- [ ] Token exported as environment variable

```bash
export GITHUB_TOKEN=your_token_here
```

## Release Process

### 1. Update CHANGELOG.md

Move all items from `[Unreleased]` section to a new version section:

```markdown
## [1.0.0] - YYYY-MM-DD

### Added
### Changed
### Fixed
### Removed

### [Unreleased]
```

Commit the changes:
```bash
git add CHANGELOG.md
git commit -m "docs: Update changelog for v1.0.0"
```

### 2. Create Version Tag

```bash
git tag -a v1.0.0 -m "Release v1.0.0: Brief description of major features/fixes"
```

The tag message should summarize what's new. Examples:

```bash
# Minor version
git tag -a v1.1.0 -m "Release v1.1.0: Add feature X, improve Y"

# Patch version
git tag -a v1.0.1 -m "Release v1.0.1: Fix bug with Z"

# Major version
git tag -a v2.0.0 -m "Release v2.0.0: Breaking change - migrate from X to Y"
```

### 3. Push Tag

```bash
git push origin v1.0.0
```

**IMPORTANT**: The GitHub Actions workflow will automatically run when the tag is pushed. GoReleaser will:
- Build binaries for all platforms
- Create a GitHub release
- Update the Homebrew formula
- Commit the formula back to main

### 4. Verify Release

Visit: https://github.com/aaronsrivastava/substack-cli/releases

Check that:
- [ ] Release created with correct version
- [ ] All platform binaries included (Linux, macOS, Windows)
- [ ] Homebrew formula updated in `homebrew-tap/Formula/substack-cli.rb`
- [ ] Formula commit pushed to main branch

## Semantic Versioning

Format: `MAJOR.MINOR.PATCH`

- **MAJOR** (1.0.0 → 2.0.0): Breaking changes
- **MINOR** (1.0.0 → 1.1.0): New features (backward compatible)
- **PATCH** (1.0.0 → 1.0.1): Bug fixes only

## Distribution

### Homebrew (macOS/Linux)
Users can install with:
```bash
brew tap aaronsrivastava/substack-cli https://github.com/aaronsrivastava/substack-cli
brew install substack-cli
```

### APT (Debian/Ubuntu)
Download `.deb` from GitHub releases or use Releases page.

### Windows
Download `.zip` from GitHub releases and add to PATH, or use a package manager.

## Example Release Sequence

```bash
# 1. Make your changes and commit them
git commit -m "feat: Add feature X"

# 2. Update changelog
# Edit CHANGELOG.md, move [Unreleased] items to [1.1.0]
git commit -m "docs: Update changelog for v1.1.0"

# 3. Create and push tag
git tag -a v1.1.0 -m "Release v1.1.0: Add feature X"
git push origin v1.1.0

# 4. Watch GitHub Actions run the release
# GitHub Actions → Actions tab → Latest "Release" workflow

# 5. Verify at GitHub releases page
# https://github.com/aaronsrivastava/substack-cli/releases
```

## Troubleshooting

### "tag already exists"
```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin :refs/tags/v1.0.0

# Recreate and push
git tag -a v1.0.0 -m "..."
git push origin v1.0.0
```

### GoReleaser workflow failed
- Check GitHub Actions logs: https://github.com/aaronsrivastava/substack-cli/actions
- Verify GITHUB_TOKEN is set in repository secrets (if using manual GoReleaser)
- Check `.goreleaser.yaml` syntax

### Homebrew formula not updated
- Verify `homebrew-tap/Formula/` directory exists
- Check GoReleaser logs for errors
- Formula should be auto-committed to main branch after release

## Testing Before Release (Optional)

Do a dry-run without publishing:

```bash
goreleaser release --skip-publish --clean
```

This creates artifacts in `dist/` without uploading to GitHub. Check that binaries work:

```bash
./dist/substack-cli_linux_amd64_v1/substack --version
./dist/substack-cli_darwin_amd64_v1/substack --version
./dist/substack-cli_windows_amd64_v1/substack.exe --version
```

## Release Checklist

- [ ] All code changes committed
- [ ] CHANGELOG.md updated
- [ ] Version tag created (`git tag -a vX.Y.Z ...`)
- [ ] Tag pushed to origin (`git push origin vX.Y.Z`)
- [ ] GitHub Actions workflow completes successfully
- [ ] Release appears on GitHub releases page
- [ ] Binaries available for all platforms
- [ ] Homebrew formula updated in `homebrew-tap/Formula/`
- [ ] Users can install with `brew install` and `apt`
