package serve

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCheckoutCommand() *cobra.Command {
	app := &Command{}

	command := &cobra.Command{
		Use:   "checkout",
		Short: "Starts a HTTP server for receiving webhooks from Kubernetes API and mutating pods",
		Run: func(command *cobra.Command, args []string) {
			err := app.Run()

			if err != nil {
				logrus.Errorf(err.Error())
			}
		},
	}

	command.Flags().StringVarP(&app.LogLevel, "log-level", "l", "info", "Logging level: error, warn, info, debug")
	command.Flags().StringVarP(&app.Path, "path", "p", "./", "GIT repository target path")
	command.Flags().StringVarP(&app.Url, "url", "p", "./", "GIT repository target path")
	command.Flags().StringVarP(&app.Username, "username", "u", "__token__", "GIT basic auth username")
	command.Flags().StringVarP(&app.Token, "token", "t", "", "GIT basic auth token/password")
	command.Flags().StringVarP(&app.Revision, "rev", "r", "", "GIT revision - commit/branch/tag")

	return command
}
