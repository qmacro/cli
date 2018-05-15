package v2v3action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/cli/actor/v2action"
	. "code.cloudfoundry.org/cli/actor/v2v3action"
	"code.cloudfoundry.org/cli/actor/v2v3action/v2v3actionfakes"
	"code.cloudfoundry.org/cli/cf/manifest"
)

var _ = FDescribe("Manifest", func() {
	var (
		actor       *Actor
		fakeV2Actor *v2v3actionfakes.FakeV2Actor
		fakeV3Actor *v2v3actionfakes.FakeV3Actor

		appName  string
		appSpace string

		warnings   Warnings
		executeErr error
	)

	BeforeEach(func() {
		fakeV2Actor = new(v2v3actionfakes.FakeV2Actor)
		fakeV3Actor = new(v2v3actionfakes.FakeV3Actor)

		actor = NewActor(fakeV2Actor, fakeV3Actor)

		appName = "some-app-name"
		appSpace = "some-space-GUID"
	})

	JustBeforeEach(func() {
		warnings, executeErr = actor.CreateApplicationManifestByNameAndSpace(appName, appSpace)
	})

	It("calls v2Actor.CreateManifestApplication with the appName and appSpace", func() {
		Expect(fakeV2Actor.CreateManifestApplicationCallCount()).To(Equal(1))
		appNameArg, appSpaceArg := fakeV2Actor.CreateManifestApplicationArgsForCall(0)
		Expect(appNameArg).To(Equal(appName))
		Expect(appSpaceArg).To(Equal(appSpace))
	})

	Context("when v2Actor.CreateManifestApplication succeeds", func() {
		var application manifest.Application

		BeforeEach(func() {
			fakeV2Actor.CreateManifestApplicationReturns(application, v2action.Warnings{"v2-action-warnings"}, nil)
		})

		Context("when there is a relevant ( >= 3.25) v3 endpoint", func() {
			BeforeEach(func() {
				fakeV3Actor.CloudControllerAPIVersionReturns("3.25")
			})

			It("Calls the v3actor.GetApplicationByNameAndSpace with the appName and appSpace", func() {
				Expect(fakeV3Actor.GetApplicationByNameAndSpaceCallCount()).To(Equal(1))
				appNameArg, appSpaceArg := fakeV3Actor.GetApplicationByNameAndSpaceArgsForCall(0)
				Expect(appNameArg).To(Equal(appName))
				Expect(appSpaceArg).To(Equal(appSpace))
			})
		})

		Context("when there is no relevant ( < 3.25) v3 endpoint", func() {
			BeforeEach(func() {
				fakeV3Actor.CloudControllerAPIVersionReturns("3.24")
			})

			It("does not call the v3actor.GetApplicationByNameAndSpace", func() {
				Expect(fakeV3Actor.GetApplicationByNameAndSpaceCallCount()).To(Equal(0))
			})
		})
	})

	Context("when v2Actor.CreateManifestApplication fails", func() {
		BeforeEach(func() {
			fakeV2Actor.CreateManifestApplicationReturns(manifest.Application{}, v2action.Warnings{"v2-action-warnings"}, errors.New("spaghetti"))
		})

		It("returns warnings and the error", func() {
			Expect(warnings).To(ConsistOf("v2-action-warnings"))
			Expect(executeErr).To(MatchError(errors.New("spaghetti")))
		})
	})
})
