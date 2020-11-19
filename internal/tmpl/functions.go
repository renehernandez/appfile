package tmpl

import (
	"fmt"
	"os"
	"text/template"
)

func createFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"requiredEnv": requiredEnv,
	}

	return funcMap
}

func requiredEnv(name string) (string, error) {
	if val, exists := os.LookupEnv(name); exists && len(val) > 0 {
		return val, nil
	}

	return "", fmt.Errorf("required env var `%s` is not set", name)
}
