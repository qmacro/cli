package v2v3action_test

import (
	"errors"

	"code.cloudfoundry.org/cli/actor/v2action"
	. "code.cloudfoundry.org/cli/actor/v2v3action"
	"code.cloudfoundry.org/cli/actor/v2v3action/v2v3actionfakes"
	"code.cloudfoundry.org/cli/actor/v3action"
	"code.cloudfoundry.org/cli/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Application Summary Actions", func() {
	var (
		actor       *Actor
		fakeV2Actor *v2v3actionfakes.FakeV2Actor
		fakeV3Actor *v2v3actionfakes.FakeV3Actor
	)

	BeforeEach(func() {
		fakeV2Actor = new(v2v3actionfakes.FakeV2Actor)
		fakeV3Actor = new(v2v3actionfakes.FakeV3Actor)
		actor = NewActor(fakeV2Actor, fakeV3Actor)
	})

	Describe("ApplicationSummary", func() {
		DescribeTable("CalculatedBuildpacks",
			func(appSummary ApplicationSummary, expectedBuildpacks []string) {
				Expect(appSummary.CalculatedBuildpacks()).To(Equal(expectedBuildpacks))
			},

			Entry("returns buildpacks when they are set in v3",
				ApplicationSummary{
					Buildpacks: []string{"ruby-bp", "java-bp"},
				},
				[]string{"ruby-bp", "java-bp"},
			),

			Entry("returns buildpack when buildpacks is not set and buildpack is set",
				ApplicationSummary{
					ApplicationSummary: v2action.ApplicationSummary{
						Application: v2action.Application{
							Buildpack: types.FilteredString{
								Value: "some-bp",
								IsSet: true,
							},
						},
					},
				},
				[]string{"some-bp"},
			),

			Entry("returns buildpacks when both buildpack and buildpacks are set",
				ApplicationSummary{
					ApplicationSummary: v2action.ApplicationSummary{
						Application: v2action.Application{
							Buildpack: types.FilteredString{
								Value: "some-bp",
								IsSet: true,
							},
						},
					},
					Buildpacks: []string{"ruby-bp", "java-bp"},
				},
				[]string{"ruby-bp", "java-bp"},
			),

			Entry("returns detected buildpack when neither buildpack nor buildpacks are set",
				ApplicationSummary{
					ApplicationSummary: v2action.ApplicationSummary{
						Application: v2action.Application{
							DetectedBuildpack: types.FilteredString{
								Value: "some-bp",
								IsSet: true,
							},
						},
					},
				},
				[]string{"some-bp"},
			),
		)
	})

	Describe("GetApplicationSummaryByNameAndSpace", func() {
		var (
			appName    string
			spaceGUID  string
			appSummary ApplicationSummary
			warnings   Warnings
			executeErr error
		)

		BeforeEach(func() {
			appName = "dora"
			spaceGUID = "some-guid"
		})

		JustBeforeEach(func() {
			appSummary, warnings, executeErr = actor.GetApplicationSummaryByNameAndSpace(appName, spaceGUID)
		})

		Context("when the app exists", func() {
			var v2AppSummary v2action.ApplicationSummary

			BeforeEach(func() {
				v2AppSummary = v2action.ApplicationSummary{
					Application: v2action.Application{
						Name:      appName,
						SpaceGUID: spaceGUID,
						Buildpack: types.FilteredString{Value: "banana", IsSet: true},
					},
				}
				fakeV2Actor.GetApplicationSummaryByNameAndSpaceReturns(
					v2AppSummary,
					v2action.Warnings{"v2-warning"},
					nil)
			})

			Context("when getting additional v3 details", func() {
				var buildpacks []string

				BeforeEach(func() {
					buildpacks = []string{"banana", "potato"}
					fakeV3Actor.GetApplicationByNameAndSpaceReturns(
						v3action.Application{
							Name:                appName,
							LifecycleBuildpacks: buildpacks,
						},
						v3action.Warnings{"v3-warning"},
						nil)
				})

				It("returns the app summary with warnings", func() {
					Expect(executeErr).ToNot(HaveOccurred())
					Expect(warnings).To(ConsistOf("v2-warning", "v3-warning"))
					Expect(appSummary).To(Equal(ApplicationSummary{
						ApplicationSummary: v2AppSummary,
						Buildpacks:         buildpacks,
					}))
				})
			})

			Context("when the v3 api errors", func() {
				var expectedErr error

				BeforeEach(func() {
					expectedErr = errors.New("v3 api error")
					fakeV3Actor.GetApplicationByNameAndSpaceReturns(
						v3action.Application{},
						v3action.Warnings{"v3-warning"},
						expectedErr)
				})

				It("returns errors and warnings", func() {
					Expect(executeErr).To(MatchError(expectedErr))
					Expect(warnings).To(ConsistOf("v2-warning", "v3-warning"))
				})
			})
		})

		Context("when the v2 API errors", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("api error")
				fakeV2Actor.GetApplicationSummaryByNameAndSpaceReturns(
					v2action.ApplicationSummary{},
					v2action.Warnings{"v2-warning"},
					expectedErr)
			})

			It("returns errors and warnings", func() {
				Expect(executeErr).To(MatchError(expectedErr))
				Expect(warnings).To(ConsistOf("v2-warning"))
			})
		})
	})
})
