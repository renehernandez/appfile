package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/renehernandez/appfile/internal/errors"
	"github.com/spf13/cobra"
)

type listCmd struct {
	*rootCmd
}

var (
	listLong = `List all apps defined in the appfile.

Optionally list components defined in a particular app.
`
	listExample = `  # List using defaults: appfile.yaml in current location, default environment and DIGITALOCEAN_ACCESS_TOKEN env var
appfile list

  # Diff using appfile.yaml in custom location, review environment and access token option
  appfile list --file /path/to/appfile.yaml --environment review --access-token $TOKEN

  # List components in app with debug output
  appfile list --log-level debug`
)

func newListCmd(rootCmd *rootCmd) *cobra.Command {
	list := listCmd{
		rootCmd: rootCmd,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all apps defined in the appfile",
		Long:    listLong,
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			list.run()
		},
	}
	return cmd
}

func (list *listCmd) run() {
	appfile := list.appfileFromSpec()

	appsStatus, err := appfile.List(list.accessToken)
	errors.CheckAndFail(err)

	table := uitable.New()

	table.AddRow("NAME", "STATUS", "DEPLOYMENT ID", "UPDATED", "URL")
	for _, status := range appsStatus {
		table.AddRow(status.Name, status.Status, status.DeploymentID, status.UpdatedAt, status.URL)
	}
	fmt.Println(table)
}
