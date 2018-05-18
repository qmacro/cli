package v2v3action

import (
	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/util/manifest"
)

//go:generate counterfeiter . V2Actor

type V2Actor interface {
	CreateApplicationManifestByNameAndSpace(name string, spaceGUID string) (manifest.Application, v2action.Warnings, error)
	GetApplicationSummaryByNameAndSpace(name string, spaceGUID string) (v2action.ApplicationSummary, v2action.Warnings, error)
	GetFeatureFlags() ([]v2action.FeatureFlag, v2action.Warnings, error)
	GetService(serviceGUID string) (v2action.Service, v2action.Warnings, error)
	GetServiceInstanceByNameAndSpace(serviceInstanceName string, spaceGUID string) (v2action.ServiceInstance, v2action.Warnings, error)
	GetServiceInstanceSharedTosByServiceInstance(serviceInstanceGUID string) ([]v2action.ServiceInstanceSharedTo, v2action.Warnings, error)
	GetSpaceByOrganizationAndName(orgGUID string, spaceName string) (v2action.Space, v2action.Warnings, error)
}
