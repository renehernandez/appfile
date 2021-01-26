package apps

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/suite"
)

type AppSpecLintSuite struct {
	suite.Suite
}

func (suite *AppSpecLintSuite) TestEmptySpec() {
	spec := NewAppSpec()
	spec.SetDefaultValues()

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(errs[0].Error(), "Spec name length () must be between 2 and 32 characters long")
	suite.Equal(errs[1].Error(), "Spec name () does not match regex ^[a-z][a-z0-9-]{0,30}[a-z0-9]$")
}

func (suite *AppSpecLintSuite) TestSpecNameLengthCannotBeSmallerThan2() {
	spec := NewAppSpec()
	spec.Name = "a"
	spec.SetDefaultValues()

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(errs[0].Error(), "Spec name length (a) must be between 2 and 32 characters long")
	suite.Equal(errs[1].Error(), "Spec name (a) does not match regex ^[a-z][a-z0-9-]{0,30}[a-z0-9]$")
}

func (suite *AppSpecLintSuite) TestSpecNameLengthCannotBeLongerThan32() {
	spec := NewAppSpec()
	spec.Name = "jkhaldkfjha760-ahdfkj-lahdfklahsd-kahfdkah"
	spec.SetDefaultValues()

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		errs[0].Error(),
		fmt.Sprintf("Spec name length (%s) must be between 2 and 32 characters long", spec.Name),
	)
	suite.Equal(
		errs[1].Error(),
		fmt.Sprintf("Spec name (%s) does not match regex ^[a-z][a-z0-9-]{0,30}[a-z0-9]$", spec.Name),
	)
}

func (suite *AppSpecLintSuite) TestServiceNameInvalid() {
	spec := NewAppSpec()
	spec.Name = "hello-world"
	svc := &godo.AppServiceSpec{
		Name: "jkhaldkfjha760-ahdfkj-lahdfklahsd-kahfdkah",
		Git: &godo.GitSourceSpec{
			Branch:       "main",
			RepoCloneURL: "https://example.com",
		},
	}
	spec.Services = []*godo.AppServiceSpec{
		svc,
	}

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		errs[0].Error(),
		fmt.Sprintf("Service name length (%s) must be between 2 and 32 characters long", svc.Name),
	)
	suite.Equal(
		errs[1].Error(),
		fmt.Sprintf("Service name (%s) does not match regex ^[a-z][a-z0-9-]{0,30}[a-z0-9]$", svc.Name),
	)
}

func (suite *AppSpecLintSuite) TestServiceNeedsAtLeastOneSource() {
	spec := NewAppSpec()
	spec.Name = "hello-world"
	svc := &godo.AppServiceSpec{
		Name: "hello-world-svc",
	}
	spec.Services = []*godo.AppServiceSpec{
		svc,
	}

	errs := spec.Validate()

	suite.Len(errs, 1)
	suite.Equal(
		errs[0].Error(),
		"Service source for hello-world-svc must be one of git, github, gitlab or image",
	)
}

func (suite *AppSpecLintSuite) TestServiceGithubSourceEmptyBranch() {
	spec := validSpec()
	spec.Services[0].GitHub.Branch = ""

	errs := spec.Validate()

	suite.Len(errs, 1)
	suite.Equal(
		errs[0].Error(),
		fmt.Sprintf("Github branch for %s service cannot be empty", spec.Services[0].Name),
	)
}

func (suite *AppSpecLintSuite) TestServiceGithubSourceInvalidRepo() {
	spec := validSpec()
	spec.Services[0].GitHub.Repo = "renehernandez_appfile"

	errs := spec.Validate()

	suite.Len(errs, 1)
	suite.Equal(
		errs[0].Error(),
		fmt.Sprintf("Github repo for %s service does not match regex ^[^/]+/[^/]+$", spec.Services[0].Name),
	)
}

func validSpec() *AppSpec {
	spec := NewAppSpec()
	spec.Name = "hello-world"
	svc := &godo.AppServiceSpec{
		Name: "hello-world-svc",
		GitHub: &godo.GitHubSourceSpec{
			Repo:   "renehernandez/appfile",
			Branch: "main",
		},
	}
	spec.Services = []*godo.AppServiceSpec{
		svc,
	}

	return spec
}

func TestTemplateSuite(t *testing.T) {
	suite.Run(t, &AppSpecLintSuite{})
}
