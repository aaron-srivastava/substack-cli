//go:build integration

package api

import (
	"fmt"
	"testing"

	"github.com/aaronsrivastava/substack-cli/internal/model"
)

// Run with: go test -tags=integration -v ./internal/api/...
//
// These tests create real drafts on your Substack. You must:
// 1. Be logged in via `substack-cli login`
// 2. Manually delete the drafts after reviewing

func TestIntegration_CreateDraftWithAllAudiences(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v (are you logged in?)", err)
	}

	audiences := []string{"everyone", "only_paid", "only_free"}

	for _, aud := range audiences {
		t.Run(aud, func(t *testing.T) {
			draft := model.DraftRequest{
				Title:     fmt.Sprintf("Test Draft - Audience: %s", aud),
				Subtitle:  "Created by integration test - please delete",
				DraftBody: `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"This is a test draft created by the substack-cli integration tests. Please delete this draft after reviewing."}]}]}`,
				Audience:  aud,
			}

			resp, err := client.CreateDraft(draft)
			if err != nil {
				t.Fatalf("CreateDraft(%s): %v", aud, err)
			}

			t.Logf("Created draft: id=%d title=%q audience=%s", resp.ID, resp.Title, aud)
		})
	}
}

func TestIntegration_CreateDraftWithSection(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v (are you logged in?)", err)
	}

	// Note: Section ID must be a valid section ID from your Substack.
	// If you don't have sections configured, this test may fail or be ignored.
	draft := model.DraftRequest{
		Title:         "Test Draft - With Section Field",
		Subtitle:      "Created by integration test - please delete",
		DraftBody:     `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"This draft tests the section field. Check if section_chosen is set correctly."}]}]}`,
		Audience:      "everyone",
		SectionChosen: false, // No section specified
	}

	resp, err := client.CreateDraft(draft)
	if err != nil {
		t.Fatalf("CreateDraft: %v", err)
	}

	t.Logf("Created draft without section: id=%d title=%q", resp.ID, resp.Title)
}
