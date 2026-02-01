package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aaronsrivastava/substack-cli/internal/auth"
	"github.com/aaronsrivastava/substack-cli/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage default settings",
	}

	configCmd.AddCommand(
		&cobra.Command{
			Use:   "show",
			Short: "Show current config",
			RunE:  configShow,
		},
		&cobra.Command{
			Use:   "set <key> <value>",
			Short: "Set a config value (send_email, audience)",
			Args:  cobra.ExactArgs(2),
			RunE:  configSet,
		},
	)

	rootCmd.AddCommand(configCmd)
}

func configPath() (string, error) {
	dir, err := auth.ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func loadConfig() (*model.Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &model.Config{Audience: "everyone"}, nil
		}
		return nil, err
	}
	var cfg model.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveConfig(cfg *model.Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func configShow(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("send_email: %v\naudience:   %s\n", cfg.SendEmail, cfg.Audience)
	return nil
}

func configSet(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	switch args[0] {
	case "send_email":
		cfg.SendEmail = args[1] == "true"
	case "audience":
		cfg.Audience = args[1]
	default:
		return fmt.Errorf("unknown config key: %s (valid: send_email, audience)", args[0])
	}
	if err := saveConfig(cfg); err != nil {
		return err
	}
	fmt.Printf("Set %s = %s\n", args[0], args[1])
	return nil
}
