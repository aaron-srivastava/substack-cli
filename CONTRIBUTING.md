# Contributing to substack-cli

Thank you for your interest in contributing to substack-cli! This document provides guidelines for contributing to the project.

## Conventional Commits

This project uses [Conventional Commits](https://www.conventionalcommits.org/) to automatically determine version bumps and generate changelogs.

### Commit Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- **feat**: A new feature
- **fix**: A bug fix
- **perf**: A performance improvement
- **refactor**: Code refactoring without feature changes
- **docs**: Documentation changes
- **test**: Adding or updating tests
- **ci**: CI/CD configuration changes
- **chore**: Other changes (e.g., dependency updates)

### Examples

```bash
# New feature
git commit -m "feat: Add support for scheduled posts"

# Bug fix
git commit -m "fix: Resolve markdown list rendering issue"

# Breaking change (MAJOR version bump)
git commit -m "feat!: Change authentication method to OAuth"
git commit -m "feat(api): Update API client\n\nBREAKING CHANGE: Old auth method no longer supported"

# With scope
git commit -m "feat(markdown): Add support for tables"

# Fix with body
git commit -m "fix: Resolve draft body JSON encoding

The draft_body field was being double-encoded.
This fix ensures the JSON is properly formatted."
```

### Version Bumping

Commits are analyzed to determine version bumps:

- **BREAKING CHANGE** or `!` suffix → MAJOR version (1.0.0 → 2.0.0)
- **feat:** commits → MINOR version (1.0.0 → 1.1.0)
- **fix:, perf:, refactor:** → PATCH version (1.0.0 → 1.0.1)

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes with conventional commit messages
4. Run tests: `go test ./...`
5. Run linter: `golangci-lint run`
6. Push your branch and create a Pull Request

## Running Tests Locally

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run linter
golangci-lint run ./...

# Build locally
go build -o substack .
```

## Signing Commits

We recommend signing your commits with GPG. To sign a commit:

```bash
git commit -S -m "feat: Add new feature"
```

Or configure signing by default:

```bash
git config --global commit.gpgsign true
```

## Code Style

- Follow Go conventions (go fmt, go vet)
- Write clear, self-documenting code
- Add comments for complex logic
- Keep functions focused and testable

## Documentation

- Update README.md if adding new commands or features
- Update CHANGELOG.md in the Unreleased section
- Add docstrings to exported functions and types

## Release Process

Releases are automated. Once your PR is merged:

1. Maintainers update CHANGELOG.md
2. Create a git tag: `git tag -a v1.1.0 -m "Release v1.1.0: Description"`
3. Push the tag: `git push origin v1.1.0`
4. GitHub Actions automatically:
   - Runs GoReleaser
   - Builds all binaries
   - Creates GitHub release
   - Updates Homebrew formula
   - Announces on Bluesky

See [RELEASING.md](RELEASING.md) for detailed release instructions.

## Questions?

Feel free to open an issue or discussion if you have questions about contributing.

Thank you for helping make substack-cli better!
