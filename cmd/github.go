package cmd

import (
	"net/http"

	"github.com/google/go-github/github"
	"github.com/yuta-masano/go-tempedit"
	"golang.org/x/oauth2"
)

// interface のメソッドシグニチャって引数名書いても OK なのか。
type gitHubAPIDoer interface {
	createRepo(org string, repo *github.Repository,
	) (*github.Repository, *github.Response, error)
	createIssue(repoName string, issue *github.IssueRequest,
	) (*github.Issue, *github.Response, error)
	issueLineFromEditor(repoName string) ([]string, error)
	getIssueTemplate(repoName, path string) (string, error)
	getIssues(repoName string) ([]*github.Issue, *github.Response, error)
}

type gitHubClient struct {
	*github.Client
}

func newGitHubClient() gitHubAPIDoer {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.gitHubToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	cl := github.NewClient(tc)
	// 構造体を無名で埋め込んだら *github.<Client> をフィールドにして代入するようだ。
	return &gitHubClient{Client: cl}
}

func (g *gitHubClient) createRepo(org string, repo *github.Repository,
) (*github.Repository, *github.Response, error) {
	return g.Repositories.Create(org, repo)
}

func (g *gitHubClient) createIssue(repoName string, issue *github.IssueRequest,
) (*github.Issue, *github.Response, error) {
	return g.Issues.Create(cfg.gitHubOwner, repoName, issue)
}

func (g *gitHubClient) issueLineFromEditor(repoName string) ([]string, error) {
	edit := tempedit.New(cfg.Editor)
	if err := edit.MakeTemp("", ".ghb."); err != nil {
		return nil, err
	}
	defer edit.CleanTempFile()

	content, err := g.getIssueTemplate(repoName, cfg.IssueTemplateRelPath)
	if err != nil {
		return nil, err
	}
	if err = edit.Write(content); err != nil {
		return nil, err
	}
	if err = edit.Run(); err != nil {
		return nil, err
	}
	if changed, err := edit.FileChanged(); !changed {
		return nil, err
	}
	return edit.Line(), nil
}

func (g *gitHubClient) getIssueTemplate(repoName, path string) (string, error) {
	file, _, res, err :=
		g.Repositories.GetContents(cfg.gitHubOwner, repoName, path, nil)
	switch res.StatusCode {
	case http.StatusOK:
		content, err := file.GetContent()
		if err != nil {
			return "", err
		}
		return content, nil
	case http.StatusNotFound: // 404 = イシューテンプレートがない = 空文字を返せばよい。
		return "", nil
	default:
		return "", err
	}
}

func (g *gitHubClient) getIssues(repoName string) ([]*github.Issue, *github.Response, error) {
	// デフォルトでステータスが Open なものを拾ってくるので nil でよかろう。
	return g.Issues.ListByRepo(cfg.gitHubOwner, repoName, nil)
}
