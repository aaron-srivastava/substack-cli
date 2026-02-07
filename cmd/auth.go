package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aaronsrivastava/substack-cli/internal/auth"
	"github.com/aaronsrivastava/substack-cli/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
	}

	authCmd.AddCommand(
		&cobra.Command{
			Use:   "login",
			Short: "Add or update an account",
			RunE:  authLogin,
		},
		&cobra.Command{
			Use:   "status",
			Short: "Show active account",
			RunE:  authStatus,
		},
		&cobra.Command{
			Use:   "list",
			Short: "List all accounts",
			RunE:  authList,
		},
		&cobra.Command{
			Use:   "switch <name>",
			Short: "Switch active account",
			Args:  cobra.ExactArgs(1),
			RunE:  authSwitch,
		},
		&cobra.Command{
			Use:   "remove <name>",
			Short: "Remove an account",
			Args:  cobra.ExactArgs(1),
			RunE:  authRemove,
		},
	)

	rootCmd.AddCommand(authCmd)
}

func prompt(scanner *bufio.Scanner, label string) string {
	fmt.Fprintf(os.Stdout, "%s: ", label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func authLogin(_ *cobra.Command, _ []string) error {
	scanner := bufio.NewScanner(os.Stdin)

	name := prompt(scanner, "Account name (e.g. my-blog)")
	pubURL := prompt(scanner, "Publication URL (e.g. https://you.substack.com)")
	userID := prompt(scanner, "User ID")
	sid := prompt(scanner, "SID (connect.sid cookie)")
	substackSID := prompt(scanner, "substack.sid cookie")
	substackLLI := prompt(scanner, "substack.lli cookie")

	acct := model.Account{
		Name:           name,
		PublicationURL: pubURL,
		UserID:         userID,
		SID:            sid,
		SubstackSID:    substackSID,
		SubstackLLI:    substackLLI,
	}

	store, err := auth.Load()
	if err != nil {
		return err
	}
	auth.AddAccount(store, acct)
	store.Active = name
	if saveErr := auth.Save(store); saveErr != nil {
		return saveErr
	}
	fmt.Fprintf(os.Stdout, "Logged in as %s (active)\n", name)
	return nil
}

func authStatus(_ *cobra.Command, _ []string) error {
	store, err := auth.Load()
	if err != nil {
		return err
	}
	acct, err := auth.GetActive(store)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Active: %s (%s)\n", acct.Name, acct.PublicationURL)
	return nil
}

func authList(_ *cobra.Command, _ []string) error {
	store, err := auth.Load()
	if err != nil {
		return err
	}
	if len(store.Accounts) == 0 {
		fmt.Fprintln(os.Stdout, "No accounts configured. Run 'substack auth login'.")
		return nil
	}
	for _, a := range store.Accounts {
		marker := "  "
		if a.Name == store.Active {
			marker = "* "
		}
		fmt.Fprintf(os.Stdout, "%s%s (%s)\n", marker, a.Name, a.PublicationURL)
	}
	return nil
}

func authSwitch(_ *cobra.Command, args []string) error {
	store, err := auth.Load()
	if err != nil {
		return err
	}
	if switchErr := auth.SwitchAccount(store, args[0]); switchErr != nil {
		return switchErr
	}
	if saveErr := auth.Save(store); saveErr != nil {
		return saveErr
	}
	fmt.Fprintf(os.Stdout, "Switched to %s\n", args[0])
	return nil
}

func authRemove(_ *cobra.Command, args []string) error {
	store, err := auth.Load()
	if err != nil {
		return err
	}
	if removeErr := auth.RemoveAccount(store, args[0]); removeErr != nil {
		return removeErr
	}
	if saveErr := auth.Save(store); saveErr != nil {
		return saveErr
	}
	fmt.Fprintf(os.Stdout, "Removed %s\n", args[0])
	return nil
}
