package cmd

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// repositoryCmd represents the repository command
var repositoryCmd = &cobra.Command{
	Use:   "repository NAME",
	Short: "create a new repository",
	Long: `Create a new repository. Specify your repository name as FILE and
short description or homepage as each flags.`,
	RunE: runRepository,
}

// Flags
var (
	description string
	url         string
)

func init() {
	addCmd.AddCommand(repositoryCmd)

	repositoryCmd.Flags().StringVarP(
		&description, "description", "d", "", "a short description of repository",
	)
	repositoryCmd.Flags().StringVarP(
		&url, "url", "", "", "a URL with more information about the repository",
	)
}

func runRepository(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("too few argments: specify a repository name")
	}

	newRepo := &github.Repository{
		Name:        github.String(args[0]),
		Description: github.String(description),
		URL:         github.String(url),
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("gitHubToken")})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	cl := github.NewClient(tc)
	addedRepo, _, err := cl.Repositories.Create("", newRepo)
	if err != nil {
		return err
	}
	fmt.Println(strings.Trim(github.Stringify(addedRepo.GitURL), `"`))
	return nil
}
