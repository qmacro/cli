package shared_test

import (
	"time"

	"code.cloudfoundry.org/bytefmt"
	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2/constant"
	. "code.cloudfoundry.org/cli/command/v2/shared"
	"code.cloudfoundry.org/cli/types"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = FDescribe("DisplayApplicationSummary", func() {
	var (
		testUI             *ui.UI
		applicationSummary v2action.ApplicationSummary
		displayCommand     bool
	)

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
	})

	JustBeforeEach(func() {
		DisplayApplicationSummary(testUI, applicationSummary, displayCommand)
	})

	Describe("basic information", func() {
		BeforeEach(func() {
			applicationSummary = v2action.ApplicationSummary{
				Application: v2action.Application{
					Name:              "some-app",
					GUID:              "some-app-guid",
					Memory:            types.NullByteSizeInMb{IsSet: true, Value: 128},
					PackageUpdatedAt:  time.Unix(0, 0),
					DetectedBuildpack: types.FilteredString{IsSet: true, Value: "some-buildpack"},
					State:             "STARTED",
				},
				Stack: v2action.Stack{
					Name: "potatos",
				},
				Routes: []v2action.Route{
					{
						Host: "banana",
						Domain: v2action.Domain{
							Name: "fruit.com",
						},
						Path: "/hi",
					},
					{
						Domain: v2action.Domain{
							Name: "foobar.com",
						},
						Port: types.NullInt{IsSet: true, Value: 13},
					},
				},
			}
		})

		It("displays all the common fields", func() {
			Expect(testUI.Out).To(Say("name:\\s+some-app"))
			Expect(testUI.Out).To(Say("requested state:\\s+started"))
			Expect(testUI.Out).To(Say("routes:\\s+banana.fruit.com/hi, foobar.com:13"))
			Expect(testUI.Out).To(Say("last uploaded:\\s+\\w{3} [0-3]\\d \\w{3} [0-2]\\d:[0-5]\\d:[0-5]\\d \\w+ \\d{4}"))
			Expect(testUI.Out).To(Say("stack:\\s+potatos"))
			Expect(testUI.Out).To(Say("buildpack:\\s+some-buildpack"))
		})

		Context("when there are running instances", func() {
			BeforeEach(func() {
				applicationSummary.Instances = types.NullInt{Value: 3, IsSet: true}
				applicationSummary.RunningInstances = []v2action.ApplicationInstanceWithStats{
					{
						ID:          0,
						State:       v2action.ApplicationInstanceState(constant.ApplicationInstanceRunning),
						Since:       1403140717.984577,
						CPU:         0.73,
						Disk:        50 * bytefmt.MEGABYTE,
						DiskQuota:   2048 * bytefmt.MEGABYTE,
						Memory:      100 * bytefmt.MEGABYTE,
						MemoryQuota: 128 * bytefmt.MEGABYTE,
						Details:     "info from the backend",
					},
					{
						ID:          1,
						State:       v2action.ApplicationInstanceState(constant.ApplicationInstanceCrashed),
						Since:       1403100000.900000,
						CPU:         0.37,
						Disk:        50 * bytefmt.MEGABYTE,
						DiskQuota:   2048 * bytefmt.MEGABYTE,
						Memory:      100 * bytefmt.MEGABYTE,
						MemoryQuota: 128 * bytefmt.MEGABYTE,
						Details:     "potato",
					},
				}
			})

			It("displays all instances related information up top", func() {
				Expect(testUI.Out).To(Say("requested state:\\s+started"))
				Expect(testUI.Out).To(Say("instances:\\s+1\\/3"))
				Expect(testUI.Out).To(Say("usage:\\s+128M x 3 instances"))
				Expect(testUI.Out).To(Say("routes:\\s+banana.fruit.com/hi, foobar.com:13"))
			})

			It("display instances table", func() {
				Expect(testUI.Out).To(Say("state\\s+since\\s+cpu\\s+memory\\s+disk\\s+details"))
				Expect(testUI.Out).To(Say(`#0\s+running\s+2014-06-19T01:18:37Z\s+73.0%\s+100M of 128M\s+50M of 2G\s+info from the backend`))
				Expect(testUI.Out).To(Say(`#1\s+crashed\s+2014-06-18T14:00:00Z\s+37.0%\s+100M of 128M\s+50M of 2G\s+potato`))

			})

		})

		Context("when there are no running instances", func() {

		})

		Describe("displayCommand", func() {
			Context("when it is false", func() {

			})

			Context("when it is true", func() {

			})
		})

		Context("when there are isolation segments", func() {

		})
	})
})
