# substack-cli

A command-line tool for managing Substack publications. Create posts from markdown files, manage drafts, publish, and handle multiple accounts.

## Features

- **Markdown to Substack**: Write posts in standard markdown, publish directly from the terminal
- **Multi-account support**: Store and switch between multiple Substack accounts
- **Full post lifecycle**: Create drafts, publish, unpublish, update metadata
- **Draft management**: List, inspect, delete, and publish drafts
- **Configurable defaults**: Set default audience and email preferences

## Install

### Via Homebrew (macOS/Linux)

```sh
brew tap aaron-srivastava/substack-cli https://github.com/aaron-srivastava/substack-cli
brew install substack-cli
```

### Via APT (Debian/Ubuntu)

Download from [releases](https://github.com/aaron-srivastava/substack-cli/releases) page.

### From Go

Requires Go 1.25.6+.

```sh
go install github.com/aaron-srivastava/substack-cli@latest
```

### Build from Source

```sh
git clone https://github.com/aaron-srivastava/substack-cli.git
cd substack-cli
go build -o substack .
mv substack /usr/local/bin/
```

See [SETUP.md](SETUP.md) for development setup with commit hooks and automation.

## Getting Started

### 1. Get your Substack cookies

Log into your Substack in a browser, then open Developer Tools > Application > Cookies for `.substack.com`. You need these values:

| Cookie | Description |
|--------|-------------|
| `connect.sid` | Session ID |
| `substack.sid` | Substack session ID |
| `substack.lli` | Substack login indicator |

You also need your **User ID** (found in network requests or your Substack profile API) and your **publication URL** (e.g. `https://you.substack.com`).

### 2. Log in

```sh
substack auth login
```

You'll be prompted for:
- Account name (a label you choose, e.g. `my-blog`)
- Publication URL
- User ID
- SID (`connect.sid` cookie)
- `substack.sid` cookie
- `substack.lli` cookie

Credentials are stored in `~/.config/substack-cli/accounts.json` with restricted file permissions.

### 3. Create a post

Write a markdown file:

```md
# My Post Title

This is a **bold** statement with a [link](https://example.com).

## Section

- Point one
- Point two

> A blockquote for emphasis.
```

Create a draft:

```sh
substack post create my-post.md
```

Create and publish immediately:

```sh
substack post create my-post.md --publish --send-email --audience everyone
```

The first `# H1` heading is used as the post title. Override it with `--title`.

### 4. Manage drafts and posts

```sh
substack draft list
substack draft publish 12345 --send-email
substack post list
substack post get 12345
substack post unpublish 12345
```

## Commands

```
substack auth login              Add or update an account
substack auth status             Show active account
substack auth list               List all accounts
substack auth switch <name>      Switch active account
substack auth remove <name>      Remove an account

substack post create <file.md>   Create post from markdown
  --title <T>                    Override title
  --subtitle <S>                 Set subtitle
  --publish                      Publish immediately
  --send-email                   Email subscribers
  --audience <A>                 "everyone" or "only_paid"
substack post list               List published posts
substack post get <id>           Show post details
substack post unpublish <id>     Unpublish a post
substack post update <id>        Update metadata (--title, --subtitle, --audience)

substack draft list              List drafts
substack draft get <id>          Show draft details
substack draft delete <id>       Delete a draft
substack draft publish <id>      Publish a draft (--send-email, --audience)

substack config show             Show default settings
substack config set <key> <val>  Set defaults (send_email, audience)
```

## Supported Markdown

| Markdown | Substack element |
|----------|-----------------|
| `# H1` through `###### H6` | Headings |
| `**bold**` | Bold text |
| `*italic*` | Italic text |
| `~~strikethrough~~` | Strikethrough |
| `` `code` `` | Inline code |
| `[text](url)` | Links |
| `> quote` | Blockquote |
| `- item` | Bullet list |
| `1. item` | Ordered list |
| ` ```lang ` | Code block with syntax |
| `---` | Horizontal rule |

## Development

### Setup

See [SETUP.md](SETUP.md) for complete development setup guide.

Quick start:

```sh
git clone https://github.com/aaron-srivastava/substack-cli.git
cd substack-cli
pip install pre-commit          # Install pre-commit framework
pre-commit install              # Install git hooks
go mod download                 # Download Go dependencies
go test ./...                   # Run tests
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make changes with [conventional commits](COMMIT_CONVENTIONS.md)
4. Push and create a pull request

All commits must follow [Conventional Commits](https://www.conventionalcommits.org/). A pre-commit hook will validate your messages automatically.

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

### Releases

Releases are automated! When you merge to `main`:

1. Commits are analyzed
2. Version is automatically bumped (major/minor/patch)
3. CHANGELOG.md is updated
4. Release is created with binaries
5. Announcement is posted on Bluesky

See [AUTO_RELEASE.md](AUTO_RELEASE.md) for details.

### Documentation

- [SETUP.md](SETUP.md) - Development setup
- [COMMIT_CONVENTIONS.md](COMMIT_CONVENTIONS.md) - Commit message format
- [AUTO_RELEASE.md](AUTO_RELEASE.md) - Automatic releases
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contributor guidelines
- [RELEASING.md](RELEASING.md) - Manual release process
- [DISTRIBUTION_SETUP.md](DISTRIBUTION_SETUP.md) - Distribution details
