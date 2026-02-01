# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Frontmatter parsing for post metadata (title, subtitle, audience, tags, etc.)
- CLI arguments override frontmatter values
- Support for YAML frontmatter in markdown files
- Multi-platform distribution (Homebrew, APT, Windows)
- GoReleaser configuration for automated releases

### Fixed
- Fixed `draft_bylines` API requirement
- Fixed `draft_body` JSON stringification for API
- Fixed tight list items in markdown conversion
- Corrected field names to match Substack API (`draft_title`, `draft_subtitle`)

## [0.1.0] - 2026-02-01

### Added
- Initial release
- Create drafts from markdown files
- List, get, delete, and publish drafts
- Create and manage posts
- Multi-account support
- Full markdown to ProseMirror conversion
- Support for: headings, bold, italic, code, links, lists, blockquotes, horizontal rules, code blocks
- Draft management commands
- Post management commands
- Account authentication and management
- Configuration management

[Unreleased]: https://github.com/aaronsrivastava/substack-cli/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/aaronsrivastava/substack-cli/releases/tag/v0.1.0
