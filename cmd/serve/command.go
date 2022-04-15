package serve

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	app := &Command{}

	command := &cobra.Command{
		Use:   "serve",
		Short: "Starts a HTTP server for receiving webhooks from Kubernetes API and mutating pods",
		Run: func(command *cobra.Command, args []string) {
			err := app.Run()

			if err != nil {
				logrus.Errorf(err.Error())
			}
		},
	}

	command.Flags().StringVarP(&app.LogLevel, "log-level", "l", "info", "Logging level: error, warn, info, debug")
	command.Flags().BoolVarP(&app.LogJSON, "log-json", "", false, "Log in JSON format")
	command.Flags().BoolVarP(&app.TLS, "tls", "t", false, "Use TLS (requires certificates)")
	command.Flags().StringVarP(&app.DefaultImage, "default-image", "i", "ghcr.io/riotkit-org/git-clone-operator:latest", "Default container image")
	command.Flags().StringVarP(&app.DefaultGitUsername, "default-git-username", "U", "__token__", "Default GIT username for HTTPS auth")
	command.Flags().StringVarP(&app.DefaultGitToken, "default-git-token", "T", "", "Default GIT token/password for HTTPS auth")

	return command
}
