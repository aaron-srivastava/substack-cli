package model

import "time"

type Account struct {
	Name           string `json:"name"`
	PublicationURL string `json:"publication_url"`
	UserID         string `json:"user_id"`
	SID            string `json:"sid"`
	SubstackSID    string `json:"substack_sid"`
	SubstackLLI    string `json:"substack_lli"`
}

type AccountStore struct {
	Active   string    `json:"active"`
	Accounts []Account `json:"accounts"`
}

type PublishOptions struct {
	SendEmail bool   `json:"send_to_everyone"`
	Audience  string `json:"audience"` // "everyone" or "only_paid"
}

type DraftBody struct {
	Type    string `json:"type"`
	Content []Node `json:"content,omitempty"`
}

type Node struct {
	Type    string         `json:"type"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Content []Node         `json:"content,omitempty"`
	Marks   []Mark         `json:"marks,omitempty"`
	Text    string         `json:"text,omitempty"`
}

type Mark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

type Byline struct {
	ID int `json:"id"`
}

type DraftRequest struct {
	Title         string   `json:"draft_title,omitempty"`
	Subtitle      string   `json:"draft_subtitle,omitempty"`
	DraftBody     string   `json:"draft_body"`
	DraftBylines  []Byline `json:"draft_bylines"`
	Audience      string   `json:"audience,omitempty"`
	Section       string   `json:"draft_section_id,omitempty"`
	SectionChosen bool     `json:"section_chosen"`
	Type          string   `json:"type"`
}

type Post struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Subtitle     string    `json:"subtitle"`
	Slug         string    `json:"slug"`
	Audience     string    `json:"audience"`
	Type         string    `json:"type"`
	DraftCreated time.Time `json:"draft_created_at"`
	PostDate     string    `json:"post_date"`
	IsPublished  bool      `json:"is_published"`
	WordCount    int       `json:"word_count"`
}

type DraftResponse struct {
	ID           int       `json:"id"`
	Title        string    `json:"draft_title"`
	Subtitle     string    `json:"draft_subtitle"`
	Slug         string    `json:"slug"`
	Audience     string    `json:"audience"`
	DraftCreated time.Time `json:"draft_created_at"`
	WordCount    int       `json:"word_count"`
}

type Config struct {
	SendEmail    bool   `json:"send_email"`
	Audience     string `json:"audience"`
	Section      string `json:"section"`
	OutputFormat string `json:"output_format"`
}
