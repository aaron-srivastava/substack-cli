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
			Short: "Set a config value (send_email, audience, section, output_format)",
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
			return &model.Config{Audience: "everyone", OutputFormat: "text"}, nil
		}
		return nil, err
	}
	var cfg model.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// Apply defaults for missing values
	if cfg.Audience == "" {
		cfg.Audience = "everyone"
	}
	if cfg.OutputFormat == "" {
		cfg.OutputFormat = "text"
	}
	return &cfg, nil
}

var validAudiences = []string{"everyone", "only_paid", "only_free"}

func validAudience(s string) bool {
	for _, v := range validAudiences {
		if s == v {
			return true
		}
	}
	return false
}

var validOutputFormats = []string{"text", "json"}

func validOutputFormat(s string) bool {
	for _, v := range validOutputFormats {
		if s == v {
			return true
		}
	}
	return false
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
	fmt.Printf("send_email:    %v\naudience:      %s\nsection:       %s\noutput_format: %s\n",
		cfg.SendEmail, cfg.Audience, cfg.Section, cfg.OutputFormat)
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
		if !validAudience(args[1]) {
			return fmt.Errorf("invalid audience: %s (valid: %v)", args[1], validAudiences)
		}
		cfg.Audience = args[1]
	case "section":
		cfg.Section = args[1]
	case "output_format":
		if !validOutputFormat(args[1]) {
			return fmt.Errorf("invalid output_format: %s (valid: %v)", args[1], validOutputFormats)
		}
		cfg.OutputFormat = args[1]
	default:
		return fmt.Errorf("unknown config key: %s (valid: send_email, audience, section, output_format)", args[0])
	}
	if err := saveConfig(cfg); err != nil {
		return err
	}
	fmt.Printf("Set %s = %s\n", args[0], args[1])
	return nil
}
