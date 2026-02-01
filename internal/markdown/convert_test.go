package markdown

import (
	"testing"
)

func TestConvert_TitleExtraction(t *testing.T) {
	src := []byte("# My Title\n\nSome text.\n")
	title, body := Convert(src)
	if title != "My Title" {
		t.Errorf("title = %q, want %q", title, "My Title")
	}
	if body.Type != "doc" {
		t.Errorf("body type = %q, want %q", body.Type, "doc")
	}
	if len(body.Content) != 1 {
		t.Fatalf("expected 1 node, got %d", len(body.Content))
	}
	if body.Content[0].Type != "paragraph" {
		t.Errorf("node type = %q, want %q", body.Content[0].Type, "paragraph")
	}
}

func TestConvert_NoTitle(t *testing.T) {
	src := []byte("Just a paragraph.\n")
	title, body := Convert(src)
	if title != "" {
		t.Errorf("title = %q, want empty", title)
	}
	if len(body.Content) != 1 || body.Content[0].Type != "paragraph" {
		t.Errorf("unexpected body: %+v", body.Content)
	}
}

func TestConvert_Heading(t *testing.T) {
	src := []byte("## Sub heading\n")
	_, body := Convert(src)
	if len(body.Content) != 1 {
		t.Fatalf("expected 1 node, got %d", len(body.Content))
	}
	n := body.Content[0]
	if n.Type != "heading" {
		t.Errorf("type = %q, want heading", n.Type)
	}
	if n.Attrs["level"] != 2 {
		t.Errorf("level = %v, want 2", n.Attrs["level"])
	}
}

func TestConvert_Bold(t *testing.T) {
	src := []byte("**bold**\n")
	_, body := Convert(src)
	text := body.Content[0].Content[0]
	if text.Text != "bold" {
		t.Errorf("text = %q", text.Text)
	}
	if len(text.Marks) != 1 || text.Marks[0].Type != "strong" {
		t.Errorf("marks = %+v, want strong", text.Marks)
	}
}

func TestConvert_Italic(t *testing.T) {
	src := []byte("*italic*\n")
	_, body := Convert(src)
	text := body.Content[0].Content[0]
	if len(text.Marks) != 1 || text.Marks[0].Type != "em" {
		t.Errorf("marks = %+v, want em", text.Marks)
	}
}

func TestConvert_Strikethrough(t *testing.T) {
	src := []byte("~~strike~~\n")
	_, body := Convert(src)
	text := body.Content[0].Content[0]
	if len(text.Marks) != 1 || text.Marks[0].Type != "strikethrough" {
		t.Errorf("marks = %+v, want strikethrough", text.Marks)
	}
}

func TestConvert_InlineCode(t *testing.T) {
	src := []byte("`code`\n")
	_, body := Convert(src)
	text := body.Content[0].Content[0]
	if len(text.Marks) != 1 || text.Marks[0].Type != "code" {
		t.Errorf("marks = %+v, want code", text.Marks)
	}
}

func TestConvert_Link(t *testing.T) {
	src := []byte("[click](https://example.com)\n")
	_, body := Convert(src)
	text := body.Content[0].Content[0]
	if len(text.Marks) != 1 || text.Marks[0].Type != "link" {
		t.Fatalf("marks = %+v, want link", text.Marks)
	}
	if text.Marks[0].Attrs["href"] != "https://example.com" {
		t.Errorf("href = %v", text.Marks[0].Attrs["href"])
	}
}

func TestConvert_Blockquote(t *testing.T) {
	src := []byte("> quoted\n")
	_, body := Convert(src)
	if body.Content[0].Type != "blockquote" {
		t.Errorf("type = %q, want blockquote", body.Content[0].Type)
	}
}

func TestConvert_BulletList(t *testing.T) {
	src := []byte("- a\n- b\n")
	_, body := Convert(src)
	if body.Content[0].Type != "bullet_list" {
		t.Errorf("type = %q, want bullet_list", body.Content[0].Type)
	}
	if len(body.Content[0].Content) != 2 {
		t.Errorf("items = %d, want 2", len(body.Content[0].Content))
	}
}

func TestConvert_OrderedList(t *testing.T) {
	src := []byte("1. a\n2. b\n")
	_, body := Convert(src)
	if body.Content[0].Type != "ordered_list" {
		t.Errorf("type = %q, want ordered_list", body.Content[0].Type)
	}
}

func TestConvert_CodeBlock(t *testing.T) {
	src := []byte("```python\nprint('hi')\n```\n")
	_, body := Convert(src)
	n := body.Content[0]
	if n.Type != "code_block" {
		t.Errorf("type = %q, want code_block", n.Type)
	}
	if n.Attrs["language"] != "python" {
		t.Errorf("lang = %v, want python", n.Attrs["language"])
	}
}

func TestConvert_HorizontalRule(t *testing.T) {
	src := []byte("---\n")
	_, body := Convert(src)
	if body.Content[0].Type != "horizontal_rule" {
		t.Errorf("type = %q, want horizontal_rule", body.Content[0].Type)
	}
}
