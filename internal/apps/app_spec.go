package apps

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/renehernandez/appfile/internal/log"
)

type AppSpec struct {
	*godo.AppSpec

	validator *specValidator
}

func NewAppSpec() *AppSpec {
	spec := &AppSpec{
		AppSpec: &godo.AppSpec{},
	}

	spec.validator = newSpecValidator(spec)

	return spec
}

func (spec *AppSpec) Validate() []error {
	return spec.validator.Validate()
}

func (spec *AppSpec) SetDefaultValues() {
	log.Debugf("Setting default values for %s spec", spec.Name)
	for _, siteSpec := range spec.StaticSites {
		spec.setEnvVarsDefaults(siteSpec.Envs)
		if len(siteSpec.Routes) == 0 {
			siteSpec.Routes = append(siteSpec.Routes, &godo.AppRouteSpec{
				Path: "/",
			})
		}
	}

	for _, svcSpec := range spec.Services {
		spec.setEnvVarsDefaults(svcSpec.Envs)

		if svcSpec.InstanceCount == 0 {
			svcSpec.InstanceCount = 1
		}

		if svcSpec.InstanceSizeSlug == "" {
			svcSpec.InstanceSizeSlug = "basic-xxs"
		}

		if len(svcSpec.InternalPorts) == 0 && len(svcSpec.Routes) == 0 {
			svcSpec.Routes = append(svcSpec.Routes, &godo.AppRouteSpec{
				Path: "/",
			})
		}
	}

	for _, jobSpec := range spec.Jobs {
		spec.setEnvVarsDefaults(jobSpec.Envs)

		if jobSpec.Kind == "" {
			jobSpec.Kind = "POST_DEPLOY"
		}

		if jobSpec.InstanceCount == 0 {
			jobSpec.InstanceCount = 1
		}

		if jobSpec.InstanceSizeSlug == "" {
			jobSpec.InstanceSizeSlug = "basic-xxs"
		}
	}

	for _, workerSpec := range spec.Workers {
		spec.setEnvVarsDefaults(workerSpec.Envs)

		if workerSpec.InstanceCount == 0 {
			workerSpec.InstanceCount = 1
		}

		if workerSpec.InstanceSizeSlug == "" {
			workerSpec.InstanceSizeSlug = "basic-xxs"
		}
	}
}

func (spec *AppSpec) setEnvVarsDefaults(envs []*godo.AppVariableDefinition) {
	for _, envVar := range envs {
		if envVar.Scope == "" {
			envVar.Scope = "RUN_AND_BUILD_TIME"
		}
	}
}

type specValidator struct {
	Spec *AppSpec
}

type regexes struct {
	Name   *regexp.Regexp
	Domain *regexp.Regexp
	Repo   *regexp.Regexp
}

type registryTypes struct {
	DOCR       godo.ImageSourceSpecRegistryType
	DOCKER_HUB godo.ImageSourceSpecRegistryType
}

var (
	SpecRegexes = &regexes{
		Name:   regexp.MustCompile(`^[a-z][a-z0-9-]{0,30}[a-z0-9]$`),
		Domain: regexp.MustCompile(`^((xn--)?[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]{2,}\.?$`),
		Repo:   regexp.MustCompile(`^[^/]+/[^/]+$`),
	}

	RegistryTypes = &registryTypes{
		DOCR:       "DOCR",
		DOCKER_HUB: "DOCKER_HUB",
	}
)

func newSpecValidator(as *AppSpec) *specValidator {
	return &specValidator{
		Spec: as,
	}
}

func (sv *specValidator) Validate() []error {
	errs := []error{}

	errs = append(errs, sv.validateName(sv.Spec.Name, sv.Spec.Name, "Spec")...)
	errs = append(errs, sv.validateServices()...)

	return errs
}

func (sv *specValidator) validateServices() []error {
	errs := []error{}

	for _, svc := range sv.Spec.Services {
		errs = append(errs, sv.validateName(sv.Spec.Name, svc.Name, "Service")...)
		errs = append(errs, (&sourceSpecValidator{
			SpecName:  sv.Spec.Name,
			Name:      svc.Name,
			FieldType: "Service",
			Git:       svc.Git,
			GitHub:    svc.GitHub,
			GitLab:    svc.GitLab,
			Image:     svc.Image,
		}).validate()...)
	}

	return errs
}

func (sv *specValidator) validateName(specName string, name string, fieldType string) []error {
	errs := []error{}

	nameLength := len(name)

	if nameLength < 2 || nameLength > 32 {
		errs = append(errs, fmt.Errorf("%s name length (%s) must be between 2 and 32 characters long", fieldType, name))
	}

	if !SpecRegexes.Name.MatchString(name) {
		errs = append(errs, fmt.Errorf("%s name (%s) does not match regex %s", fieldType, name, SpecRegexes.Name))
	}

	return errs
}

type sourceSpecValidator struct {
	SpecName  string
	Name      string
	FieldType string
	Git       *godo.GitSourceSpec
	GitHub    *godo.GitHubSourceSpec
	GitLab    *godo.GitLabSourceSpec
	Image     *godo.ImageSourceSpec
}

func (sources *sourceSpecValidator) validate() []error {
	errs := []error{}
	sourcesConfigured := 0
	if sources.Git != nil {
		sourcesConfigured++
	}
	if sources.GitHub != nil {
		sourcesConfigured++
	}
	if sources.GitLab != nil {
		sourcesConfigured++
	}
	if sources.Image != nil {
		sourcesConfigured++
	}

	if sourcesConfigured == 0 {
		errs = append(errs,
			fmt.Errorf(
				"%s source for %s must be one of git, github, gitlab or image",
				sources.FieldType,
				sources.Name,
			),
		)

		return errs
	} else if sourcesConfigured > 1 {
		errs = append(errs,
			fmt.Errorf(
				"%s source for %s can only be one of git, github, gitlab or image",
				sources.FieldType,
				sources.Name,
			),
		)

		return errs
	}

	if sources.Git != nil {
		if sources.Git.Branch == "" {
			errs = append(errs, fmt.Errorf("Git branch for %s %s cannot be empty",
				sources.Name,
				strings.ToLower(sources.FieldType),
			))
		}

		if sources.Git.RepoCloneURL == "" {
			errs = append(errs, fmt.Errorf("Repo clone URL for %s %s cannot be empty",
				sources.Name,
				strings.ToLower(sources.FieldType),
			))
		}
	}

	if sources.GitHub != nil {
		if sources.GitHub.Branch == "" {
			errs = append(errs, fmt.Errorf("Github branch for %s %s cannot be empty",
				sources.Name,
				strings.ToLower(sources.FieldType),
			))
		}

		if !SpecRegexes.Repo.MatchString(sources.GitHub.Repo) {
			errs = append(errs, fmt.Errorf("Github repo for %s %s does not match regex %s",
				sources.Name,
				strings.ToLower(sources.FieldType),
				SpecRegexes.Repo,
			))
		}
	}

	if sources.GitLab != nil {
		if sources.GitLab.Branch == "" {
			errs = append(errs, fmt.Errorf("GitLab branch for %s %s cannot be empty",
				sources.Name,
				strings.ToLower(sources.FieldType),
			))
		}

		if !SpecRegexes.Repo.MatchString(sources.GitLab.Repo) {
			errs = append(errs, fmt.Errorf("GitLab repo for %s %s does not match regex %s",
				sources.Name,
				strings.ToLower(sources.FieldType),
				SpecRegexes.Repo,
			))
		}
	}

	if sources.Image != nil {
		if sources.Image.RegistryType != RegistryTypes.DOCR && sources.Image.RegistryType != RegistryTypes.DOCKER_HUB {
			errs = append(errs, fmt.Errorf("Image registry type for %s %s is invalid",
				sources.Name,
				strings.ToLower(sources.FieldType),
			))
		}

		if sources.Image.RegistryType == RegistryTypes.DOCR && sources.Image.Registry != "" {
			errs = append(errs, fmt.Errorf("Image registry for %s %s of type %s must be empty",
				sources.Name,
				strings.ToLower(sources.FieldType),
				sources.Image.RegistryType,
			))
		}

		if sources.Image.RegistryType == RegistryTypes.DOCKER_HUB && sources.Image.Registry == "" {
			errs = append(errs, fmt.Errorf("Image registry for %s %s of type cannot be empty",
				sources.Name,
				strings.ToLower(sources.FieldType),
			))
		}

		if sources.Image.Repository == "" {
			errs = append(errs, fmt.Errorf("Image repository for %s %s cannot be empty",
				sources.Name,
				strings.ToLower(sources.FieldType),
			))
		}
	}

	return errs
}
