package cmd

var cfgPath string

type config struct {
	APIKey   string
	Template string
	Editor   string
}

var cfg config
