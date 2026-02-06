package markdown

import (
	"bytes"
	"strings"

	"github.com/aaronsrivastava/substack-cli/internal/model"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

// Frontmatter holds parsed YAML frontmatter fields.
type Frontmatter struct {
	Title           string
	Subtitle        string
	Date            string
	Tags            []string
	Audience        string
	Draft           bool
	Slug            string
	CanonicalURL    string
	MetaDescription string
	SocialImage     string
	ScheduledAt     string
	Section         string
	PodcastURL      string
}

// ParseFrontmatter extracts YAML frontmatter and returns it along with the
// remaining body. If no frontmatter is present, returns nil and the full source.
func ParseFrontmatter(source []byte) (*Frontmatter, []byte) {
	s := string(source)
	if !strings.HasPrefix(s, "---\n") {
		return nil, source
	}
	end := strings.Index(s[4:], "\n---")
	if end < 0 {
		return nil, source
	}
	yamlBlock := s[4 : 4+end]
	body := []byte(s[4+end+4:]) // skip past closing "---\n"

	fm := &Frontmatter{}
	for _, line := range strings.Split(yamlBlock, "\n") {
		key, val, ok := parseYAMLLine(line)
		if !ok {
			continue
		}
		switch key {
		case "title":
			fm.Title = val
		case "subtitle":
			fm.Subtitle = val
		case "date":
			fm.Date = val
		case "tags":
			fm.Tags = parseYAMLList(val)
		case "audience":
			fm.Audience = val
		case "draft":
			fm.Draft = val == "true"
		case "slug":
			fm.Slug = val
		case "canonical_url":
			fm.CanonicalURL = val
		case "meta_description":
			fm.MetaDescription = val
		case "social_image":
			fm.SocialImage = val
		case "scheduled_at":
			fm.ScheduledAt = val
		case "section":
			fm.Section = val
		case "podcast_url":
			fm.PodcastURL = val
		}
	}
	return fm, body
}

func parseYAMLLine(line string) (string, string, bool) {
	idx := strings.Index(line, ":")
	if idx < 0 {
		return "", "", false
	}
	key := strings.TrimSpace(line[:idx])
	val := strings.TrimSpace(line[idx+1:])
	// Strip surrounding quotes
	val = strings.Trim(val, "\"'")
	return key, val, true
}

func parseYAMLList(val string) []string {
	// Handle [a, b, c] inline syntax
	val = strings.Trim(val, "[]")
	parts := strings.Split(val, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, "\"'")
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// Convert parses markdown source (with optional frontmatter) and returns (title, substackBody).
// The first H1 is extracted as the title (if present and frontmatter has no title).
func Convert(source []byte) (string, model.DraftBody) {
	_, body := ParseFrontmatter(source)
	return convertBody(body)
}

// ConvertWithFrontmatter parses markdown source, returning frontmatter, title, and body.
// Title priority: frontmatter title > first H1.
func ConvertWithFrontmatter(source []byte) (*Frontmatter, string, model.DraftBody) {
	fm, body := ParseFrontmatter(source)
	title, draftBody := convertBody(body)
	if fm != nil && fm.Title != "" {
		title = fm.Title
	}
	return fm, title, draftBody
}

func convertBody(source []byte) (string, model.DraftBody) {
	source = bytes.TrimLeft(source, "\n")
	md := goldmark.New(goldmark.WithExtensions(extension.Strikethrough))
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	var title string
	var nodes []model.Node

	for child := doc.FirstChild(); child != nil; child = child.NextSibling() {
		if h, ok := child.(*ast.Heading); ok && h.Level == 1 && title == "" {
			title = string(nodeText(child, source))
			continue
		}
		n := convertBlock(child, source)
		if n != nil {
			nodes = append(nodes, *n)
		}
	}

	return title, model.DraftBody{Type: "doc", Content: nodes}
}

func convertBlock(node ast.Node, source []byte) *model.Node {
	switch n := node.(type) {
	case *ast.Heading:
		return &model.Node{
			Type:    "heading",
			Attrs:   map[string]any{"level": n.Level},
			Content: convertInlineChildren(n, source),
		}
	case *ast.Paragraph:
		content := convertInlineChildren(n, source)
		return &model.Node{Type: "paragraph", Content: content}
	case *ast.Blockquote:
		var children []model.Node
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if cn := convertBlock(c, source); cn != nil {
				children = append(children, *cn)
			}
		}
		return &model.Node{Type: "blockquote", Content: children}
	case *ast.FencedCodeBlock:
		lang := string(n.Language(source))
		attrs := map[string]any{}
		if lang != "" {
			attrs["language"] = lang
		}
		code := codeBlockText(n, source)
		return &model.Node{
			Type:    "code_block",
			Attrs:   attrs,
			Content: []model.Node{{Type: "text", Text: code}},
		}
	case *ast.List:
		typ := "bullet_list"
		if n.IsOrdered() {
			typ = "ordered_list"
		}
		var items []model.Node
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if item := convertListItem(c, source); item != nil {
				items = append(items, *item)
			}
		}
		return &model.Node{Type: typ, Content: items}
	case *ast.ThematicBreak:
		return &model.Node{Type: "horizontal_rule"}
	default:
		// Handle TextBlock and other container nodes (e.g. tight list items)
		content := convertInlineChildren(node, source)
		if len(content) > 0 {
			return &model.Node{Type: "paragraph", Content: content}
		}
		return nil
	}
}

func convertListItem(node ast.Node, source []byte) *model.Node {
	var children []model.Node
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		if cn := convertBlock(c, source); cn != nil {
			children = append(children, *cn)
		}
	}
	return &model.Node{Type: "list_item", Content: children}
}

func convertInlineChildren(node ast.Node, source []byte) []model.Node {
	var nodes []model.Node
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		nodes = append(nodes, convertInline(c, source, nil)...)
	}
	return nodes
}

func convertInline(node ast.Node, source []byte, marks []model.Mark) []model.Node {
	switch n := node.(type) {
	case *ast.Text:
		t := string(n.Value(source))
		result := []model.Node{{Type: "text", Text: t, Marks: marks}}
		if n.SoftLineBreak() {
			result[0].Text += "\n"
		}
		return result
	case *ast.CodeSpan:
		newMarks := appendMark(marks, model.Mark{Type: "code"})
		t := string(nodeText(n, source))
		return []model.Node{{Type: "text", Text: t, Marks: newMarks}}
	case *ast.Emphasis:
		typ := "em"
		if n.Level == 2 {
			typ = "strong"
		}
		newMarks := appendMark(marks, model.Mark{Type: typ})
		var nodes []model.Node
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			nodes = append(nodes, convertInline(c, source, newMarks)...)
		}
		return nodes
	case *east.Strikethrough:
		newMarks := appendMark(marks, model.Mark{Type: "strikethrough"})
		var nodes []model.Node
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			nodes = append(nodes, convertInline(c, source, newMarks)...)
		}
		return nodes
	case *ast.Link:
		newMarks := appendMark(marks, model.Mark{
			Type:  "link",
			Attrs: map[string]any{"href": string(n.Destination)},
		})
		var nodes []model.Node
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			nodes = append(nodes, convertInline(c, source, newMarks)...)
		}
		return nodes
	case *ast.AutoLink:
		url := string(n.URL(source))
		newMarks := appendMark(marks, model.Mark{
			Type:  "link",
			Attrs: map[string]any{"href": url},
		})
		return []model.Node{{Type: "text", Text: url, Marks: newMarks}}
	default:
		// Recurse into unknown inline containers
		var nodes []model.Node
		for c := node.FirstChild(); c != nil; c = c.NextSibling() {
			nodes = append(nodes, convertInline(c, source, marks)...)
		}
		return nodes
	}
}

func appendMark(existing []model.Mark, m model.Mark) []model.Mark {
	newMarks := make([]model.Mark, len(existing), len(existing)+1)
	copy(newMarks, existing)
	return append(newMarks, m)
}

func nodeText(n ast.Node, source []byte) []byte {
	var buf []byte
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			buf = append(buf, t.Value(source)...)
		} else {
			buf = append(buf, nodeText(c, source)...)
		}
	}
	return buf
}

func codeBlockText(n *ast.FencedCodeBlock, source []byte) string {
	var buf []byte
	lines := n.Lines()
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		buf = append(buf, seg.Value(source)...)
	}
	return string(buf)
}
