package yaml

import (
	"bytes"
	"encoding/json"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	"github.com/renehernandez/appfile/internal/env"
	"github.com/renehernandez/appfile/internal/log"
)

// ParseAppfileSpec parses an AppfileSpec object
func ParseAppfileSpec(yamlData *bytes.Buffer, out interface{}) error {
	log.Debugf("Unmarshaling appfile spec from yaml")

	err := yaml.Unmarshal(yamlData.Bytes(), out)

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

func ParseAppSpec(yamlData *bytes.Buffer, out interface{}) error {
	log.Debugf("Unmarshaling app spec from yaml")

	bytes := yamlData.Bytes()
	jsonData, err := yaml.YAMLToJSON(bytes)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal app spec from yaml")
	}

	err = json.Unmarshal(jsonData, out)

	if err != nil {
		return errors.Wrapf(err, "Failed to unmarshal app spec from yaml. Yaml content:\n%s", string(bytes))
	}

	return nil
}
