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
	suite.Equal("Spec name length (a) must be between 2 and 32 characters long", errs[0].Error())
	suite.Equal("Spec name (a) does not match regex ^[a-z][a-z0-9-]{0,30}[a-z0-9]$", errs[1].Error())
}

func (suite *AppSpecLintSuite) TestSpecNameLengthCannotBeLongerThan32() {
	spec := NewAppSpec()
	spec.Name = "jkhaldkfjha760-ahdfkj-lahdfklahsd-kahfdkah"
	spec.SetDefaultValues()

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		fmt.Sprintf("Spec name length (%s) must be between 2 and 32 characters long", spec.Name),
		errs[0].Error(),
	)
	suite.Equal(
		fmt.Sprintf("Spec name (%s) does not match regex ^[a-z][a-z0-9-]{0,30}[a-z0-9]$", spec.Name),
		errs[1].Error(),
	)
}

type ServiceSpecLintSuite struct {
	suite.Suite
}

func (suite *ServiceSpecLintSuite) TestNameInvalid() {
	spec := validSpec()
	spec.Services[0].Name = "jkhaldkfjha760-ahdfkj-lahdfklahsd-kahfdkah"

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		fmt.Sprintf("Service name length (%s) must be between 2 and 32 characters long", spec.Services[0].Name),
		errs[0].Error(),
	)
	suite.Equal(
		fmt.Sprintf("Service name (%s) does not match regex ^[a-z][a-z0-9-]{0,30}[a-z0-9]$", spec.Services[0].Name),
		errs[1].Error(),
	)
}

func (suite *ServiceSpecLintSuite) TestNeedsOneSource() {
	spec := validSpec()
	spec.Name = "hello-world"
	spec.Services[0].GitHub = nil

	errs := spec.Validate()

	suite.Len(errs, 1)
	suite.Equal(
		fmt.Sprintf("Service %s source must be exactly one of git, github, gitlab or image", spec.Services[0].Name),
		errs[0].Error(),
	)
}

func (suite *ServiceSpecLintSuite) TestCannotHaveMoreThanOneSource() {
	spec := validSpec()
	spec.Name = "hello-world"
	spec.Services[0].GitLab = &godo.GitLabSourceSpec{}

	errs := spec.Validate()

	suite.Len(errs, 1)
	suite.Equal(
		fmt.Sprintf("Service %s source must be exactly one of git, github, gitlab or image", spec.Services[0].Name),
		errs[0].Error(),
	)
}

func (suite *ServiceSpecLintSuite) TestInvalidGitHubSource() {
	spec := validSpec()
	spec.Services[0].GitHub.Branch = ""
	spec.Services[0].GitHub.Repo = "renehernandez_appfile"

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		fmt.Sprintf("Github branch for %s service cannot be empty", spec.Services[0].Name),
		errs[0].Error(),
	)
	suite.Equal(
		fmt.Sprintf("Github repo for %s service does not match regex ^[^/]+/[^/]+$", spec.Services[0].Name),
		errs[1].Error(),
	)
}

func (suite *ServiceSpecLintSuite) TestInvalidGitLabSource() {
	spec := validSpec()
	svc := spec.Services[0]
	svc.GitHub = nil
	svc.GitLab = &godo.GitLabSourceSpec{
		Repo: "renehernandez_appfile",
	}

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		fmt.Sprintf("GitLab branch for %s service cannot be empty", spec.Services[0].Name),
		errs[0].Error(),
	)
	suite.Equal(
		fmt.Sprintf("GitLab repo for %s service does not match regex ^[^/]+/[^/]+$", spec.Services[0].Name),
		errs[1].Error(),
	)
}

func (suite *ServiceSpecLintSuite) TestInvalidGitSource() {
	spec := validSpec()
	svc := spec.Services[0]
	svc.GitHub = nil
	svc.Git = &godo.GitSourceSpec{}

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		fmt.Sprintf("Git branch for %s service cannot be empty", spec.Services[0].Name),
		errs[0].Error(),
	)
	suite.Equal(
		fmt.Sprintf("Repo clone URL for %s service cannot be empty", spec.Services[0].Name),
		errs[1].Error(),
	)
}

func (suite *ServiceSpecLintSuite) TestInvalidEmptyImageSource() {
	spec := validSpecWithImageSource()
	spec.Services[0].Image = &godo.ImageSourceSpec{}

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		fmt.Sprintf("Image registry type for %s service is invalid", spec.Services[0].Name),
		errs[0].Error(),
	)
	suite.Equal(
		fmt.Sprintf("Image repository for %s service cannot be empty", spec.Services[0].Name),
		errs[1].Error(),
	)
}

func (suite *ServiceSpecLintSuite) TestInvalidDOCRImageSource() {
	spec := validSpecWithImageSource()
	spec.Services[0].Image.Registry = "custom"
	spec.Services[0].Image.Repository = ""

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		fmt.Sprintf("Image registry for %s service of type %s must be empty",
			spec.Services[0].Name,
			RegistryTypes.DOCR,
		),
		errs[0].Error(),
	)
	suite.Equal(
		fmt.Sprintf("Image repository for %s service cannot be empty", spec.Services[0].Name),
		errs[1].Error(),
	)
}

func (suite *ServiceSpecLintSuite) TestInvalidDockerHubImageSource() {
	spec := validSpecWithImageSource()
	svc := spec.Services[0]
	svc.Image.RegistryType = RegistryTypes.DOCKER_HUB
	svc.Image.Registry = ""
	svc.Image.Repository = ""

	errs := spec.Validate()

	suite.Len(errs, 2)
	suite.Equal(
		fmt.Sprintf("Image registry for %s service of type %s cannot be empty",
			svc.Name,
			RegistryTypes.DOCKER_HUB,
		),
		errs[0].Error(),
	)
	suite.Equal(
		fmt.Sprintf("Image repository for %s service cannot be empty", svc.Name),
		errs[1].Error(),
	)
}

func (suite *ServiceSpecLintSuite) TestInvalidEnvs() {
	spec := validSpecWithImageSource()
	svc := spec.Services[0]
	svc.Envs = []*godo.AppVariableDefinition{
		{
			Key:   "ase@/;afajd",
			Scope: "INVALID_SCOPE",
			Type:  "Invalid type",
		},
	}

	errs := spec.Validate()

	suite.Len(errs, 3)

	suite.Equal(
		fmt.Sprintf("Service %s env key %s does not match regex %s",
			svc.Name,
			"ase@/;afajd",
			"^[_A-Za-z][_A-Za-z0-9]*$",
		),
		errs[0].Error(),
	)
	suite.Equal(
		fmt.Sprintf(
			"Service %s env scope 'INVALID_SCOPE' is not valid. Must be one of [RUN_TIME BUILD_TIME RUN_AND_BUILD_TIME]",
			svc.Name,
		),
		errs[1].Error(),
	)
	suite.Equal(
		fmt.Sprintf(
			"Service %s env type 'Invalid type' is not valid. Must be one of [GENERAL SECRET]",
			svc.Name,
		),
		errs[2].Error(),
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

func TestAppSpecLintSuite(t *testing.T) {
	suite.Run(t, &AppSpecLintSuite{})
}

func TestServiceSpecLintSuite(t *testing.T) {
	suite.Run(t, &ServiceSpecLintSuite{})
}
