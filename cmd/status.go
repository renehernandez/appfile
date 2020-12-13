package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/renehernandez/appfile/internal/errors"
	"github.com/spf13/cobra"
)

type statusCmd struct {
	*rootCmd
}

var (
	statusLong = `Show status for apps defined in the appfile.
`
	statusExample = `  # Status using defaults: appfile.yaml in current location, default environment and DIGITALOCEAN_ACCESS_TOKEN env var
appfile status

  # Status using appfile.yaml in custom location, review environment and access token option
  appfile status --file /path/to/appfile.yaml --environment review --access-token $TOKEN

  # Status with debug output
  appfile status --log-level debug`
)

func newStatusCmd(rootCmd *rootCmd) *cobra.Command {
	status := statusCmd{
		rootCmd: rootCmd,
	}

	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Show status for apps defined in the appfile",
		Long:    statusLong,
		Example: statusExample,
		Run: func(cmd *cobra.Command, args []string) {
			status.run()
		},
	}
	return cmd
}

func (status *statusCmd) run() {
	appfile := status.appfileFromSpec()

	appsStatus, err := appfile.Status(status.accessToken)
	errors.CheckAndFail(err)

	table := uitable.New()
	table.Wrap = true
	table.MaxColWidth = 80

	for _, status := range appsStatus {
		table.AddRow("Name:", status.Name)
		table.AddRow("Status:", status.Status)
		table.AddRow("Deployment ID:", status.DeploymentID)
		table.AddRow("Updated:", status.UpdatedAt)
		table.AddRow("URL:", status.URL)
		table.AddRow("")
	}
	fmt.Println(table)
}
