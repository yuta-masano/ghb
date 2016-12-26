package cmd

var cfgPath string

type config struct {
	Template string
	Editor   string

	GitHubOwner string
	GitHubToken string
}

var cfg config
