package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
)

type ghb struct {
	c gitHubAPIDoer
}

func newGHB() *ghb {
	cl := newGitHubClient()
	return &ghb{c: cl}
}

func (g *ghb) createRepo(org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
	return g.c.createRepo(org, repo)
}

func (g *ghb) createIssue(repoName string, labels []string) (*github.Issue, *github.Response, error) {
	// NG! var issue *github.IssueRequest
	// issue.Title = ... <- nil pointer dereference
	var issue github.IssueRequest

	lines, err := g.c.issueLineFromEditor(repoName)
	if err != nil {
		return nil, nil, err
	}
	if len(lines) == 1 {
		issue.Title = github.String(lines[0])
	} else {
		issue.Title, issue.Body =
			github.String(lines[0]), github.String(strings.Join(lines[1:], "\n"))
	}
	if len(labels) > 0 {
		issue.Labels = &labels
	}
	return g.c.createIssue(repoName, &issue)
}

func (g *ghb) fPrintIssues(out io.Writer, repoName string) (*github.Response, error) {
	issues, _, err := g.c.getIssues(repoName)
	if err != nil {
		return nil, err
	}

	tbl := tablewriter.NewWriter(out)
	tbl.SetBorder(false)
	tbl.SetColumnSeparator("")
	for _, issue := range issues {
		var lables []string
		for _, l := range issue.Labels {
			lables = append(lables, l.String())
		}
		tbl.Append([]string{
			fmt.Sprintf("#%d", *issue.Number),
			*issue.Title,
			strings.Join(lables, ", "),
		})
	}
	tbl.Render()
	return nil, nil
}
