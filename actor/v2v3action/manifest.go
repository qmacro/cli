package v2v3action

import (
	"errors"

	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/cf/manifest"
)

type ManifestV2Actor interface {
	CreateManifestApplication(string, string) (manifest.Application, v2action.Warnings, error)
}

type ManifestV3Actor interface {
	GetApplicationByNameAndSpace(string, string) (manifest.Application, v2action.Warnings, error)
}

func (actor *Actor) CreateApplicationManifestByNameAndSpace(appName string, appSpace string) (Warnings, error) {

	actor.V2Actor.CreateManifestApplication(appName, appSpace)

	if actor.V3Actor.CloudControllerAPIVersion() {
		actor.V3Actor.GetApplicationByNameAndSpace(appName, appSpace)
	}

	return Warnings{"v2-action-warnings"}, errors.New("spaghetti")
}
