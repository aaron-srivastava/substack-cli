package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aaronsrivastava/substack-cli/internal/api"
	"github.com/aaronsrivastava/substack-cli/internal/markdown"
	"github.com/aaronsrivastava/substack-cli/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	postCmd := &cobra.Command{
		Use:   "post",
		Short: "Manage posts",
	}

	createCmd := &cobra.Command{
		Use:   "create <file.md>",
		Short: "Create a post from markdown",
		Args:  cobra.ExactArgs(1),
		RunE:  postCreate,
	}
	createCmd.Flags().String("title", "", "Post title (overrides H1 in file)")
	createCmd.Flags().String("subtitle", "", "Post subtitle")
	createCmd.Flags().Bool("publish", false, "Publish immediately")
	createCmd.Flags().Bool("send-email", false, "Send email to subscribers")
	createCmd.Flags().String("audience", "", "Audience: everyone, only_paid, only_free")
	createCmd.Flags().String("section", "", "Section/category for the post")

	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update post metadata",
		Args:  cobra.ExactArgs(1),
		RunE:  postUpdate,
	}
	updateCmd.Flags().String("title", "", "New title")
	updateCmd.Flags().String("subtitle", "", "New subtitle")
	updateCmd.Flags().String("audience", "", "New audience")
	updateCmd.Flags().String("send-email", "", "true/false")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List published posts",
		RunE:  postList,
	}
	listCmd.Flags().String("format", "", "Output format: text or json")

	postCmd.AddCommand(
		createCmd,
		listCmd,
		&cobra.Command{
			Use:   "get <id>",
			Short: "Get post details",
			Args:  cobra.ExactArgs(1),
			RunE:  postGet,
		},
		&cobra.Command{
			Use:   "unpublish <id>",
			Short: "Unpublish a post",
			Args:  cobra.ExactArgs(1),
			RunE:  postUnpublish,
		},
		updateCmd,
	)

	rootCmd.AddCommand(postCmd)
}

func postCreate(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	source, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	fm, title, body := markdown.ConvertWithFrontmatter(source)

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshaling body: %w", err)
	}

	// Start from config defaults, then let frontmatter override, then CLI args override.
	subtitle := ""
	audience := cfg.Audience
	section := cfg.Section

	if fm != nil {
		if fm.Subtitle != "" {
			subtitle = fm.Subtitle
		}
		if fm.Audience != "" {
			audience = fm.Audience
		}
		if fm.Section != "" {
			section = fm.Section
		}
	}

	// CLI args override frontmatter
	if cmd.Flags().Changed("title") {
		title, _ = cmd.Flags().GetString("title")
	}
	if cmd.Flags().Changed("subtitle") {
		subtitle, _ = cmd.Flags().GetString("subtitle")
	}
	if cmd.Flags().Changed("audience") {
		audience, _ = cmd.Flags().GetString("audience")
	}
	if cmd.Flags().Changed("section") {
		section, _ = cmd.Flags().GetString("section")
	}

	draft := model.DraftRequest{
		Title:         title,
		Subtitle:      subtitle,
		DraftBody:     string(bodyJSON),
		Audience:      audience,
		Section:       section,
		SectionChosen: section != "",
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	resp, err := client.CreateDraft(draft)
	if err != nil {
		return fmt.Errorf("creating draft: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Draft created: id=%d title=%q\n", resp.ID, resp.Title)

	publish, _ := cmd.Flags().GetBool("publish")
	if cmd.Flags().Changed("publish") && publish {
		sendEmail, _ := cmd.Flags().GetBool("send-email")
		opts := model.PublishOptions{
			SendEmail: sendEmail,
			Audience:  audience,
		}
		post, publishErr := client.PublishDraft(resp.ID, opts)
		if publishErr != nil {
			return fmt.Errorf("publishing: %w", publishErr)
		}
		fmt.Fprintf(os.Stdout, "Published: id=%d slug=%q\n", post.ID, post.Slug)
	}

	return nil
}

func postList(cmd *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	format := cfg.OutputFormat
	if cmd.Flags().Changed("format") {
		format, _ = cmd.Flags().GetString("format")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}
	posts, err := client.ListPosts()
	if err != nil {
		return err
	}

	if format == "json" {
		data, marshalErr := json.MarshalIndent(posts, "", "  ")
		if marshalErr != nil {
			return marshalErr
		}
		fmt.Fprintln(os.Stdout, string(data))
		return nil
	}

	if len(posts) == 0 {
		fmt.Fprintln(os.Stdout, "No published posts.")
		return nil
	}
	for _, p := range posts {
		fmt.Fprintf(os.Stdout, "%-8d %s  %s\n", p.ID, p.PostDate, p.Title)
	}
	return nil
}

func postGet(_ *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid post id: %s", args[0])
	}
	client, err := api.NewClient()
	if err != nil {
		return err
	}
	post, err := client.GetPost(id)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "ID:       %d\nTitle:    %s\nSubtitle: %s\nSlug:     %s\nAudience: %s\nDate:     %s\n",
		post.ID, post.Title, post.Subtitle, post.Slug, post.Audience, post.PostDate)
	return nil
}

func postUnpublish(_ *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid post id: %s", args[0])
	}
	client, err := api.NewClient()
	if err != nil {
		return err
	}
	if unpublishErr := client.UnpublishPost(id); unpublishErr != nil {
		return unpublishErr
	}
	fmt.Fprintf(os.Stdout, "Post %d unpublished.\n", id)
	return nil
}

func postUpdate(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid post id: %s", args[0])
	}

	updates := map[string]any{}
	if t, _ := cmd.Flags().GetString("title"); t != "" {
		updates["title"] = t
	}
	if s, _ := cmd.Flags().GetString("subtitle"); s != "" {
		updates["subtitle"] = s
	}
	if a, _ := cmd.Flags().GetString("audience"); a != "" {
		updates["audience"] = a
	}

	if len(updates) == 0 {
		return errors.New("no updates specified")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}
	post, err := client.UpdatePost(id, updates)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Updated: id=%d title=%q\n", post.ID, post.Title)
	return nil
}
