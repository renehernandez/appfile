package tmpl

import (
	"bytes"
	"io/ioutil"

	"github.com/pkg/errors"
)

func RenderFromFile(file string, data ...interface{}) (*bytes.Buffer, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return &bytes.Buffer{}, err
	}

	templatedYaml, err := RenderTemplateToBuffer(string(content), data...)
	if err != nil {
		return &bytes.Buffer{}, errors.Wrapf(err, "Could not templatize %s", file)
	}

	return templatedYaml, nil
}
