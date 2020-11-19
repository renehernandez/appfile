package cmd

import (
	"github.com/renehernandez/appfile/internal/errors"
	"github.com/spf13/cobra"
)

type destroyCmd struct {
	*rootCmd
}

func newDestroyCmd(rootCmd *rootCmd) *cobra.Command {
	destroy := destroyCmd{
		rootCmd: rootCmd,
	}

	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Deploy app platform specifications to DigitalOcean",
		Run: func(cmd *cobra.Command, args []string) {
			destroy.run()
		},
	}
	return cmd
}

func (destroy *destroyCmd) run() {
	appfile := destroy.appfileFromSpec()

	err := appfile.Destroy(destroy.accessToken)
	errors.CheckAndFail(err)
}
