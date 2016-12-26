package cmd

var cfgPath string

type config struct {
	Template string
	Editor   string
}

var cfg config
