# Quick Release Guide

This is a quick reference for making releases. See `RELEASE.md` for detailed setup instructions.

## Prerequisites

- [ ] GitHub account and repository access
- [ ] octo-sts configured in GitHub Actions (uses OIDC, no manual token needed)
- [ ] GoReleaser installed locally for testing (`go install github.com/goreleaser/goreleaser@latest`)
- [ ] Bluesky account for announcements (optional)
- [ ] Bluesky app password set in GitHub secrets as `BLUESKY_APP_PASSWORD` (optional)

### Bluesky Setup (Optional)

To announce releases on Bluesky:
1. Go to https://bsky.app/settings/app-passwords
2. Create a new app password
3. Add it to GitHub repository secrets as `BLUESKY_APP_PASSWORD`

Announcements will automatically post to @aaronsrivastava.com during releases.

## Release Process

### 0. Understand Conventional Commits

This project uses [Conventional Commits](https://www.conventionalcommits.org/) to automatically determine version bumps:

- **feat:** → MINOR version bump (1.0.0 → 1.1.0)
- **fix:** → PATCH version bump (1.0.0 → 1.0.1)
- **BREAKING CHANGE** or `!` → MAJOR version bump (1.0.0 → 2.0.0)

See [CONTRIBUTING.md](CONTRIBUTING.md) for commit format details.

### 1. Update CHANGELOG.md

Move all items from `[Unreleased]` section to a new version section:

```markdown
## [1.0.0] - YYYY-MM-DD

### Features
### Bug fixes
### Performance
### Refactoring

## [Unreleased]
```

Commit the changes using conventional commits:
```bash
git add CHANGELOG.md
git commit -m "docs: Update changelog for v1.0.0"
```

### 2. Create Version Tag

The version is determined by analyzing commit messages. Create an annotated tag:

```bash
git tag -a v1.0.0 -m "Release v1.0.0: Brief description of major features/fixes"
```

Examples based on conventional commits:

```bash
# Patch release (fix commits)
git tag -a v1.0.1 -m "Release v1.0.1: Fix bug with Z"

# Minor release (feat commits)
git tag -a v1.1.0 -m "Release v1.1.0: Add feature X"

# Major release (breaking changes)
git tag -a v2.0.0 -m "Release v2.0.0: Breaking change"
```

### 3. Push Tag

```bash
git push origin v1.0.0
```

**IMPORTANT**: The GitHub Actions workflow automatically triggers when the tag is pushed. Here's what happens:

1. **octo-sts** obtains a short-lived GitHub token using OIDC
2. **GoReleaser** analyzes commits and:
   - Builds binaries for all platforms
   - Creates GitHub release with auto-generated changelog
   - Updates Homebrew formula
   - Commits formula changes back to main
3. **Bluesky announcer** posts release announcement to @aaronsrivastava.com (if configured)

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
# First ensure you have a GitHub token for testing
export GITHUB_TOKEN=your_token_here

goreleaser release --skip-publish --clean
```

This creates artifacts in `dist/` without uploading to GitHub. Check that binaries work:

```bash
./dist/substack-cli_linux_amd64_v1/substack --version
./dist/substack-cli_darwin_amd64_v1/substack --version
./dist/substack-cli_windows_amd64_v1/substack.exe --version
```

### Testing Bluesky Announcements

To test announcements without releasing:

```bash
export GITHUB_TOKEN=your_token_here
export BLUESKY_APP_PASSWORD=your_app_password

goreleaser release --skip-publish --skip-sign --clean
```

Note: Bluesky announcements won't actually post during dry-runs.

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
