package cmd

import (
	"errors"
	"os"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

var addRepositoryCmd = &cobra.Command{
	Aliases: []string{"r", "repo"},
	Use:     "repository NAME",
	Short:   "add a new repository",
	Long: `add a new repository. Specify your repository name as NAME and
short description or homepage as each flags.`,
	RunE: runAddRepository,
}

// Flags
var (
	description string
	url         string
)

func init() {
	addCmd.AddCommand(addRepositoryCmd)

	addRepositoryCmd.Flags().StringVarP(
		&description, "description", "d", "", "a short description of repository",
	)
	addRepositoryCmd.Flags().StringVarP(
		&url, "url", "u", "", "a URL with more information about the repository",
	)
}

func runAddRepository(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("too few argments: specify a repository name")
	}

	newRepo := &github.Repository{
		Name:        github.String(args[0]),
		Description: github.String(description),
		URL:         github.String(url),
	}
	ghb := newGHB()
	addedRepo, _, err := ghb.createRepo("", newRepo)
	if err != nil {
		return err
	}
	ghb.fPrintGitClone(os.Stdout, addedRepo)
	return nil
}
