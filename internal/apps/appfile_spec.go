package apps

import (
	"fmt"
	"path/filepath"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/renehernandez/appfile/internal/env"
	"github.com/renehernandez/appfile/internal/log"
	"github.com/renehernandez/appfile/internal/tmpl"
	"github.com/renehernandez/appfile/internal/yaml"
)

type AppfileSpec struct {
	AppSpecs     []string            `yaml:"specs"`
	Environments map[string][]string `yaml:"environments"`

	path string
}

func (spec *AppfileSpec) Path() string {
	return spec.path
}

func (spec *AppfileSpec) SetPath(path string) error {
	var err error
	spec.path, err = filepath.Abs(path)
	return err
}

func (spec *AppfileSpec) IsValid() bool {
	return len(spec.AppSpecs) > 0
}

func (spec *AppfileSpec) hasEnvironment(name string) bool {
	_, ok := spec.Environments[name]
	return ok
}

func (spec *AppfileSpec) ReadEnvironment(name string) (*env.Environment, error) {
	if !spec.hasEnvironment(name) {
		if name != "default" {
			return &env.Environment{}, fmt.Errorf("Environment %s not found in appfile spec at %s", name, spec.Path())
		}

		log.Debugf("Using default environment without any defined values")
		return &env.Environment{}, nil
	}

	fullEnv := &env.Environment{
		Name: name,
	}

	for _, envPath := range spec.Environments[name] {
		file := filepath.Join(filepath.Dir(spec.Path()), envPath)
		log.Debugf("Reading environment values from %s", file)
		templatedYaml, err := tmpl.RenderFromFile(file)
		if err != nil {
			return &env.Environment{}, err
		}

		currentEnvPart, err := yaml.ParseEnvironment(templatedYaml)
		if err != nil {
			return &env.Environment{}, errors.Wrapf(err, "Could not parse resulting yaml from file %s in env %s", file, name)
		}

		fullEnv, err = fullEnv.Merge(currentEnvPart)
		if err != nil {
			return &env.Environment{}, errors.Wrapf(err, "Could not merge values from file %s in env %s", file, name)
		}
	}

	return fullEnv, nil
}

func (spec *AppfileSpec) readApps(state *StateData) ([]*godo.App, error) {
	apps := []*godo.App{}

	for _, appSpecPath := range spec.AppSpecs {
		file := filepath.Join(filepath.Dir(spec.Path()), appSpecPath)
		log.Debugf("Reading app spec from %s", file)
		templatedYaml, err := tmpl.RenderFromFile(file, state)
		if err != nil {
			return []*godo.App{}, err
		}

		var appSpec AppSpec
		err = yaml.ParseAppSpec(templatedYaml, &appSpec)
		if err != nil {
			return []*godo.App{}, errors.Wrapf(err, "Could not parse resulting yaml for app spec from file %s", file)
		}

		appSpec.SetDefaultValues()

		apps = append(apps, &godo.App{Spec: &appSpec.AppSpec})
	}

	return apps, nil
}
