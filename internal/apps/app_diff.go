package apps

import (
	"github.com/digitalocean/godo"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type AppDiff struct {
	Name string

	localApp  *godo.App
	remoteApp *godo.App
}

func (diff *AppDiff) CalculateDiff() ([]diffmatchpatch.Diff, error) {
	var localSpec, remoteSpec string
	var err error

	if localSpec, err = appSpecToString(diff.localApp); err != nil {
		return []diffmatchpatch.Diff{}, err
	}

	if remoteSpec, err = appSpecToString(diff.remoteApp); err != nil {
		return []diffmatchpatch.Diff{}, err
	}

	dmp := diffmatchpatch.New()

	fileAdmp, fileBdmp, dmpStrings := dmp.DiffLinesToChars(remoteSpec, localSpec)
	diffs := dmp.DiffMain(fileAdmp, fileBdmp, false)
	diffs = dmp.DiffCharsToLines(diffs, dmpStrings)
	diffs = dmp.DiffCleanupSemantic(diffs)

	return diffs, nil
}

func appSpecToString(app *godo.App) (string, error) {
	if app.Spec == nil {
		return "", nil
	}

	b, err := yaml.Marshal(app.Spec)
	if err != nil {
		return "", errors.Wrapf(err, "Error converting spec to json string for app %s", app.Spec.Name)
	}

	return string(b), nil
}
