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
	spec := validSpec()
	spec.Services[0].Name = "jkhaldkfjha760-ahdfkj-lahdfklahsd-kahfdkah"

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		errs[0].Error(),
		fmt.Sprintf("Service name length (%s) must be between 2 and 32 characters long", spec.Services[0].Name),
	)
	suite.Equal(
		errs[1].Error(),
		fmt.Sprintf("Service name (%s) does not match regex ^[a-z][a-z0-9-]{0,30}[a-z0-9]$", spec.Services[0].Name),
	)
}

func (suite *AppSpecLintSuite) TestServiceNeedsAtLeastOneSource() {
	spec := validSpec()
	spec.Name = "hello-world"
	spec.Services[0].GitHub = nil

	errs := spec.Validate()

	suite.Len(errs, 1)
	suite.Equal(
		errs[0].Error(),
		"Service source for hello-world-svc must be one of git, github, gitlab or image",
	)
}

func (suite *AppSpecLintSuite) TestServiceInvalidGitHubSource() {
	spec := validSpec()
	spec.Services[0].GitHub.Branch = ""
	spec.Services[0].GitHub.Repo = "renehernandez_appfile"

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		errs[0].Error(),
		fmt.Sprintf("Github branch for %s service cannot be empty", spec.Services[0].Name),
	)
	suite.Equal(
		errs[1].Error(),
		fmt.Sprintf("Github repo for %s service does not match regex ^[^/]+/[^/]+$", spec.Services[0].Name),
	)
}

func (suite *AppSpecLintSuite) TestServiceInvalidGitLabSource() {
	spec := validSpec()
	svc := spec.Services[0]
	svc.GitHub = nil
	svc.GitLab = &godo.GitLabSourceSpec{
		Repo: "renehernandez_appfile",
	}

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		errs[0].Error(),
		fmt.Sprintf("GitLab branch for %s service cannot be empty", spec.Services[0].Name),
	)
	suite.Equal(
		errs[1].Error(),
		fmt.Sprintf("GitLab repo for %s service does not match regex ^[^/]+/[^/]+$", spec.Services[0].Name),
	)
}

func (suite *AppSpecLintSuite) TestServiceInvalidGitSource() {
	spec := validSpec()
	svc := spec.Services[0]
	svc.GitHub = nil
	svc.Git = &godo.GitSourceSpec{}

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		errs[0].Error(),
		fmt.Sprintf("Git branch for %s service cannot be empty", spec.Services[0].Name),
	)
	suite.Equal(
		errs[1].Error(),
		fmt.Sprintf("Repo clone URL for %s service cannot be empty", spec.Services[0].Name),
	)
}

func (suite *AppSpecLintSuite) TestServiceInvalidEmptyImageSource() {
	spec := validSpecWithImageSource()
	spec.Services[0].Image = &godo.ImageSourceSpec{}

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		errs[0].Error(),
		fmt.Sprintf("Image registry type for %s service is invalid", spec.Services[0].Name),
	)
	suite.Equal(
		errs[1].Error(),
		fmt.Sprintf("Image repository for %s service cannot be empty", spec.Services[0].Name),
	)
}

func (suite *AppSpecLintSuite) TestServiceInvalidDOCRImageSource() {
	spec := validSpecWithImageSource()
	spec.Services[0].Image.Registry = "custom"
	spec.Services[0].Image.Repository = ""

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		errs[0].Error(),
		fmt.Sprintf("Image registry for %s service of type %s must be empty",
			spec.Services[0].Name,
			RegistryTypes.DOCR,
		),
	)
	suite.Equal(
		errs[1].Error(),
		fmt.Sprintf("Image repository for %s service cannot be empty", spec.Services[0].Name),
	)
}

func (suite *AppSpecLintSuite) TestServiceInvalidDockerHubImageSource() {
	spec := validSpecWithImageSource()
	svc := spec.Services[0]
	svc.Image.RegistryType = RegistryTypes.DOCKER_HUB
	svc.Image.Registry = ""
	svc.Image.Repository = ""

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		errs[0].Error(),
		fmt.Sprintf("Image registry for %s service of type %s cannot be empty",
			svc.Name,
			RegistryTypes.DOCKER_HUB,
		),
	)
	suite.Equal(
		errs[1].Error(),
		fmt.Sprintf("Image repository for %s service cannot be empty", svc.Name),
	)
}

func validSpecWithImageSource() *AppSpec {
	spec := validSpec()
	spec.Services[0].GitHub = nil
	spec.Services[0].Image = &godo.ImageSourceSpec{
		RegistryType: RegistryTypes.DOCR,
		Repository:   "my_repository",
	}

	return spec
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
