package tmpl

import (
	"fmt"
	"os"
	"text/template"

	"github.com/goccy/go-yaml"
)

func createFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"requiredEnv": requiredEnv,
		"toYaml":      toYaml,
	}

	return funcMap
}

func requiredEnv(name string) (string, error) {
	if val, exists := os.LookupEnv(name); exists && len(val) > 0 {
		return val, nil
	}

	return "", fmt.Errorf("required env var `%s` is not set", name)
}

func toYaml(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
