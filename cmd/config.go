package cmd

var cfgPath string

type config struct {
	APIKey   string
	Owner    string
	Template string
	Editor   string
}

var cfg config
