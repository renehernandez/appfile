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
	Spec     *AppfileSpec
	AppSpecs []*AppSpec
	State    *StateData

	token string
}

func NewAppfileFromAppSpec(spec *AppSpec, token string) (*Appfile, error) {
	spec.SetDefaultValues()

	return &Appfile{
		Spec: &AppfileSpec{},
		AppSpecs: []*AppSpec{
			spec,
		},
		State: &StateData{},
		token: token,
	}, nil
}

func NewAppfileFromSpec(spec *AppfileSpec, envName string, token string) (*Appfile, error) {
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

	appSpecs, err := spec.loadAppSpecs(&state)
	if err != nil {
		return &Appfile{}, err
	}

	return &Appfile{
		Spec:     spec,
		State:    &state,
		AppSpecs: appSpecs,
		token:    token,
	}, nil
}

func (appfile *Appfile) Sync() error {
	remoteApps, err := appfile.readAppsFromRemote()
	if err != nil {
		return err
	}

	svc := do.NewAppService(appfile.token)

	for _, appSpec := range appfile.AppSpecs {
		log.Infof("Syncing app %s", appSpec.Name)
		remoteApp, ok := remoteApps[appSpec.Name]
		localApp := &godo.App{Spec: appSpec.AppSpec}
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

func (appfile *Appfile) Destroy() error {
	remoteApps, err := appfile.readAppsFromRemote()
	if err != nil {
		return err
	}

	appSvc := do.NewAppService(appfile.token)
	domainSvc := do.NewDomainService(appfile.token)
	remoteList := []*godo.App{}

	for _, appSpec := range appfile.AppSpecs {
		remoteApp, ok := remoteApps[appSpec.Name]
		if !ok {
			return fmt.Errorf("No app to destroy with name %s", appSpec.Name)
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

func (appfile *Appfile) Diff() ([]*AppDiff, error) {
	remoteApps, err := appfile.readAppsFromRemote()
	if err != nil {
		return []*AppDiff{}, err
	}
	appDiffs := []*AppDiff{}

	for _, appSpec := range appfile.AppSpecs {
		remoteApp, ok := remoteApps[appSpec.Name]
		if !ok {
			remoteApp = &godo.App{}
		}

		appDiffs = append(appDiffs, &AppDiff{
			Name:       appSpec.Name,
			localSpec:  appSpec.AppSpec,
			remoteSpec: remoteApp.Spec,
		})
	}

	return appDiffs, nil
}

func (appfile *Appfile) Status() ([]*AppStatus, error) {
	remoteApps, err := appfile.readAppsFromRemote()
	if err != nil {
		return []*AppStatus{}, err
	}

	appsStatus := []*AppStatus{}

	for _, appSpec := range appfile.AppSpecs {
		appStatus := &AppStatus{
			Name:         appSpec.Name,
			Status:       DeploymentStatusUnknown,
			URL:          "-",
			UpdatedAt:    "-",
			DeploymentID: "-",
		}

		remoteApp, ok := remoteApps[appSpec.Name]

		if ok {
			if remoteApp.InProgressDeployment != nil {
				appStatus.DeploymentID = remoteApp.InProgressDeployment.ID
				appStatus.Status = DeploymentStatusInProgress
			} else if remoteApp.ActiveDeployment != nil {
				appStatus.DeploymentID = remoteApp.ActiveDeployment.ID
				appStatus.Status = DeploymentStatusDeployed
			}

			if appStatus.Status != DeploymentStatusUnknown {
				appStatus.UpdatedAt = remoteApp.UpdatedAt.String()
				appStatus.URL = remoteApp.LiveDomain
			}

			appsStatus = append(appsStatus, appStatus)
		} else {
			log.Warningf("%s app not found in App Platform", appSpec.Name)
		}
	}

	return appsStatus, nil
}

func (appfile *Appfile) Lint() ([]AppLint, error) {
	lints := []AppLint{}

	remoteApps, err := appfile.readAppsFromRemote()
	if err != nil {
		return []AppLint{}, err
	}

	svc := do.NewAppService(appfile.token)

	for _, appSpec := range appfile.AppSpecs {
		lint := AppLint{
			Name:     appSpec.Name,
			FileName: appSpec.FileName,
			Errors:   appSpec.Validate(),
		}

		lints = append(lints, lint)

		if len(lint.Errors) != 0 {
			continue
		}

		remoteApp, ok := remoteApps[appSpec.Name]
		localApp := &godo.App{
			Spec: appSpec.AppSpec,
		}

		if ok {
			localApp.ID = remoteApp.ID
		}

		lint.Errors = append(lint.Errors, svc.Propose(localApp))
	}

	return lints, nil
}

func (appfile *Appfile) readAppsFromRemote() (map[string]*godo.App, error) {
	log.Debugln("Get apps running in DigitalOcean")
	svc := do.NewAppService(appfile.token)

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
