package apps

import (
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/renehernandez/appfile/internal/do"
	"github.com/renehernandez/appfile/internal/log"
)

type EnvMetadata struct {
	Name string
}

type StateData struct {
	Environment EnvMetadata
	Values      map[string]interface{}
}

type Appfile struct {
	Spec *AppfileSpec

	State *StateData
	Apps  []*godo.App
}

func NewAppfileFromSpec(spec *AppfileSpec, envName string) (*Appfile, error) {
	env, err := spec.ReadEnvironment(envName)
	if err != nil {
		return &Appfile{}, err
	}

	state := StateData{
		Environment: EnvMetadata{
			Name: env.Name,
		},
		Values: env.Values,
	}

	apps, err := spec.readApps(&state)
	if err != nil {
		return &Appfile{}, err
	}

	return &Appfile{
		Spec:  spec,
		State: &state,
		Apps:  apps,
	}, nil
}

func (appfile *Appfile) Sync(token string) error {
	remoteApps, err := appfile.readAppsFromRemote(token)
	if err != nil {
		return err
	}

	svc := do.NewAppService(token)

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

	appSvc := do.NewAppService(token)
	domainSvc := do.NewDomainService(token)
	remoteList := []*godo.App{}

	for _, localApp := range appfile.Apps {
		remoteApp, ok := remoteApps[localApp.Spec.Name]
		if !ok {
			return fmt.Errorf("No app to destroy with name %s", localApp.Spec.Name)
		}

		remoteList = append(remoteList, remoteApp)
	}

	for _, app := range remoteList {
		log.Debugf("Destroying app %s", app.Spec.Name)
		err := appSvc.Destroy(app)
		if err != nil {
			return err
		}
		log.Infof("App %s destroyed successfully", app.Spec.Name)

		for _, domain := range app.Spec.Domains {
			if domain.Domain != "" && domain.Zone != "" {
				log.Debugf("Deleting %s hostname in %s zone", domain.Domain, domain.Zone)
				err = domainSvc.DeleteRecord(domain)
				if err != nil {
					return err
				}
			}
		}
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

func (appfile *Appfile) List(token string) ([]*AppStatus, error) {
	remoteApps, err := appfile.readAppsFromRemote(token)
	if err != nil {
		return []*AppStatus{}, err
	}

	appsStatus := []*AppStatus{}

	for _, localApp := range appfile.Apps {
		appStatus := &AppStatus{
			Name:         localApp.Spec.Name,
			Status:       DeploymentStatusUnknown,
			URL:          "-",
			DeploymentID: "-",
			UpdatedAt:    "-",
		}

		remoteApp, ok := remoteApps[localApp.Spec.Name]

		if ok {
			appStatus.UpdatedAt = remoteApp.UpdatedAt.String()
			appStatus.URL = remoteApp.LiveDomain

			if remoteApp.InProgressDeployment != nil {
				appStatus.DeploymentID = remoteApp.InProgressDeployment.ID
				appStatus.Status = DeploymentStatusInProgress
			} else {
				appStatus.DeploymentID = remoteApp.ActiveDeployment.ID
				appStatus.Status = DeploymentStatusDeployed
			}
		} else {
			log.Debugf("%s app not found in App Platform", localApp.Spec.Name)
		}

		appsStatus = append(appsStatus, appStatus)
	}

	return appsStatus, nil
}

func (appfile *Appfile) readAppsFromRemote(token string) (map[string]*godo.App, error) {
	log.Debugln("Get apps running in DigitalOcean")
	svc := do.NewAppService(token)

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
