package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// These values are embedded when building.
var (
	buildVersion  string
	buildRevision string
	buildWith     string
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show program's version information and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\nrevision: %s\nwith: %s\n",
			buildVersion, buildRevision, buildWith)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
