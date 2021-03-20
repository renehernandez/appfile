package apps

import (
	"fmt"
	"regexp"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/digitalocean/godo"
	"github.com/renehernandez/appfile/internal/log"
)

type regexes struct {
	Name   *regexp.Regexp
	Domain *regexp.Regexp
	Repo   *regexp.Regexp
	EnvKey *regexp.Regexp
}

type registryTypes struct {
	DOCR       godo.ImageSourceSpecRegistryType
	DOCKER_HUB godo.ImageSourceSpecRegistryType
}

type workloadType string

type workloadTypes struct {
	ServiceType    workloadType
	JobType        workloadType
	WorkerType     workloadType
	StaticSiteType workloadType
}

var (
	SpecRegexes = &regexes{
		Name:   regexp.MustCompile(`^[a-z][a-z0-9-]{0,30}[a-z0-9]$`),
		Domain: regexp.MustCompile(`^((xn--)?[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]{2,}\.?$`),
		Repo:   regexp.MustCompile(`^[^/]+/[^/]+$`),
		EnvKey: regexp.MustCompile(`^[_A-Za-z][_A-Za-z0-9]*$`),
	}

	RegistryTypes = &registryTypes{
		DOCR:       "DOCR",
		DOCKER_HUB: "DOCKER_HUB",
	}

	WorkloadTypes = &workloadTypes{
		ServiceType:    "Service",
		JobType:        "Job",
		StaticSiteType: "Static Site",
		WorkerType:     "Worker",
	}
)

type AppSpec struct {
	*godo.AppSpec

	FileName  string
	validator *specValidator
}

func NewAppSpec() *AppSpec {
	spec := &AppSpec{
		AppSpec: &godo.AppSpec{},
	}

	spec.validator = newSpecValidator(spec)

	return spec
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
		if envVar.Type == "" {
			envVar.Type = "GENERAL"
		}
	}
}

func (spec *AppSpec) Validate() []error {
	return spec.validator.Validate()
}

type specValidator struct {
	Spec *AppSpec
}

func newSpecValidator(as *AppSpec) *specValidator {
	return &specValidator{
		Spec: as,
	}
}

func (sv *specValidator) Validate() []error {
	errs := []error{}

	errs = append(errs, validateName(sv.Spec.Name, sv.Spec.Name, "Spec")...)
	errs = append(errs, sv.validateServices()...)
	errs = append(errs, sv.validateWorkers()...)
	errs = append(errs, sv.validateJobs()...)

	return errs
}

func (sv *specValidator) validateServices() []error {
	errs := []error{}

	for _, svc := range sv.Spec.Services {
		errs = append(errs, (&workloadValidator{
			SpecName:     sv.Spec.Name,
			WorkloadName: svc.Name,
			WorkloadType: string(WorkloadTypes.ServiceType),
			SourceSpecValidator: &sourceSpecValidator{
				Git:    svc.Git,
				GitHub: svc.GitHub,
				GitLab: svc.GitLab,
				Image:  svc.Image,
			},
			EnvsSpecValidator: &envsSpecValidator{
				Global: false,
				Envs:   svc.Envs,
			},
			InstanceSpecValidator: &instanceSpecValidator{
				SizeSlug: svc.InstanceSizeSlug,
				Count:    svc.InstanceCount,
			},
		}).validate()...)
	}

	return errs
}

func (sv *specValidator) validateWorkers() []error {
	errs := []error{}

	for _, worker := range sv.Spec.Workers {
		errs = append(errs, (&workloadValidator{
			SpecName:     sv.Spec.Name,
			WorkloadName: worker.Name,
			WorkloadType: string(WorkloadTypes.WorkerType),
			SourceSpecValidator: &sourceSpecValidator{
				Git:    worker.Git,
				GitHub: worker.GitHub,
				GitLab: worker.GitLab,
				Image:  worker.Image,
			},
			EnvsSpecValidator: &envsSpecValidator{
				Global: false,
				Envs:   worker.Envs,
			},
			InstanceSpecValidator: &instanceSpecValidator{
				SizeSlug: worker.InstanceSizeSlug,
				Count:    worker.InstanceCount,
			},
		}).validate()...)
	}

	return errs
}

func (sv *specValidator) validateJobs() []error {
	errs := []error{}

	for _, job := range sv.Spec.Jobs {
		errs = append(errs, (&workloadValidator{
			SpecName:     sv.Spec.Name,
			WorkloadName: job.Name,
			WorkloadType: string(WorkloadTypes.JobType),
			SourceSpecValidator: &sourceSpecValidator{
				Git:    job.Git,
				GitHub: job.GitHub,
				GitLab: job.GitLab,
				Image:  job.Image,
			},
			EnvsSpecValidator: &envsSpecValidator{
				Global: false,
				Envs:   job.Envs,
			},
			InstanceSpecValidator: &instanceSpecValidator{
				SizeSlug: job.InstanceSizeSlug,
				Count:    job.InstanceCount,
			},
		}).validate()...)
	}

	return errs
}

func (sv *specValidator) validateStaticSites() []error {
	errs := []error{}

	for _, site := range sv.Spec.StaticSites {
		errs = append(errs, (&workloadValidator{
			SpecName:     sv.Spec.Name,
			WorkloadName: site.Name,
			WorkloadType: string(WorkloadTypes.StaticSiteType),
			SourceSpecValidator: &sourceSpecValidator{
				Git:    site.Git,
				GitHub: site.GitHub,
				GitLab: site.GitLab,
			},
			EnvsSpecValidator: &envsSpecValidator{
				Global: false,
				Envs:   site.Envs,
			},
		}).validate()...)
	}

	return errs
}

type workloadValidator struct {
	SpecName              string
	WorkloadName          string
	WorkloadType          string
	SourceSpecValidator   *sourceSpecValidator
	EnvsSpecValidator     *envsSpecValidator
	InstanceSpecValidator *instanceSpecValidator
}

func (workload *workloadValidator) validate() []error {
	errs := []error{}

	errs = append(errs, validateName(workload.SpecName, workload.WorkloadName, workload.WorkloadType)...)
	errs = append(errs, workload.SourceSpecValidator.validate(workload.WorkloadName, workload.WorkloadType)...)
	errs = append(errs, workload.EnvsSpecValidator.validate(workload.WorkloadName, workload.WorkloadType)...)
	errs = append(errs, workload.InstanceSpecValidator.validate(strings.ToLower(workload.WorkloadName), workload.WorkloadType)...)

	return errs
}

type sourceSpecValidator struct {
	Git    *godo.GitSourceSpec
	GitHub *godo.GitHubSourceSpec
	GitLab *godo.GitLabSourceSpec
	Image  *godo.ImageSourceSpec
}

func (validator *sourceSpecValidator) validate(name string, fieldType string) []error {
	errs := []error{}
	sourcesConfigured := 0
	if validator.Git != nil {
		sourcesConfigured++
	}
	if validator.GitHub != nil {
		sourcesConfigured++
	}
	if validator.GitLab != nil {
		sourcesConfigured++
	}
	if validator.Image != nil {
		sourcesConfigured++
	}

	if sourcesConfigured == 0 || sourcesConfigured > 1 {
		errs = append(errs,
			fmt.Errorf(
				"%s %s source must be exactly one of git, github, gitlab or image",
				fieldType,
				name,
			),
		)

		return errs
	}

	if validator.Git != nil {
		if validator.Git.Branch == "" {
			errs = append(errs, fmt.Errorf("Git branch for %s %s cannot be empty",
				name,
				strings.ToLower(fieldType),
			))
		}

		if validator.Git.RepoCloneURL == "" {
			errs = append(errs, fmt.Errorf("Repo clone URL for %s %s cannot be empty",
				name,
				strings.ToLower(fieldType),
			))
		}
	}

	if validator.GitHub != nil {
		if validator.GitHub.Branch == "" {
			errs = append(errs, fmt.Errorf("Github branch for %s %s cannot be empty",
				name,
				strings.ToLower(fieldType),
			))
		}

		if !SpecRegexes.Repo.MatchString(validator.GitHub.Repo) {
			errs = append(errs, fmt.Errorf("Github repo for %s %s does not match regex %s",
				name,
				strings.ToLower(fieldType),
				SpecRegexes.Repo,
			))
		}
	}

	if validator.GitLab != nil {
		if validator.GitLab.Branch == "" {
			errs = append(errs, fmt.Errorf("GitLab branch for %s %s cannot be empty",
				name,
				strings.ToLower(fieldType),
			))
		}

		if !SpecRegexes.Repo.MatchString(validator.GitLab.Repo) {
			errs = append(errs, fmt.Errorf("GitLab repo for %s %s does not match regex %s",
				name,
				strings.ToLower(fieldType),
				SpecRegexes.Repo,
			))
		}
	}

	if validator.Image != nil {
		if validator.Image.RegistryType != RegistryTypes.DOCR && validator.Image.RegistryType != RegistryTypes.DOCKER_HUB {
			errs = append(errs, fmt.Errorf("Image registry type for %s %s is invalid",
				name,
				strings.ToLower(fieldType),
			))
		}

		if validator.Image.RegistryType == RegistryTypes.DOCR && validator.Image.Registry != "" {
			errs = append(errs, fmt.Errorf("Image registry for %s %s of type %s must be empty",
				name,
				strings.ToLower(fieldType),
				validator.Image.RegistryType,
			))
		}

		if validator.Image.RegistryType == RegistryTypes.DOCKER_HUB && validator.Image.Registry == "" {
			errs = append(errs, fmt.Errorf("Image registry for %s %s of type %s cannot be empty",
				name,
				strings.ToLower(fieldType),
				validator.Image.RegistryType,
			))
		}

		if validator.Image.Repository == "" {
			errs = append(errs, fmt.Errorf("Image repository for %s %s cannot be empty",
				name,
				strings.ToLower(fieldType),
			))
		}
	}

	return errs
}

type envsSpecValidator struct {
	Global bool
	Envs   []*godo.AppVariableDefinition
}

func (validator *envsSpecValidator) validate(name string, fieldType string) []error {
	errs := []error{}
	types := []interface{}{"GENERAL", "SECRET"}
	allowedTypes := mapset.NewSetFromSlice(types)
	scopes := []interface{}{"RUN_TIME", "BUILD_TIME", "RUN_AND_BUILD_TIME"}
	allowedScopes := mapset.NewSetFromSlice(scopes)

	prefixMsg := "Global env"
	if !validator.Global {
		prefixMsg = fmt.Sprintf("%s %s env", fieldType, name)
	}

	for _, envDef := range validator.Envs {
		if !SpecRegexes.EnvKey.MatchString(envDef.Key) {
			errs = append(errs, fmt.Errorf("%s key %s does not match regex %s",
				prefixMsg,
				envDef.Key,
				SpecRegexes.EnvKey,
			))
		}

		if !allowedScopes.Contains(envDef.Scope) {
			errs = append(errs, fmt.Errorf("%s scope '%s' is not valid. Must be one of %s",
				prefixMsg,
				envDef.Scope,
				scopes,
			))
		}

		if !allowedTypes.Contains(envDef.Type) {
			errs = append(errs, fmt.Errorf("%s type '%s' is not valid. Must be one of %s",
				prefixMsg,
				envDef.Type,
				types,
			))
		}
	}

	return errs
}

type instanceSpecValidator struct {
	SizeSlug string
	Count    int64
}

func (validator *instanceSpecValidator) validate(name string, fieldType string) []error {
	errs := []error{}

	if validator.Count < 0 {
		errs = append(errs, fmt.Errorf("Instance count for %s %s cannot be negative",
			fieldType,
			name,
		))
	}

	instanceTypes := mapset.NewSetFromSlice([]interface{}{
		"basic-xxs",
		"basic-xs",
		"basic-s",
		"basic-m",
		"professional-xs",
		"professional-s",
		"professional-m",
		"professional-1l",
		"professional-l",
		"professional-xl",
	})

	if !instanceTypes.Contains(validator.SizeSlug) {
		errs = append(errs, fmt.Errorf("Size slug invalid for %s %s",
			fieldType,
			name,
		))
	}

	return errs
}

func validateName(specName string, name string, fieldType string) []error {
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
