package cmd

import (
	"github.com/renehernandez/appfile/internal/errors"
	"github.com/renehernandez/appfile/internal/log"
	"github.com/spf13/cobra"
)

type lintCmd struct {
	*rootCmd
}

var (
	lintLong = `Lint the apps specifications against the App Specification Reference.

	For more details, check the Reference at https://www.digitalocean.com/docs/app-platform/references/app-specification-reference/
`
	lintExample = `  # Status using defaults: appfile.yaml in current location, default environment and DIGITALOCEAN_ACCESS_TOKEN env var
appfile status

  # Status using appfile.yaml in custom location, review environment and access token option
  appfile status --file /path/to/appfile.yaml --environment review --access-token $TOKEN

  # Status with debug output
  appfile status --log-level debug`
)

func newLintCmd(rootCmd *rootCmd) *cobra.Command {
	lint := lintCmd{
		rootCmd: rootCmd,
	}

	cmd := &cobra.Command{
		Use:     "lint",
		Short:   "Lint the apps definitions against the App Specification Reference",
		Long:    lintLong,
		Example: lintExample,
		Run: func(cmd *cobra.Command, args []string) {
			lint.run()
		},
	}
	return cmd
}

func (lint *lintCmd) run() {
	appfile := lint.appfileFromSpec()

	lints, err := appfile.Lint()
	errors.CheckAndFail(err)

	for _, lint := range lints {
		if lint.Error != nil {
			log.Errorf("[%s] %s", lint.Name, lint.Error)
		} else {
			log.Infof("[%s] lint ran successfully", lint.Name)
		}
	}
}
