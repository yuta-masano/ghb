package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
)

// listIssue represents the issue command
var listIssue = &cobra.Command{
	Aliases: []string{"i", "issues"},
	Use:     "issue [REPO_NAME]",
	Short:   "list issues",
	Long: `List opening issues about specified repository.
If REPO_NAME is omitted, ghb tries to read it from .git/config file.`,
	RunE: runListIssue,
}

func init() {
	listCmd.AddCommand(listIssue)
}

func runListIssue(cmd *cobra.Command, args []string) error {
	var repoName string
	var err error

	repoName, err = gitconfig.Repository()
	if err != nil {
		if len(args) < 1 {
			return errors.New("too few argments: specify a repository name to list issues")
		}
		repoName = args[0]
	}

	ghb := newGHB()
	_, err = ghb.fPrintIssues(os.Stdout, repoName)
	return err
}
