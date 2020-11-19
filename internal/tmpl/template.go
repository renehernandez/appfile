package tmpl

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
)

func newTemplate() *template.Template {
	funcMap := sprig.TxtFuncMap()

	for name, f := range createFuncMap() {
		funcMap[name] = f
	}

	return template.New("stringTemplate").Funcs(funcMap)
}

func RenderTemplateToBuffer(s string, data ...interface{}) (*bytes.Buffer, error) {
	tpl, err := newTemplate().Parse(s)

	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	var d interface{}

	if len(data) > 0 {
		d = data[0]
	}

	if err = tpl.Execute(&buffer, d); err != nil {
		return &buffer, errors.Wrapf(err, "Failed to execute template with data %++v", d)
	}

	return &buffer, nil
}
