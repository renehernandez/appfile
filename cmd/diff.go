package cmd

import (
	"fmt"

	"github.com/renehernandez/appfile/internal/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

type diffCmd struct {
	*rootCmd
}

var (
	diffLong = `Diff local app spec against app spec running in DigitalOcean

`
	diffExample = `  # Diff using defaults: appfile.yaml in current location, default environment and DIGITALOCEAN_ACCESS_TOKEN env var
appfile diff

  # Diff using appfile.yaml in custom location, review environment and access token option
  appfile diff --file /path/to/appfile.yaml --environment review --access-token $TOKEN

  # Diff with debug output
  appfile sync --log-level debug`
)

func newDiffCmd(rootCmd *rootCmd) *cobra.Command {
	diff := diffCmd{
		rootCmd: rootCmd,
	}

	cmd := &cobra.Command{
		Use:     "diff",
		Short:   "Diff local app spec against app spec running in DigitalOcean",
		Long:    diffLong,
		Example: diffExample,
		Run: func(cmd *cobra.Command, args []string) {
			diff.run()
		},
	}
	return cmd
}

func (diff *diffCmd) run() {
	appfile := diff.appfileFromSpec()

	diffs, err := appfile.Diff()
	errors.CheckAndFail(err)

	dmp := diffmatchpatch.New()

	for _, appDiff := range diffs {
		appSpecDiffs, err := appDiff.CalculateDiff()
		errors.CheckAndFailf(err, "Failed to calculate diff for app %s", appDiff.Name)

		fmt.Printf("Diff for app %s\n", appDiff.Name)
		fmt.Println(dmp.DiffPrettyText(appSpecDiffs))
	}
}
