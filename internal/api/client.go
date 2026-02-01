package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aaronsrivastava/substack-cli/internal/auth"
	"github.com/aaronsrivastava/substack-cli/internal/model"
)

type Client struct {
	HTTP    *http.Client
	Account *model.Account
}

func NewClient() (*Client, error) {
	store, err := auth.Load()
	if err != nil {
		return nil, err
	}
	acct, err := auth.GetActive(store)
	if err != nil {
		return nil, err
	}
	return &Client{HTTP: &http.Client{}, Account: acct}, nil
}

func NewClientWith(acct *model.Account) *Client {
	return &Client{HTTP: &http.Client{}, Account: acct}
}

func (c *Client) baseURL() string {
	url := c.Account.PublicationURL
	url = strings.TrimRight(url, "/")
	return url
}

func (c *Client) do(method, url string, body any) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.AddCookie(&http.Cookie{Name: "substack.sid", Value: c.Account.SubstackSID, Domain: ".substack.com"})
	req.AddCookie(&http.Cookie{Name: "substack.lli", Value: c.Account.SubstackLLI, Domain: ".substack.com"})
	req.AddCookie(&http.Cookie{Name: "connect.sid", Value: c.Account.SID, Domain: ".substack.com"})

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}
	return resp, nil
}

func decodeJSON[T any](resp *http.Response) (T, error) {
	var result T
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}

func (c *Client) userID() (int, error) {
	url := fmt.Sprintf("%s/api/v1/publication/users", c.baseURL())
	resp, err := c.do("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("fetching publication users: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var users []struct {
		ID   int    `json:"id"`
		Role string `json:"role"`
	}
	if err := json.Unmarshal(raw, &users); err != nil {
		return 0, fmt.Errorf("decoding publication users: %w", err)
	}
	for _, u := range users {
		if u.Role == "admin" || u.Role == "owner" {
			return u.ID, nil
		}
	}
	if len(users) > 0 {
		return users[0].ID, nil
	}
	return 0, fmt.Errorf("no users found for publication")
}

func (c *Client) CreateDraft(draft model.DraftRequest) (*model.DraftResponse, error) {
	if len(draft.DraftBylines) == 0 {
		uid, err := c.userID()
		if err != nil {
			return nil, fmt.Errorf("resolving user ID for bylines: %w", err)
		}
		draft.DraftBylines = []model.Byline{{ID: uid}}
	}
	if draft.Type == "" {
		draft.Type = "newsletter"
	}
	url := fmt.Sprintf("%s/api/v1/drafts/", c.baseURL())
	resp, err := c.do("POST", url, draft)
	if err != nil {
		return nil, err
	}
	return ptr(decodeJSON[model.DraftResponse](resp))
}

func (c *Client) GetDraft(id int) (*model.DraftResponse, error) {
	url := fmt.Sprintf("%s/api/v1/drafts/%d", c.baseURL(), id)
	resp, err := c.do("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return ptr(decodeJSON[model.DraftResponse](resp))
}

func (c *Client) ListDrafts() ([]model.DraftResponse, error) {
	url := fmt.Sprintf("%s/api/v1/drafts/", c.baseURL())
	resp, err := c.do("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]model.DraftResponse](resp)
}

func (c *Client) DeleteDraft(id int) error {
	url := fmt.Sprintf("%s/api/v1/drafts/%d", c.baseURL(), id)
	resp, err := c.do("DELETE", url, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) PublishDraft(id int, opts model.PublishOptions) (*model.Post, error) {
	url := fmt.Sprintf("%s/api/v1/drafts/%d/publish", c.baseURL(), id)
	resp, err := c.do("PUT", url, opts)
	if err != nil {
		return nil, err
	}
	return ptr(decodeJSON[model.Post](resp))
}

func (c *Client) ListPosts() ([]model.Post, error) {
	url := fmt.Sprintf("%s/api/v1/posts/", c.baseURL())
	resp, err := c.do("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]model.Post](resp)
}

func (c *Client) GetPost(id int) (*model.Post, error) {
	url := fmt.Sprintf("%s/api/v1/posts/%d", c.baseURL(), id)
	resp, err := c.do("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return ptr(decodeJSON[model.Post](resp))
}

func (c *Client) UnpublishPost(id int) error {
	url := fmt.Sprintf("%s/api/v1/posts/%d/unpublish", c.baseURL(), id)
	resp, err := c.do("PUT", url, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) UpdatePost(id int, updates map[string]any) (*model.Post, error) {
	url := fmt.Sprintf("%s/api/v1/drafts/%d", c.baseURL(), id)
	resp, err := c.do("PUT", url, updates)
	if err != nil {
		return nil, err
	}
	return ptr(decodeJSON[model.Post](resp))
}

func ptr[T any](v T, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	return &v, nil
}
