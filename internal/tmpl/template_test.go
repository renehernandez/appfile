package tmpl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type templateSuite struct {
	suite.Suite
}

func (s *templateSuite) TestWithData() {
	valuesYamlContent := `foo:
  bar: {{ .foo.bar }}
`
	expected := `foo:
  bar: FOO_BAR
`
	data := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "FOO_BAR",
		},
	}
	buf, err := RenderTemplateToBuffer(valuesYamlContent, data)

	s.NoError(err, "unexpected error")

	actual := buf.String()

	s.Equal(expected, actual)
}

func (s *templateSuite) TestRequiredEnvFailsIfNotSet() {
	valuesYamlContent := `foo:
  bar: {{ requiredEnv "FOO_VALUE" }}
`
	_, err := RenderTemplateToBuffer(valuesYamlContent)

	s.Error(err, "required env var `FOO_VALUE` is not set")
}

func (s *templateSuite) TestRequiredEnv() {
	valuesYamlContent := `foo:
  bar: {{ requiredEnv "FOO_VALUE" }}
`
	os.Setenv("FOO_VALUE", "foo")
	buf, err := RenderTemplateToBuffer(valuesYamlContent)

	s.NoError(err, "unexpected error")

	expected := `foo:
  bar: foo
`
	s.Equal(expected, buf.String())
	os.Unsetenv("FOO_VALUE")
}

func TestTemplateSuite(t *testing.T) {
	suite.Run(t, &templateSuite{})
}
