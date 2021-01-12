package cmd

import (
	"github.com/renehernandez/appfile/internal/errors"
	"github.com/renehernandez/appfile/internal/log"
	"github.com/spf13/cobra"
)

type syncCmd struct {
	*rootCmd
}

var (
	syncLong = `Sync all resources from app platform specs to DigitalOcean

If there is no app with the existing name, a new app will be create.
Otherwise the existing app will be updated with the changes in the spec.
`
	syncExample = `  # Sync using defaults: appfile.yaml in current location, default environment and DIGITALOCEAN_ACCESS_TOKEN env var
appfile sync

  # Sync using appfile.yaml in custom location, review environment and access token option
  appfile sync --file /path/to/appfile.yaml --environment review --access-token $TOKEN

  # Sync with debug output
  appfile sync --log-level debug`
)

func newSyncCmd(rootCmd *rootCmd) *cobra.Command {
	sync := syncCmd{
		rootCmd: rootCmd,
	}

	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "Sync all resources from app platform specs to DigitalOcean",
		Long:    syncLong,
		Example: syncExample,
		Run: func(cmd *cobra.Command, args []string) {
			sync.run()
		},
	}
	return cmd
}

func (sync *syncCmd) run() {
	appfile := sync.appfileFromSpec()

	err := appfile.Sync()
	errors.CheckAndFail(err)

	for _, app := range appfile.LocalApps {
		for _, domain := range app.Spec.Domains {
			log.Infof("%s app will be accessible at %s", app.Spec.Name, domain.Domain)
		}
	}
}
