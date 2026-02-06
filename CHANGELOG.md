# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.1] - 2026-02-06

### Fixed
- Fixed test failures from lint feedback

### CI/CD
- Switched to go-semantic-release for automated releases
- Fixed GITHUB_TOKEN environment variable in auto-release workflow
- Fixed repository owner name in goreleaser config

## [0.2.0] - 2026-02-02

### Added
- Octo-STS integration for secure token management
- Conventional commits enforcement
- Auto-release workflow triggered on push to main
- Bluesky announcements for new releases
- Pre-commit framework for code quality

### Documentation
- Added comprehensive distribution setup guide
- Added quick release reference guide
- Added development setup instructions

### Refactoring
- Replaced husky with pre-commit framework

## [0.1.0] - 2026-02-01

### Added
- Initial release
- Create drafts from markdown files with frontmatter support
- Frontmatter parsing for post metadata (title, subtitle, audience, tags, etc.)
- CLI arguments override frontmatter values
- List, get, delete, and publish drafts
- Create and manage posts
- Multi-account support
- Full markdown to ProseMirror conversion
- Support for: headings, bold, italic, code, links, lists, blockquotes, horizontal rules, code blocks
- Multi-platform distribution via Homebrew, APT, and Windows
- GoReleaser configuration for automated builds

[0.2.1]: https://github.com/aaron-srivastava/substack-cli/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/aaron-srivastava/substack-cli/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/aaron-srivastava/substack-cli/releases/tag/v0.1.0
