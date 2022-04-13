package cmd

import (
	"github.com/riotkit-org/git-clone-operator/cmd/checkout"
	"github.com/riotkit-org/git-clone-operator/cmd/serve"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCheckCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "check-binary",
		Short: "Does nothing. Allows to check if this binary requires libc at all",
		Run: func(command *cobra.Command, args []string) {
			println("Looks OK")
			return
		},
	}
	return command
}

func Main() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git-clone-operator",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				logrus.Errorf(err.Error())
			}
		},
	}
	cmd.AddCommand(serve.NewServeCommand())
	cmd.AddCommand(checkout.NewCheckoutCommand())
	cmd.AddCommand(NewCheckCommand())

	return cmd
}
