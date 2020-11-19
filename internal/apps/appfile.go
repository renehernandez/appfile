package apps

import (
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/renehernandez/appfile/internal/env"
	"github.com/renehernandez/appfile/internal/log"
)

type Appfile struct {
	Spec *AppfileSpec

	Environment *env.Environment
	Apps        []*godo.App
}

func NewAppfileFromSpec(spec *AppfileSpec, envName string) (*Appfile, error) {
	env, err := spec.ReadEnvironment(envName)
	if err != nil {
		return &Appfile{}, err
	}

	apps, err := spec.readApps(env)
	if err != nil {
		return &Appfile{}, err
	}

	return &Appfile{
		Spec:        spec,
		Environment: env,
		Apps:        apps,
	}, nil
}

func (appfile *Appfile) Sync(token string) error {
	remoteApps, err := appfile.readAppsFromRemote(token)
	if err != nil {
		return err
	}

	svc := newService(token)

	for _, localApp := range appfile.Apps {
		log.Infof("Syncing app %s", localApp.Spec.Name)
		remoteApp, ok := remoteApps[localApp.Spec.Name]
		if !ok {
			err = svc.Create(localApp)
		} else {
			err = svc.Update(localApp, remoteApp)
		}

		if err != nil {
			return err
		}
		log.Infof("App %s synced successfully", localApp.Spec.Name)
	}

	return nil
}

func (appfile *Appfile) Destroy(token string) error {
	remoteApps, err := appfile.readAppsFromRemote(token)
	if err != nil {
		return err
	}

	svc := newService(token)
	remoteList := []*godo.App{}

	for _, localApp := range appfile.Apps {
		remoteApp, ok := remoteApps[localApp.Spec.Name]
		if !ok {
			return fmt.Errorf("No app to destroy with name %s", localApp.Spec.Name)
		}

		remoteList = append(remoteList, remoteApp)
	}

	for _, app := range remoteList {
		log.Infof("Destroying app %s", app.Spec.Name)
		err := svc.Destroy(app)
		if err != nil {
			return err
		}
		log.Infof("App %s destroyed successfully", app.Spec.Name)
	}

	return nil
}

func (appfile *Appfile) Diff(token string) ([]*AppDiff, error) {
	remoteApps, err := appfile.readAppsFromRemote(token)
	if err != nil {
		return []*AppDiff{}, err
	}

	appDiffs := []*AppDiff{}

	for _, localApp := range appfile.Apps {
		remoteApp, ok := remoteApps[localApp.Spec.Name]
		if !ok {
			remoteApp = &godo.App{}
		}

		appDiffs = append(appDiffs, &AppDiff{
			Name:      localApp.Spec.Name,
			localApp:  localApp,
			remoteApp: remoteApp,
		})
	}

	return appDiffs, nil
}

func (appfile *Appfile) readAppsFromRemote(token string) (map[string]*godo.App, error) {
	svc := newService(token)

	remoteApps, err := svc.ListApps()
	if err != nil {
		return map[string]*godo.App{}, errors.Wrap(err, "Failed to get apps data from DigitalOcean")
	}

	mapping := map[string]*godo.App{}

	for _, app := range remoteApps {
		mapping[app.Spec.Name] = app
	}

	return mapping, nil
}
