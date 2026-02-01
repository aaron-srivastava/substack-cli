package auth

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aaronsrivastava/substack-cli/internal/model"
)

func tmpStore(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "accounts.json")
}

func TestAddAndGetActive(t *testing.T) {
	store := &model.AccountStore{}
	acct := model.Account{Name: "test", PublicationURL: "https://test.substack.com", UserID: "1"}
	AddAccount(store, acct)

	if store.Active != "test" {
		t.Errorf("active = %q, want test", store.Active)
	}
	got, err := GetActive(store)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "test" {
		t.Errorf("name = %q", got.Name)
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tmpStore(t)
	store := &model.AccountStore{}
	AddAccount(store, model.Account{Name: "a", PublicationURL: "https://a.substack.com"})

	if err := SaveTo(store, path); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Accounts) != 1 || loaded.Active != "a" {
		t.Errorf("loaded = %+v", loaded)
	}
}

func TestRemoveAccount(t *testing.T) {
	store := &model.AccountStore{}
	AddAccount(store, model.Account{Name: "a"})
	AddAccount(store, model.Account{Name: "b"})
	store.Active = "a"

	if err := RemoveAccount(store, "a"); err != nil {
		t.Fatal(err)
	}
	if len(store.Accounts) != 1 {
		t.Errorf("accounts = %d, want 1", len(store.Accounts))
	}
	if store.Active != "b" {
		t.Errorf("active = %q, want b", store.Active)
	}
}

func TestRemoveNotFound(t *testing.T) {
	store := &model.AccountStore{}
	if err := RemoveAccount(store, "nope"); err == nil {
		t.Error("expected error")
	}
}

func TestSwitchAccount(t *testing.T) {
	store := &model.AccountStore{}
	AddAccount(store, model.Account{Name: "a"})
	AddAccount(store, model.Account{Name: "b"})

	if err := SwitchAccount(store, "b"); err != nil {
		t.Fatal(err)
	}
	if store.Active != "b" {
		t.Errorf("active = %q", store.Active)
	}
}

func TestLoadNonexistent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nope.json")
	store, err := LoadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(store.Accounts) != 0 {
		t.Errorf("expected empty store")
	}
}

func TestFilePermissions(t *testing.T) {
	path := tmpStore(t)
	store := &model.AccountStore{}
	AddAccount(store, model.Account{Name: "sec"})
	if err := SaveTo(store, path); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("permissions = %o, want 0600", perm)
	}
}
