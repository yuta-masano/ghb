package cmd

import (
	"testing"

	"github.com/google/go-github/github"
)

type fakeGitHub struct {
	gitHubAPIDoer
	// 関数型のメソッドシグニチャも引数名書いて OK なのか。
	fakeCreateRepo  func(org string, repo *github.Repository) (*github.Repository, *github.Response, error)
	fakeCreateIssue func(repoName string, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
}

func (f *fakeGitHub) createRepo(org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
	return f.fakeCreateRepo(org, repo)
}

func (f *fakeGitHub) createIssue(repoName string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	return f.fakeCreateIssue(repoName, issue)
}

func TestCreateRepo(t *testing.T) {
	t.Parallel()

	fakeGHB := &fakeGitHub{
		fakeCreateRepo: func(org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
			return nil, nil, nil
		},
		// fakeCreateIssue: func(repoName string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
		// 	return nil, nil, nil
		// },
	}

	testGHB := &ghb{c: fakeGHB}

	repo, res, error := testGHB.createRepo("", nil)
	if repo != nil || res != nil || error != nil {
		t.Fail()
	}
}
