package cmd

import (
	"github.com/riotkit-org/git-clone-operator/cmd/checkout"
	"github.com/riotkit-org/git-clone-operator/cmd/serve"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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

	return cmd
}
