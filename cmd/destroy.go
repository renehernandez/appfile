package cmd

import (
	"github.com/renehernandez/appfile/internal/errors"
	"github.com/spf13/cobra"
)

type destroyCmd struct {
	*rootCmd
}

var (
	destroyLong = `Destroy apps running in DigitalOcean

It fails without deleting any app if any of the apps declared in the appfile spec is not found in DigitalOcean
`

	destroyExample = `  # Destroy using defaults: appfile.yaml in current location, default environment and DIGITALOCEAN_ACCESS_TOKEN env var
appfile destroy

  # Destroy using appfile.yaml in custom location, review environment and access token option
  appfile destroy --file /path/to/appfile.yaml --environment review --access-token $TOKEN

  # Destroy with debug output
  appfile destroy --log-level debug`
)

func newDestroyCmd(rootCmd *rootCmd) *cobra.Command {
	destroy := destroyCmd{
		rootCmd: rootCmd,
	}

	cmd := &cobra.Command{
		Use:     "destroy",
		Short:   "Destroy apps running in DigitalOcean",
		Long:    destroyLong,
		Example: destroyExample,
		Run: func(cmd *cobra.Command, args []string) {
			destroy.run()
		},
	}
	return cmd
}

func (destroy *destroyCmd) run() {
	appfile := destroy.appfileFromSpec()

	err := appfile.Destroy()
	errors.CheckAndFail(err)
}
