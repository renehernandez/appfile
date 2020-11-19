package cmd

import (
	"github.com/renehernandez/appfile/internal/errors"
	"github.com/spf13/cobra"
)

type syncCmd struct {
	*rootCmd
}

func newSyncCmd(rootCmd *rootCmd) *cobra.Command {
	sync := syncCmd{
		rootCmd: rootCmd,
	}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Deploy app platform specifications to DigitalOcean",
		Run: func(cmd *cobra.Command, args []string) {
			sync.run()
		},
	}
	return cmd
}

func (sync *syncCmd) run() {
	appfile := sync.appfileFromSpec()

	err := appfile.Sync(sync.accessToken)
	errors.CheckAndFail(err)
}
