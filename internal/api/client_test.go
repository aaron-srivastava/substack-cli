package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaronsrivastava/substack-cli/internal/model"
)

func testClient(handler http.HandlerFunc) (*Client, *httptest.Server) {
	srv := httptest.NewServer(handler)
	acct := &model.Account{
		Name:          "test",
		PublicationURL: srv.URL,
		UserID:        "123",
		SID:           "sid-val",
		SubstackSID:   "ssid-val",
		SubstackLLI:   "lli-val",
	}
	return NewClientWith(acct), srv
}

func TestCookiesSet(t *testing.T) {
	var gotCookies []*http.Cookie
	client, srv := testClient(func(w http.ResponseWriter, r *http.Request) {
		gotCookies = r.Cookies()
		_ = json.NewEncoder(w).Encode([]model.DraftResponse{})
	})
	defer srv.Close()

	_, _ = client.ListDrafts()

	names := map[string]string{}
	for _, c := range gotCookies {
		names[c.Name] = c.Value
	}
	if names["substack.sid"] != "ssid-val" {
		t.Errorf("substack.sid = %q", names["substack.sid"])
	}
	if names["substack.lli"] != "lli-val" {
		t.Errorf("substack.lli = %q", names["substack.lli"])
	}
	if names["connect.sid"] != "sid-val" {
		t.Errorf("connect.sid = %q", names["connect.sid"])
	}
}

func TestCreateDraft(t *testing.T) {
	callCount := 0
	client, srv := testClient(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/publication/users":
			_ = json.NewEncoder(w).Encode([]struct {
				ID   int    `json:"id"`
				Role string `json:"role"`
			}{{ID: 99, Role: "admin"}})
		case r.URL.Path == "/api/v1/drafts/":
			callCount++
			if r.Method != "POST" {
				t.Errorf("method = %s, want POST", r.Method)
			}
			_ = json.NewEncoder(w).Encode(model.DraftResponse{ID: 42, Title: "Test"})
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	})
	defer srv.Close()

	resp, err := client.CreateDraft(model.DraftRequest{Title: "Test"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 42 {
		t.Errorf("id = %d, want 42", resp.ID)
	}
	if callCount != 1 {
		t.Errorf("draft create called %d times, want 1", callCount)
	}
}

func TestListPosts(t *testing.T) {
	client, srv := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/posts/" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]model.Post{{ID: 1, Title: "A"}, {ID: 2, Title: "B"}})
	})
	defer srv.Close()

	posts, err := client.ListPosts()
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 2 {
		t.Errorf("len = %d, want 2", len(posts))
	}
}

func TestAPIError(t *testing.T) {
	client, srv := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("forbidden"))
	})
	defer srv.Close()

	_, err := client.ListPosts()
	if err == nil {
		t.Fatal("expected error")
	}
}
