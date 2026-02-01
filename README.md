# substack-cli

A command-line tool for managing Substack publications. Create posts from markdown files, manage drafts, publish, and handle multiple accounts.

## Features

- **Markdown to Substack**: Write posts in standard markdown, publish directly from the terminal
- **Multi-account support**: Store and switch between multiple Substack accounts
- **Full post lifecycle**: Create drafts, publish, unpublish, update metadata
- **Draft management**: List, inspect, delete, and publish drafts
- **Configurable defaults**: Set default audience and email preferences

## Install

Requires Go 1.21+.

```sh
go install github.com/aaronsrivastava/substack-cli@latest
```

Or build from source:

```sh
git clone https://github.com/aaronsrivastava/substack-cli.git
cd substack-cli
go build -o substack .
# Move to somewhere on your PATH
mv substack /usr/local/bin/
```

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
