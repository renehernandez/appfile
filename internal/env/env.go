// Heavily based on equivalent package in https://github.com/roboll/helmfile/

package env

import (
	"github.com/goccy/go-yaml"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/renehernandez/appfile/internal/maputil"
)

type Environment struct {
	Name   string
	Values map[string]interface{}
}

func (e Environment) deepCopy() (Environment, error) {
	valuesBytes, err := yaml.Marshal(e.Values)
	if err != nil {
		return Environment{}, errors.Wrapf(err, "Failed to marshal content of %s environment", e.Name)
	}

	var values map[string]interface{}
	if err := yaml.Unmarshal(valuesBytes, &values); err != nil {
		return Environment{}, errors.Wrapf(err, "Failed to unmarshal content of %s environment", e.Name)
	}
	values, err = maputil.CastKeysToStrings(values)
	if err != nil {
		return Environment{}, errors.Wrapf(err, "Failed to convert keys to string on %s environment", e.Name)
	}

	return Environment{
		Name:   e.Name,
		Values: values,
	}, nil
}

func (e *Environment) Merge(other *Environment) (*Environment, error) {
	copy, err := e.deepCopy()

	if other != nil && err == nil {
		if err = mergo.Merge(&copy.Values, other.Values, mergo.WithOverride, mergo.WithOverwriteWithEmptyValue); err != nil {
			return nil, err
		}
	}
	return &copy, err
}
