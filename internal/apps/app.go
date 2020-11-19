package apps

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
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

type AppService struct {
	client *godo.Client
}

func newService(token string) AppService {
	return AppService{
		client: godo.NewFromToken(token),
	}
}

func (svc *AppService) ListApps() ([]*godo.App, error) {
	list := []*godo.App{}
	ctx := context.TODO()

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}

	for {
		apps, resp, err := svc.client.Apps.List(ctx, opt)
		if err != nil {
			return []*godo.App{}, err
		}

		// append the current page's droplets to our list
		list = append(list, apps...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return []*godo.App{}, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return list, nil
}

func (svc *AppService) FindByName(appName string) (*godo.App, error) {
	apps, err := svc.ListApps()
	if err != nil {
		return &godo.App{}, err
	}

	for _, app := range apps {
		if app.Spec.Name == appName {
			return app, nil
		}
	}

	return &godo.App{}, errors.New("App with name %s not found")
}

func (svc *AppService) Create(app *godo.App) error {
	ctx := context.TODO()
	request := &godo.AppCreateRequest{Spec: app.Spec}

	_, _, err := svc.client.Apps.Create(ctx, request)
	if err != nil {
		return errors.Wrapf(err, "Failed to create new app from spec %s", app.Spec.Name)
	}

	return nil
}

func (svc *AppService) Update(local *godo.App, remote *godo.App) error {
	ctx := context.TODO()
	request := &godo.AppUpdateRequest{Spec: local.Spec}

	_, _, err := svc.client.Apps.Update(ctx, remote.ID, request)
	if err != nil {
		return errors.Wrapf(err, "Failed to update app from spec %s", local.Spec.Name)
	}

	return nil
}

func (svc *AppService) Destroy(app *godo.App) error {
	ctx := context.TODO()

	_, err := svc.client.Apps.Delete(ctx, app.ID)
	if err != nil {
		return errors.Wrapf(err, "Failed to update app from spec %s", app.Spec.Name)
	}

	return nil
}
