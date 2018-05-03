package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// autocompleteCmd represents the autocomplete command.
var autocompleteCmd = &cobra.Command{
	Use:   "autocomplete",
	Short: "generate shell autocompletion script for ghb",
	Long: `Generates a shell autocompletion script for ghb.

NOTE: The current version supports Bash only.
      This should work for *nix systems with Bash installed.

By default, the file is written directly to /etc/bash_completion.d
for convenience, and the command may need superuser rights, e.g.:

	$ sudo ghb autocomplete

Add ` + "`--completionfile=/path/to/file`" + ` flag to set alternative
file-path and name.

Logout and in again to reload the completion scripts,
or just source them in directly:

	$ . /etc/bash_completion`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if autocompleteType != "bash" {
			return errors.New("only Bash is supported")
		}

		err := cmd.Root().GenBashCompletionFile(autocompleteTarget)
		if err != nil {
			return err
		}

		fmt.Println("Bash completion file for ghb saved to", autocompleteTarget)

		return nil
	},
}

// Flags
var (
	autocompleteTarget string
	autocompleteType   string
)

func init() {
	RootCmd.AddCommand(autocompleteCmd)

	autocompleteCmd.Flags().StringVarP(&autocompleteTarget, "completionfile", "",
		"/etc/bash_completion.d/ghb.sh", "autocompletion file")
	autocompleteCmd.Flags().StringVarP(&autocompleteType, "type", "",
		"bash", "autocompletion type (currently only bash supported)")
	// For bash-completion
	err := autocompleteCmd.Flags().SetAnnotation("completionfile", cobra.BashCompFilenameExt, []string{})
	if err != nil {
		panic(err)
	}
}
