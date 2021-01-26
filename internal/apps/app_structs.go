package apps

import (
	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/yaml.v2"
)

type AppDiff struct {
	Name string

	localSpec  *godo.AppSpec
	remoteSpec *godo.AppSpec
}

func (diff *AppDiff) CalculateDiff() ([]diffmatchpatch.Diff, error) {
	var localYaml, remoteYaml string
	var err error

	if localYaml, err = appSpecToString(diff.localSpec); err != nil {
		return []diffmatchpatch.Diff{}, err
	}

	if remoteYaml, err = appSpecToString(diff.remoteSpec); err != nil {
		return []diffmatchpatch.Diff{}, err
	}

	dmp := diffmatchpatch.New()

	fileAdmp, fileBdmp, dmpStrings := dmp.DiffLinesToChars(remoteYaml, localYaml)
	diffs := dmp.DiffMain(fileAdmp, fileBdmp, false)
	diffs = dmp.DiffCharsToLines(diffs, dmpStrings)
	diffs = dmp.DiffCleanupSemantic(diffs)

	return diffs, nil
}

func appSpecToString(spec *godo.AppSpec) (string, error) {
	if spec == nil {
		return "", nil
	}

	b, err := yaml.Marshal(spec)
	if err != nil {
		return "", errors.Wrapf(err, "Error converting spec to json string for app %s", spec.Name)
	}

	return string(b), nil
}

type AppStatus struct {
	Name         string
	Status       DeploymentStatus
	DeploymentID string
	UpdatedAt    string
	URL          string
}

type DeploymentStatus string

const (
	DeploymentStatusUnknown    DeploymentStatus = "unknown"
	DeploymentStatusDeployed   DeploymentStatus = "deployed"
	DeploymentStatusInProgress DeploymentStatus = "in progress"
)

type AppLint struct {
	Name   string
	Errors []error
}
