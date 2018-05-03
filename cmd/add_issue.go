package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
)

var addIssueCmd = &cobra.Command{
	Aliases: []string{"i"},
	Use:     "issue [REPO_NAME]",
	Short:   "add a new issue",
	Long: `Add a new issue about specified repository allowing you to edit
the issue subject and description via your editor.
If REPO_NAME is omitted, ghb tries to read it from .git/config file.`,
	RunE: runAddIssue,
}

// Flags
var (
	labels []string
)

func init() {
	addCmd.AddCommand(addIssueCmd)

	addIssueCmd.Flags().StringSliceVarP(
		&labels, "label", "l", nil, "a list of comma separated label names",
	)
}

func runAddIssue(cmd *cobra.Command, args []string) error {
	var repoName string
	var err error

	repoName, err = gitconfig.Repository()
	if err != nil {
		if len(args) < 1 {
			return errors.New("too few argments: specify a repository name to add a new issue")
		}
		repoName = args[0]
	}

	ghb := newGHB()
	addedIssue, _, err := ghb.createIssue(repoName, labels)
	if err != nil {
		return err
	}
	fmt.Printf("Issue #%d created\n", *addedIssue.Number)
	return nil
}
