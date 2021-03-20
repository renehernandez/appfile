package do

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
)

type AppService struct {
	client *godo.Client
}

func NewAppService(token string) *AppService {
	return &AppService{
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

func (svc *AppService) ListInstancesSizes() ([]*godo.AppInstanceSize, error) {
	ctx := context.TODO()

	sizes, _, err := svc.client.Apps.ListInstanceSizes(ctx)
	if err != nil {
		return []*godo.AppInstanceSize{}, err
	}

	return sizes, nil
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
		return errors.Wrapf(err, "Failed to delete app %s", app.Spec.Name)
	}

	return nil
}

func (svc *AppService) Propose(app *godo.App) error {
	ctx := context.TODO()
	request := &godo.AppProposeRequest{
		Spec: app.Spec,
	}

	if app.ID != "" {
		request.AppID = app.ID
	}

	_, _, err := svc.client.Apps.Propose(ctx, request)

	return err
}
