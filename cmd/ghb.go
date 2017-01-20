package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/go-github/github"
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
	// var issue *github.IssueRequest にしてはいけない。
	// issue.Labels が nil pointer dereference になる。
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

func (g *ghb) printGitClone(out io.Writer, repo *github.Repository) {
	fmt.Fprintf(out, "git clone %s\n", strings.Trim(github.Stringify(repo.GitURL), `"`))
}

func (g *ghb) printIssueNum(out io.Writer, issue *github.Issue) {
	fmt.Fprintln(out, "Issue num:", strings.Trim(github.Stringify(issue.Number), `"`))
}

func (g *ghb) fPrintIssues(out io.Writer, repoName string) (*github.Response, error) {
	issues, _, err := g.c.getIssues(repoName)
	if err != nil {
		return nil, err
	}

	for _, issue := range issues {
		fmt.Fprintf(out, "#%d\t%s\n", *issue.Number, *issue.Title)
	}
	return nil, nil
}
