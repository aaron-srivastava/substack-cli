# Automatic Release Workflow

The auto-release workflow automatically creates version tags and releases when commits are pushed to the `main` branch. No manual tagging needed!

## How It Works

### Trigger

The workflow runs on every push to `main` (excluding documentation-only changes):

```yaml
on:
  push:
    branches:
      - main
    paths-ignore:
      - "README.md"
      - "CHANGELOG.md"
      - "docs/**"
      - ".gitignore"
```

### Process

1. **Analyze Commits** (.analyze-commits job)
   - Gets the latest tag (or v0.0.0 if none exists)
   - Examines all commits since the last release
   - Detects breaking changes, features, and fixes

2. **Determine Version Bump**
   - Breaking changes (`!` or `BREAKING CHANGE`) → MAJOR
   - Features (`feat:` commits) → MINOR
   - Fixes (`fix:` commits) → PATCH
   - Example: v1.2.3 + feat → v1.3.0

3. **Create Release** (.create-release job, only if version bumps)
   - Updates CHANGELOG.md
   - Commits changelog
   - Creates annotated git tag
   - Pushes tag to GitHub
   - Triggers GoReleaser workflow automatically

4. **Build & Announce** (release.yml workflow, triggered by tag)
   - GoReleaser builds all binaries
   - Creates GitHub release
   - Updates Homebrew formula
   - Announces on Bluesky

## Example Flow

### Step 1: Make commits with conventional format

```bash
git commit -m "feat: Add post scheduling"
git commit -m "fix: Handle markdown lists"
git push origin my-feature-branch
```

### Step 2: Create pull request & merge to main

Create PR and merge to main. The auto-release workflow runs automatically!

### Step 3: Auto-release detects changes

```
✓ Analyzing commits since v1.0.0
✓ Found: 1 feature, 1 fix
✓ Version bump: v1.0.0 → v1.1.0 (MINOR)
✓ Updating CHANGELOG.md
✓ Creating tag v1.1.0
✓ Pushing tag...
```

### Step 4: GoReleaser takes over

```
✓ Building binaries
✓ Creating GitHub release
✓ Updating Homebrew formula
✓ Announcing on Bluesky
```

## Version Bumping Rules

Based on [Semantic Versioning](https://semver.org/):

### MAJOR (x.0.0)

Breaking changes - use when upgrading requires user action:

```bash
feat!: Change authentication to OAuth
# or
feat: Update API

BREAKING CHANGE: Old session tokens no longer work
```

### MINOR (1.x.0)

New features - backward compatible:

```bash
feat: Add post scheduling
feat(cli): Add interactive mode
```

### PATCH (1.0.x)

Bug fixes:

```bash
fix: Resolve markdown rendering
fix(api): Handle null responses
perf: Optimize database queries
```

### NO AUTO-RELEASE

These commit types don't trigger releases:

```bash
docs: Update readme
test: Add unit tests
ci: Update GitHub Actions
chore: Update dependencies
```

## Manual Release Override

If you need to release without new code changes (e.g., documentation fixes):

1. Create a simple commit:
   ```bash
   git commit --allow-empty -m "chore: Release v1.2.0"
   git push origin main
   ```

2. Or manually create the tag:
   ```bash
   git tag -a v1.2.0 -m "Release v1.2.0: Documentation fixes"
   git push origin v1.2.0
   ```

## Monitoring

### Watch the Workflow

1. Go to https://github.com/aaronsrivastava/substack-cli/actions
2. Look for "Auto-Release" workflow
3. View logs for details

### Common Outputs

**No release needed:**
```
✓ Analyzing commits since v1.0.0
✓ No significant changes detected
✓ Release type: none
✓ Skipping release
```

**Release created:**
```
✓ Analyzing commits since v1.0.0
✓ Found: 2 features, 1 fix
✓ Version bump: v1.0.0 → v1.1.0
✓ Updated CHANGELOG.md
✓ Created tag v1.1.0
✓ Pushed tag - GoReleaser workflow triggered
```

## Changelog Updates

The workflow automatically updates `CHANGELOG.md` with:
- New version section (e.g., `[1.1.0] - 2024-02-01`)
- Empty subsections (Added, Changed, Fixed, Removed)
- Link to the version tag

**You should manually populate the subsections** with descriptions of changes.

Example:

```markdown
## [1.1.0] - 2024-02-01

### Added
- Post scheduling support
- Interactive CLI mode

### Fixed
- Markdown list rendering
- API timeout handling

## [1.0.0] - 2024-01-15
...
```

## Disable Auto-Release

To temporarily disable auto-releases, skip the workflow in your commit:

```bash
git commit -m "feat: Work in progress

[skip ci]"
```

Or delete the `.github/workflows/auto-release.yml` file.

## Frequently Asked Questions

### Why no release for my commit?

Check if:
1. Commit message starts with `feat:`, `fix:`, `perf:`, `refactor:`, etc.
2. Commit was on `main` branch (not a branch push)
3. Commit didn't only touch docs/tests/ci/gitignore

### Can I force a release?

Yes, create an empty commit:

```bash
git commit --allow-empty -m "chore: Bump version for release"
git push origin main
```

This will trigger the workflow (no significant changes detected though).

### What if the workflow fails?

1. Check the Actions tab for error logs
2. Common issues:
   - Git config not set (workflow sets this)
   - octo-sts token expired (usually not an issue)
   - Tag already exists (delete and retry)

### Can I manually edit the version?

The version is determined by commits. To force a specific version:

1. Manually create the tag:
   ```bash
   git tag -a v2.0.0 -m "Release v2.0.0"
   git push origin v2.0.0
   ```

2. The release.yml workflow will trigger and build

### Do I need to update CHANGELOG.md?

The workflow updates it automatically with the version and date. You should then manually add descriptions under the Added/Changed/Fixed sections.

## Configuration

Main configuration file: `.github/workflows/auto-release.yml`

Key settings:

- **paths-ignore**: Changes to these paths don't trigger releases
- **MAJOR/MINOR/PATCH**: Determined by commit analysis
- **commitlint**: Validates commit format in PRs

## Best Practices

1. **Use meaningful commit messages** following conventional commits
2. **Group related changes** in one feature/fix commit
3. **Squash fixup commits** before merging to main
4. **Don't merge work-in-progress** PRs to main
5. **Let the automation handle versioning** - don't manually tag unless necessary

## Examples

### Feature Release

```
Commits:
✓ feat: Add image upload support
✓ fix: Handle concurrent requests

Result: v1.2.0 → v1.3.0 (MINOR)
Changelog: "Added image upload support"
```

### Hotfix Release

```
Commits:
✓ fix: Critical security vulnerability

Result: v1.3.0 → v1.3.1 (PATCH)
Changelog: "Fixed critical security vulnerability"
```

### Major Release

```
Commits:
✓ feat!: Migrate to new API
✓ feat: Redesign CLI
✓ perf: Improve performance

Result: v1.3.1 → v2.0.0 (MAJOR)
Changelog: Multiple entries
```

## Troubleshooting

### Workflow shows "should-release: false"

**Cause**: No commits with feat/fix since last tag

**Solution**: Add a proper commit:
```bash
git commit -m "feat: Add new feature"
git push origin main
```

### Tag creation fails

**Cause**: Tag already exists locally/remotely

**Solution**:
```bash
git tag -d v1.2.0  # Local
git push origin :refs/tags/v1.2.0  # Remote
# Then retry
```

### CHANGELOG.md merge conflicts

If multiple auto-releases happen fast, resolve conflicts:
```bash
git checkout --ours CHANGELOG.md  # or --theirs
git add CHANGELOG.md
git commit -m "chore: Resolve changelog conflict"
```

## See Also

- [COMMIT_CONVENTIONS.md](COMMIT_CONVENTIONS.md) - How to write commits
- [RELEASING.md](RELEASING.md) - Manual release process
- [.github/workflows/auto-release.yml](.github/workflows/auto-release.yml) - Workflow config
