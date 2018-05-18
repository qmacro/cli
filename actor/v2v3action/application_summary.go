package v2v3action

import "code.cloudfoundry.org/cli/actor/v2action"

type ApplicationSummary struct {
	v2action.ApplicationSummary
	Buildpacks []string
}

func (appSummary ApplicationSummary) CalculatedBuildpacks() []string {
	if len(appSummary.Buildpacks) == 0 {
		return []string{appSummary.CalculatedBuildpack()}
	}
	return appSummary.Buildpacks
}

func (actor Actor) GetApplicationSummaryByNameAndSpace(name string, spaceGUID string) (ApplicationSummary, Warnings, error) {
	var allWarnings Warnings

	v2appSummary, v2warnings, err := actor.V2Actor.GetApplicationSummaryByNameAndSpace(name, spaceGUID)
	allWarnings = append(allWarnings, v2warnings...)
	if err != nil {
		return ApplicationSummary{}, allWarnings, err
	}

	v3app, v3warnings, err := actor.V3Actor.GetApplicationByNameAndSpace(name, spaceGUID)
	allWarnings = append(allWarnings, v3warnings...)
	if err != nil {
		return ApplicationSummary{}, allWarnings, err
	}

	return ApplicationSummary{
		ApplicationSummary: v2appSummary,
		Buildpacks:         v3app.LifecycleBuildpacks,
	}, allWarnings, err
}
