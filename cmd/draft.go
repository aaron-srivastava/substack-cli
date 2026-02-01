package cmd

import (
	"fmt"
	"strconv"

	"github.com/aaronsrivastava/substack-cli/internal/api"
	"github.com/aaronsrivastava/substack-cli/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	draftCmd := &cobra.Command{
		Use:   "draft",
		Short: "Manage drafts",
	}

	publishCmd := &cobra.Command{
		Use:   "publish <id>",
		Short: "Publish a draft",
		Args:  cobra.ExactArgs(1),
		RunE:  draftPublish,
	}
	publishCmd.Flags().Bool("send-email", false, "Send email to subscribers")
	publishCmd.Flags().String("audience", "everyone", "Audience: everyone or only_paid")

	draftCmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List drafts",
			RunE:  draftList,
		},
		&cobra.Command{
			Use:   "get <id>",
			Short: "Get draft details",
			Args:  cobra.ExactArgs(1),
			RunE:  draftGet,
		},
		&cobra.Command{
			Use:   "delete <id>",
			Short: "Delete a draft",
			Args:  cobra.ExactArgs(1),
			RunE:  draftDelete,
		},
		publishCmd,
	)

	rootCmd.AddCommand(draftCmd)
}

func draftList(_ *cobra.Command, _ []string) error {
	client, err := api.NewClient()
	if err != nil {
		return err
	}
	drafts, err := client.ListDrafts()
	if err != nil {
		return err
	}
	if len(drafts) == 0 {
		fmt.Println("No drafts.")
		return nil
	}
	for _, d := range drafts {
		fmt.Printf("%-8d %s\n", d.ID, d.Title)
	}
	return nil
}

func draftGet(_ *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid draft id: %s", args[0])
	}
	client, err := api.NewClient()
	if err != nil {
		return err
	}
	d, err := client.GetDraft(id)
	if err != nil {
		return err
	}
	fmt.Printf("ID:       %d\nTitle:    %s\nSubtitle: %s\nSlug:     %s\nAudience: %s\n",
		d.ID, d.Title, d.Subtitle, d.Slug, d.Audience)
	return nil
}

func draftDelete(_ *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid draft id: %s", args[0])
	}
	client, err := api.NewClient()
	if err != nil {
		return err
	}
	if err := client.DeleteDraft(id); err != nil {
		return err
	}
	fmt.Printf("Draft %d deleted.\n", id)
	return nil
}

func draftPublish(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid draft id: %s", args[0])
	}
	sendEmail, _ := cmd.Flags().GetBool("send-email")
	audience, _ := cmd.Flags().GetString("audience")

	client, err := api.NewClient()
	if err != nil {
		return err
	}
	post, err := client.PublishDraft(id, model.PublishOptions{
		SendEmail: sendEmail,
		Audience:  audience,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Published: id=%d slug=%q\n", post.ID, post.Slug)
	return nil
}
