package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aaronsrivastava/substack-cli/internal/model"
)

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "substack-cli"), nil
}

func ConfigDir() (string, error) {
	return configDir()
}

func storeFile() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "accounts.json"), nil
}

func Load() (*model.AccountStore, error) {
	path, err := storeFile()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &model.AccountStore{}, nil
		}
		return nil, err
	}
	var store model.AccountStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return &store, nil
}

func LoadFrom(path string) (*model.AccountStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &model.AccountStore{}, nil
		}
		return nil, err
	}
	var store model.AccountStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return &store, nil
}

func Save(store *model.AccountStore) error {
	path, err := storeFile()
	if err != nil {
		return err
	}
	return SaveTo(store, path)
}

func SaveTo(store *model.AccountStore, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func AddAccount(store *model.AccountStore, acct model.Account) {
	for i, a := range store.Accounts {
		if a.Name == acct.Name {
			store.Accounts[i] = acct
			return
		}
	}
	store.Accounts = append(store.Accounts, acct)
	if store.Active == "" {
		store.Active = acct.Name
	}
}

func RemoveAccount(store *model.AccountStore, name string) error {
	for i, a := range store.Accounts {
		if a.Name == name {
			store.Accounts = append(store.Accounts[:i], store.Accounts[i+1:]...)
			if store.Active == name {
				store.Active = ""
				if len(store.Accounts) > 0 {
					store.Active = store.Accounts[0].Name
				}
			}
			return nil
		}
	}
	return fmt.Errorf("account %q not found", name)
}

func GetActive(store *model.AccountStore) (*model.Account, error) {
	for _, a := range store.Accounts {
		if a.Name == store.Active {
			return &a, nil
		}
	}
	return nil, fmt.Errorf("no active account; run 'substack auth login'")
}

func SwitchAccount(store *model.AccountStore, name string) error {
	for _, a := range store.Accounts {
		if a.Name == name {
			store.Active = name
			return nil
		}
	}
	return fmt.Errorf("account %q not found", name)
}
