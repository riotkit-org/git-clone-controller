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
	command.Flags().BoolVarP(&app.TLS, "tls", "-t", false, "Use TLS (requires certificates)")

	return command
}
