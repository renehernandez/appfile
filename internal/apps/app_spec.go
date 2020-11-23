package apps

import (
	"github.com/digitalocean/godo"
)

type AppSpec struct {
	godo.AppSpec
}

func (spec *AppSpec) SetDefaultValues() {
	if len(spec.StaticSites) > 0 {
		for _, siteSpec := range spec.StaticSites {
			if len(siteSpec.Routes) == 0 {
				siteSpec.Routes = append(siteSpec.Routes, &godo.AppRouteSpec{
					Path: "/",
				})
			}
		}
	}

	if len(spec.Services) > 0 {
		for _, svcSpec := range spec.Services {
			if len(svcSpec.InternalPorts) == 0 && len(svcSpec.Routes) == 0 {
				svcSpec.Routes = append(svcSpec.Routes, &godo.AppRouteSpec{
					Path: "/",
				})
			}
		}
	}
}
