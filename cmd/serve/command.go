package serve

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
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

	command.Flags().StringVarP(&app.LogLevel, "log-level", "l", getEnvOrDefault("LOG_LEVEL", "info").(string), "Logging level: error, warn, info, debug")
	command.Flags().BoolVarP(&app.LogJSON, "log-json", "", getEnvOrDefault("LOG_JSON", false).(bool), "Log in JSON format")
	command.Flags().BoolVarP(&app.TLS, "tls", "t", getEnvOrDefault("USE_TLS", false).(bool), "Use TLS (requires certificates)")
	command.Flags().StringVarP(&app.DefaultImage, "default-image", "i", getEnvOrDefault("DEFAULT_IMAGE", "ghcr.io/riotkit-org/git-clone-operator:master").(string), "Default container image")
	command.Flags().StringVarP(&app.DefaultGitUsername, "default-git-username", "U", getEnvOrDefault("DEFAULT_GIT_USERNAME", "__token__").(string), "Default GIT username for HTTPS auth")
	command.Flags().StringVarP(&app.DefaultGitToken, "default-git-token", "T", getEnvOrDefault("DEFAULT_GIT_TOKEN", "").(string), "Default GIT token/password for HTTPS auth")

	return command
}

func getEnvOrDefault(name string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}
	if value == "1" || value == "true" || value == "TRUE" || value == "yes" || value == "YES" || value == "Y" || value == "y" {
		return true
	}
	if value == "0" || value == "false" || value == "FALSE" || value == "no" || value == "NO" || value == "N" || value == "n" {
		return false
	}
	return value
}
