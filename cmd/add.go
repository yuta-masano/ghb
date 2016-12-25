package cmd

import "github.com/spf13/cobra"

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "create a new issue or repository",
}

func init() {
	RootCmd.AddCommand(addCmd)
}
