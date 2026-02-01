# Distribution & Release Setup

This document summarizes the complete distribution and release infrastructure for substack-cli.

## 🔐 Security: octo-sts Integration

The release workflow uses **octo-sts** (Octo Secure Token Service) for secure GitHub token handling:

- **No long-lived secrets**: Short-lived tokens obtained via OIDC
- **Automatic expiration**: Tokens expire after 1 hour
- **Fine-grained permissions**: Only `contents: write` and `id-token: write` needed
- **Audit trail**: All token requests logged in GitHub Actions

### How it Works

1. GitHub Actions runs the release workflow
2. octo-sts action requests a short-lived token using OIDC
3. Token is passed to GoReleaser
4. Token automatically expires after use

**No manual token management needed!**

## 📝 Conventional Commits

All commits follow [Conventional Commits](https://www.conventionalcommits.org/) format to automatically determine version bumps:

```
<type>[optional scope]: <description>
```

### Supported Types

- **feat**: Features (→ MINOR version bump)
- **fix**: Bug fixes (→ PATCH version bump)
- **perf**: Performance improvements
- **refactor**: Code refactoring
- **docs**: Documentation
- **test**: Tests
- **ci**: CI/CD
- **chore**: Other changes

### Version Bumping

GoReleaser automatically determines the next version based on commits:

```
1.0.0 + feat commit   → 1.1.0 (MINOR)
1.0.0 + fix commit    → 1.0.1 (PATCH)
1.0.0 + BREAKING      → 2.0.0 (MAJOR)
```

### Example Commits

```bash
git commit -m "feat: Add scheduled post support"
git commit -m "fix: Resolve markdown list issue"
git commit -m "feat!: Change authentication method"  # Breaking change
git commit -m "perf: Optimize API calls"
git commit -m "docs: Update readme"
```

## 📢 Bluesky Announcements

Releases are automatically announced on Bluesky (@aaronsrivastava.com).

### Setup (One-time)

1. Create Bluesky account and app password:
   - Visit https://bsky.app/settings/app-passwords
   - Create new app password
   - Save the password

2. Add to GitHub repository secrets:
   - Go to repository Settings → Secrets and variables → Actions
   - Create new secret: `BLUESKY_APP_PASSWORD`
   - Paste the app password

### What Gets Posted

Each release announcement includes:
- Version number
- Release notes from changelog
- Direct link to release on GitHub

Example:
```
🚀 substack-cli v1.2.0 is out!

Features:
- Add scheduled post support
- Improve markdown rendering

Download: https://github.com/aaronsrivastava/substack-cli/releases/tag/v1.2.0
```

## 📦 Distribution Channels

### Homebrew (macOS/Linux)

Users install with:
```bash
brew tap aaronsrivastava/substack-cli https://github.com/aaronsrivastava/substack-cli
brew install substack-cli
```

Formula location: `homebrew-tap/Formula/substack-cli.rb` (auto-generated)

### APT/Debian/Ubuntu

Download `.deb` from GitHub releases:
```bash
sudo dpkg -i substack-cli_*.deb
```

Or manually from: https://github.com/aaronsrivastava/substack-cli/releases

### Windows

Download `.zip` from releases and extract to PATH, or use a package manager.

### Direct Download

All releases: https://github.com/aaronsrivastava/substack-cli/releases

Includes:
- Linux: amd64, arm64
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Windows: amd64, arm64

## 🚀 Release Workflow

### Manual Release Process

```bash
# 1. Make and commit your changes using conventional commits
git commit -m "feat: New feature"
git commit -m "fix: Bug fix"

# 2. Update CHANGELOG.md manually (move items from Unreleased)
git add CHANGELOG.md
git commit -m "docs: Update changelog for v1.1.0"

# 3. Create a git tag
git tag -a v1.1.0 -m "Release v1.1.0: New features and fixes"

# 4. Push the tag (this triggers GitHub Actions)
git push origin v1.1.0
```

### What Happens Automatically

When you push a version tag, GitHub Actions:

1. **Obtains token** via octo-sts (no manual secrets needed!)
2. **Runs GoReleaser** which:
   - Builds binaries for all platforms
   - Generates changelog from conventional commits
   - Creates GitHub release
   - Updates Homebrew formula in `homebrew-tap/Formula/`
   - Commits formula back to main
3. **Announces on Bluesky** (if BLUESKY_APP_PASSWORD is set)

### Monitoring

Watch the release progress:
- GitHub: https://github.com/aaronsrivastava/substack-cli/actions
- Check "Release" workflow for the tagged version
- View logs to troubleshoot any issues

## 🔧 Configuration Files

### `.goreleaser.yaml`

Main release configuration:
- **builds**: Target platforms (Linux, macOS, Windows)
- **git**: Conventional commit parsing
- **changelog**: Grouped by commit type
- **announce.bluesky**: Bluesky announcer settings
- **brews**: Homebrew formula generation
- **nfpms**: Debian/APT packages

### `.github/workflows/release.yml`

GitHub Actions workflow:
- **trigger**: Runs on version tags (v*)
- **octo-sts**: Gets short-lived token
- **goreleaser**: Builds and releases
- **environment**: BLUESKY_APP_PASSWORD secret

### `CONTRIBUTING.md`

Contributor guidelines:
- Conventional commit format
- How version bumping works
- PR process
- Code style

## 📋 Checklist for First Release

- [ ] Update CHANGELOG.md (move Unreleased → v1.0.0)
- [ ] Commit with conventional message: `docs: Update changelog for v1.0.0`
- [ ] Create tag: `git tag -a v1.0.0 -m "Release v1.0.0: Description"`
- [ ] Push tag: `git push origin v1.0.0`
- [ ] (Optional) Add BLUESKY_APP_PASSWORD secret to GitHub
- [ ] Monitor GitHub Actions workflow
- [ ] Verify release on GitHub releases page
- [ ] Test Homebrew installation: `brew install substack-cli` (after tap setup)

## 📚 Related Documentation

- [RELEASING.md](RELEASING.md) - Quick reference for making releases
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contributor guidelines
- [CHANGELOG.md](CHANGELOG.md) - Version history

## 🔗 External Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [GoReleaser Documentation](https://goreleaser.com/)
- [octo-sts Action](https://github.com/octo-sts/action)
- [Bluesky API Docs](https://docs.bsky.app/)

## Frequently Asked Questions

### Do I need to set GITHUB_TOKEN manually?

No! octo-sts handles it automatically. Just make sure the repository has Actions enabled.

### How do I announce to Bluesky?

1. Create/get your Bluesky app password from https://bsky.app/settings/app-passwords
2. Add it to GitHub repository secrets as `BLUESKY_APP_PASSWORD`
3. Releases will automatically announce to @aaronsrivastava.com

### Can I skip the Bluesky announcement?

Yes, just don't set the `BLUESKY_APP_PASSWORD` secret. GoReleaser will skip it silently.

### What if a release fails?

Check the GitHub Actions logs: https://github.com/aaronsrivastava/substack-cli/actions

Common issues:
- Missing `BLUESKY_APP_PASSWORD` → announcement skipped (not an error)
- Invalid git tag → check tag format (v*.*)
- Conventional commit parsing → ensure commits follow format

### Can I manually trigger a release?

Currently, releases are triggered by pushing tags. To test locally:

```bash
export GITHUB_TOKEN=your_token
goreleaser release --skip-publish --clean
```

### How do I fix a broken release?

1. Delete the tag locally: `git tag -d v1.0.0`
2. Delete remote tag: `git push origin :refs/tags/v1.0.0`
3. Fix the issue
4. Recreate tag and push
