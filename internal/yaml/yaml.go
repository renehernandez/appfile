package yaml

import (
	"bytes"
	"encoding/json"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	"github.com/renehernandez/appfile/internal/env"
	"github.com/renehernandez/appfile/internal/log"
)

type Spec interface {
}

// ParseAppfileSpec parses an AppfileSpec object
func ParseAppfileSpec(yamlData *bytes.Buffer, spec Spec) error {
	log.Debugf("Unmarshaling appfile spec from yaml")

	err := yaml.Unmarshal(yamlData.Bytes(), spec)

	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal appfile spec from yaml")
	}

	return nil
}

// ParseEnvironment parses state values from an environment file
func ParseEnvironment(yamlData *bytes.Buffer) (*env.Environment, error) {
	log.Debugf("Unmarshaling environment from yaml")

	var newEnv env.Environment
	err := yaml.Unmarshal(yamlData.Bytes(), &newEnv.Values)

	if err != nil {
		return &env.Environment{}, errors.Wrap(err, "Failed to unmarshal environment from yaml")
	}

	return &newEnv, nil
}

func ParseAppSpec(yamlData *bytes.Buffer, spec Spec) error {
	log.Debugf("Unmarshaling app spec from yaml")

	jsonData, err := yaml.YAMLToJSON(yamlData.Bytes())
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal app spec from yaml")
	}

	err = json.Unmarshal(jsonData, spec)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal app spec from yaml")
	}

	return nil
}
