package cmd

var cfgPath string

type config struct {
	Editor               string
	IssueTemplateRelPath string

	gitHubOwner string
	gitHubToken string
}

var cfg config
