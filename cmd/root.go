package cmd

import (
	"os"

	"github.com/renehernandez/appfile/internal/apps"
	"github.com/renehernandez/appfile/internal/errors"
	"github.com/renehernandez/appfile/internal/log"
	"github.com/renehernandez/appfile/internal/tmpl"
	"github.com/renehernandez/appfile/internal/version"
	"github.com/renehernandez/appfile/internal/yaml"
	"github.com/spf13/cobra"
)

type rootCmd struct {
	environment string
	file        string
	logLevel    string
	accessToken string
}

func (root *rootCmd) Environment() string {
	return root.environment
}

func (root *rootCmd) File() string {
	return root.file
}

func (root *rootCmd) LogLevel() string {
	return root.logLevel
}

func (root *rootCmd) AccessToken() string {
	return root.accessToken
}

func NewRootCmd() *cobra.Command {
	root := rootCmd{}

	cmd := &cobra.Command{
		Use:     "appfile",
		Short:   "Deploy app platform specifications to DigitalOcean",
		Version: version.Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			root.initialize()
		},
	}

	cmd.PersistentFlags().StringVarP(&root.environment, "environment", "e", "default", "root all resources from spec file")
	cmd.PersistentFlags().StringVarP(&root.file, "file", "f", "appfile.yaml", "load appfile spec from file")
	cmd.PersistentFlags().StringVar(&root.logLevel, "log-level", "info", "Set log level")
	cmd.PersistentFlags().StringVarP(&root.accessToken, "access-token", "t", "", "API V2 access token")
	cmd.AddCommand(newDiffCmd(&root))
	cmd.AddCommand(newSyncCmd(&root))
	cmd.AddCommand(newDestroyCmd(&root))

	return cmd
}

func (root *rootCmd) initialize() {
	log.Initialize(root.LogLevel())

	if root.AccessToken() == "" {
		token, ok := os.LookupEnv("DIGITALOCEAN_ACCESS_TOKEN")
		if !ok || token == "" {
			log.Fatalf("No access token option specified and DIGITALOCEAN_ACCESS_TOKEN environment variable is not defined")
		}

		root.accessToken = token
	}
}

func (root *rootCmd) logOptions(cmd *cobra.Command) {
	log.Debugf("Invoking %s command with options: environment=%s; file=%s; log-level=%s", cmd.Name(), root.Environment(), root.File(), root.LogLevel())
}

func (root *rootCmd) appfileFromSpec() *apps.Appfile {
	log.Debugln("Start reading appfile spec")
	templatedYaml, err := tmpl.RenderFromFile(root.File())
	errors.CheckAndFail(err)

	var spec apps.AppfileSpec
	err = yaml.ParseAppfileSpec(templatedYaml, &spec)
	errors.CheckAndFailf(err, "Could not parse resulting yaml from file %s", root.File())

	err = spec.SetPath(root.File())
	errors.CheckAndFailf(err, "Could not generate absolute path for file %s", root.File())
	log.Debugln("Finished reading appfile spec")

	appfile, err := apps.NewAppfileFromSpec(&spec, root.Environment())
	errors.CheckAndFailln(err, "Could not create appfile from spec")

	return appfile
}
