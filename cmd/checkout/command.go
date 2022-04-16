package checkout

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func NewCheckoutCommand() *cobra.Command {
	app := &Command{}

	command := &cobra.Command{
		Use:   "checkout",
		Short: "Starts a HTTP server for receiving webhooks from Kubernetes API and mutating pods",
		Run: func(command *cobra.Command, args []string) {
			if len(args) == 0 {
				logrus.Errorf("Please enter a GIT url as an argument")
				os.Exit(1)
			}

			app.Url = args[0]
			err := app.Run()

			if err != nil {
				logrus.Errorf(err.Error())
				os.Exit(1)
			}
		},
	}

	command.Flags().StringVarP(&app.LogLevel, "log-level", "l", "info", "Logging level: error, warn, info, debug")
	command.Flags().StringVarP(&app.Path, "path", "p", "./", "GIT repository target path")
	command.Flags().StringVarP(&app.Username, "username", "U", "__token__", "GIT basic auth username")
	command.Flags().StringVarP(&app.Token, "token", "t", "", "GIT basic auth token/password")
	command.Flags().StringVarP(&app.Revision, "rev", "r", "", "GIT revision - commit/branch/tag (defaults to: main)")
	command.Flags().BoolVarP(&app.CleanUpRemotes, "clean-remotes", "", true, "Delete `git remote` from local repository to prevent token leak")
	app.IsBare = false

	return command
}
