package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issue REPO_NAME",
	Short: "add a new issue",
	Long: `Create a new issue about specified repository allowing you to edit
the issue subject and description via your editor.`,
	RunE: runIssue,
}

// Flags
var (
	labels []string
)

func init() {
	addCmd.AddCommand(issueCmd)

	issueCmd.Flags().StringSliceVarP(
		&labels, "labels", "l", nil, "a list of comma separated label names",
	)
}

func validateContents(before, after []byte) ([]string, error) {
	// see: TestValidateContents(t *testing.T)
	if bytes.Equal(before, after) {
		return nil, errors.New("edit aborted")
	}
	if bytes.Equal(bytes.TrimRight(before, "\n"), bytes.TrimRight(after, "\n")) {
		return nil, errors.New("no changed")
	}
	// 空文字は不許可です。
	if string(after) == "" {
		return nil, errors.New("canceled")
	}
	// Split は分割対象がないとそのまま対象がスライスの 0 番目の要素になる。
	lines := strings.Split(string(after), "\n")
	if len(lines) == 0 { // 起こり得るのか？
		return nil, errors.New("canceled")
	}
	return lines, nil
}

func issueFromEditor(content string) (*github.IssueRequest, error) {
	tmpFile, err := ioutil.TempFile("", ".ghb.")
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := tmpFile.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	err = tmpWrite(tmpFile, content)
	if err != nil {
		return nil, err
	}

	before, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	editor := getEditor()
	if err = run(editor, tmpFile.Name()); err != nil {
		return nil, err
	}
	after, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	defer func() {
		if removeErr := os.Remove(tmpFile.Name()); removeErr != nil {
			panic(removeErr)
		}
	}()
	lines, err := validateContents(before, after)
	if err != nil {
		return nil, err
	}

	var issue github.IssueRequest
	if len(lines) == 1 {
		issue.Title = github.String(lines[0])
	} else {
		issue.Title, issue.Body =
			github.String(lines[0]), github.String(strings.Join(lines[1:], "\n"))
	}
	return &issue, nil
}

func runIssue(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("too few argments: specify a repository name to add a new issue")
	}
	repoName := args[0]

	newIssue, err := issueFromEditor(cfg.Template)
	if err != nil {
		return err
	}
	if len(labels) > 0 {
		newIssue.Labels = &labels
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	cl := github.NewClient(tc)
	addedIssue, _, err := cl.Issues.Create(cfg.GitHubOwner, repoName, newIssue)
	if err != nil {
		return err
	}
	fmt.Println(strings.Trim(github.Stringify(addedIssue.Number), `"`))
	return nil
}
