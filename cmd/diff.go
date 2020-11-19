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

func newDiffCmd(rootCmd *rootCmd) *cobra.Command {
	diff := diffCmd{
		rootCmd: rootCmd,
	}

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Deploy app platform specifications to DigitalOcean",
		Run: func(cmd *cobra.Command, args []string) {
			diff.run()
		},
	}
	return cmd
}

func (diff *diffCmd) run() {
	appfile := diff.appfileFromSpec()

	diffs, err := appfile.Diff(diff.accessToken)
	errors.CheckAndFail(err)

	dmp := diffmatchpatch.New()

	for _, appDiff := range diffs {
		appSpecDiffs, err := appDiff.CalculateDiff()
		errors.CheckAndFailf(err, "Failed to calculate diff for app %s", appDiff.Name)

		fmt.Printf("Diff for app %s\n", appDiff.Name)
		fmt.Println(dmp.DiffPrettyText(appSpecDiffs))
	}
}
