package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

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
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return &model.Config{Audience: "everyone", OutputFormat: "text"}, nil
		}
		return nil, readErr
	}
	var cfg model.Config
	if unmarshalErr := json.Unmarshal(data, &cfg); unmarshalErr != nil {
		return nil, unmarshalErr
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
	return slices.Contains(validAudiences, s)
}

var validOutputFormats = []string{"text", "json"}

func validOutputFormat(s string) bool {
	return slices.Contains(validOutputFormats, s)
}

func saveConfig(cfg *model.Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if mkdirErr := os.MkdirAll(filepath.Dir(path), 0700); mkdirErr != nil {
		return mkdirErr
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
	fmt.Fprintf(os.Stdout, "send_email:    %v\naudience:      %s\nsection:       %s\noutput_format: %s\n",
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
	if saveErr := saveConfig(cfg); saveErr != nil {
		return saveErr
	}
	fmt.Fprintf(os.Stdout, "Set %s = %s\n", args[0], args[1])
	return nil
}
