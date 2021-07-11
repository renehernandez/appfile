package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/renehernandez/appfile/internal/apps"
	"github.com/renehernandez/appfile/internal/errors"
	"github.com/renehernandez/appfile/internal/log"
	"github.com/renehernandez/appfile/internal/tmpl"
	"github.com/renehernandez/appfile/internal/version"
	"github.com/renehernandez/appfile/internal/yaml"
	"github.com/spf13/cobra"
)

type loadEnvVarsError struct {
	message string
}

func (error loadEnvVarsError) Error() string {
	return error.message
}

type configureAccessTokenError struct {
	message string
}

func (error configureAccessTokenError) Error() string {
	return error.message
}

type rootCmd struct {
	environment string
	file        string
	logLevel    string
	accessToken string
	envFile     string
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
			if cmd.Name() == "help" {
				return
			}
			if err := root.initialize(); err != nil {
				log.Fatalln(err.Error())
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&root.environment, "environment", "e", "default", "specify the environment name")
	cmd.PersistentFlags().StringVarP(&root.file, "file", "f", "appfile.yaml", "load appfile spec from file")
	cmd.PersistentFlags().StringVar(&root.logLevel, "log-level", "info", "set log level")
	cmd.PersistentFlags().StringVarP(&root.accessToken, "access-token", "t", "", "API V2 access token")
	cmd.PersistentFlags().StringVar(&root.envFile, "env-file", ".env", "path to env file")
	cmd.AddCommand(newDiffCmd(&root))
	cmd.AddCommand(newSyncCmd(&root))
	cmd.AddCommand(newDestroyCmd(&root))
	cmd.AddCommand(newStatusCmd(&root))
	cmd.AddCommand(newLintCmd(&root))

	return cmd
}

func (root *rootCmd) initialize() error {
	log.Initialize(root.LogLevel())

	if err := root.loadEnvVars(); err != nil {
		log.Debugln(err.Error())
	}

	return root.verifyAccessToken()
}

func (root *rootCmd) loadEnvVars() error {
	if err := godotenv.Load(root.envFile); err != nil {
		return loadEnvVarsError{
			message: fmt.Sprintf("Unable to load env file at %s. Error: %s", root.envFile, err),
		}
	}

	return nil
}

func (root *rootCmd) verifyAccessToken() error {
	if root.AccessToken() == "" {
		token, ok := os.LookupEnv("DIGITALOCEAN_ACCESS_TOKEN")
		if !ok || token == "" {
			return &configureAccessTokenError{
				message: "No access token option specified and DIGITALOCEAN_ACCESS_TOKEN environment variable is not defined",
			}
		}

		root.accessToken = token
	}

	return nil
}

func (root *rootCmd) logOptions(cmd *cobra.Command) {
	log.Debugf("Invoking %s command with options: environment=%s; file=%s; log-level=%s", cmd.Name(), root.Environment(), root.File(), root.LogLevel())
}

func (root *rootCmd) appfileFromSpec() *apps.Appfile {
	log.Debugln("Start parsing appfile spec")
	templatedYaml, err := tmpl.RenderFromFile(root.File())
	errors.CheckAndFail(err)

	var spec apps.AppfileSpec
	err = yaml.ParseAppfileSpec(templatedYaml, &spec)
	var appfile *apps.Appfile

	if err != nil || !spec.IsValid() {
		log.Debugf("Could not parse appfile specification from file %s", root.File())
		log.Debugf("Try parsing app specification instead")

		appSpec := apps.NewAppSpec()
		err = yaml.ParseAppSpec(templatedYaml, &appSpec)
		errors.CheckAndFailf(err, "Could not parse app specification from file %s", root.File())

		appfile, err = apps.NewAppfileFromAppSpec(appSpec, root.AccessToken())
	} else {
		err = spec.SetPath(root.File())
		errors.CheckAndFailf(err, "Could not generate absolute path for file %s", root.File())
		log.Debugln("Finished reading appfile spec")

		appfile, err = apps.NewAppfileFromSpec(&spec, root.Environment(), root.AccessToken())
	}

	errors.CheckAndFail(err)

	return appfile
}
