package apps

import (
	"github.com/digitalocean/godo"
	"github.com/renehernandez/appfile/internal/log"
)

type AppSpec struct {
	godo.AppSpec
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

	if spec.Region == "" {
		spec.Region = "nyc"
	}
}

func (spec *AppSpec) setEnvVarsDefaults(envs []*godo.AppVariableDefinition) {
	for _, envVar := range envs {
		if envVar.Scope == "" {
			envVar.Scope = "RUN_AND_BUILD_TIME"
		}
	}
}
