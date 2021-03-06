package organization_test

import (
	"github.com/cloudfoundry/cli/cf/api/organizations/organizationsfakes"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/cf/trace/fakes"
	"github.com/cloudfoundry/cli/flags"
	"github.com/cloudfoundry/cli/plugin/models"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/cli/cf/commands/organization"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
)

var _ = Describe("orgs command", func() {
	var (
		ui                  *testterm.FakeUI
		orgRepo             *organizationsfakes.FakeOrganizationRepository
		configRepo          core_config.Repository
		requirementsFactory *testreq.FakeReqFactory
		deps                command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = configRepo
		deps.RepoLocator = deps.RepoLocator.SetOrganizationRepository(orgRepo)
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("orgs").SetDependency(deps, pluginCall))
	}

	runCommand := func(args ...string) bool {
		return testcmd.RunCliCommand("orgs", args, requirementsFactory, updateCommandDependency, false)
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		configRepo = testconfig.NewRepositoryWithDefaults()
		orgRepo = new(organizationsfakes.FakeOrganizationRepository)
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}

		deps = command_registry.NewDependency(new(fakes.FakePrinter))
	})

	Describe("requirements", func() {
		It("fails when not logged in", func() {
			requirementsFactory.LoginSuccess = false

			Expect(runCommand()).To(BeFalse())
		})

		Context("when arguments are provided", func() {
			var cmd command_registry.Command
			var flagContext flags.FlagContext

			BeforeEach(func() {
				cmd = &organization.ListOrgs{}
				cmd.SetDependency(deps, false)
				flagContext = flags.NewFlagContext(cmd.MetaData().Flags)
			})

			It("should fail with usage", func() {
				flagContext.Parse("blahblah")

				reqs := cmd.Requirements(requirementsFactory, flagContext)

				err := testcmd.RunRequirements(reqs)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage"))
				Expect(err.Error()).To(ContainSubstring("No argument required"))
			})
		})
	})

	Describe("when invoked by a plugin", func() {
		var (
			pluginOrgsModel []plugin_models.GetOrgs_Model
		)

		BeforeEach(func() {
			org1 := models.Organization{}
			org1.Name = "Organization-1"
			org1.Guid = "org-1-guid"

			org2 := models.Organization{}
			org2.Name = "Organization-2"

			org3 := models.Organization{}
			org3.Name = "Organization-3"

			orgRepo.ListOrgsReturns([]models.Organization{org1, org2, org3}, nil)

			pluginOrgsModel = []plugin_models.GetOrgs_Model{}
			deps.PluginModels.Organizations = &pluginOrgsModel
		})

		It("populates the plugin models upon execution", func() {
			testcmd.RunCliCommand("orgs", []string{}, requirementsFactory, updateCommandDependency, true)
			Expect(pluginOrgsModel[0].Name).To(Equal("Organization-1"))
			Expect(pluginOrgsModel[0].Guid).To(Equal("org-1-guid"))
			Expect(pluginOrgsModel[1].Name).To(Equal("Organization-2"))
			Expect(pluginOrgsModel[2].Name).To(Equal("Organization-3"))
		})
	})

	Context("when there are orgs to be listed", func() {
		BeforeEach(func() {
			org1 := models.Organization{}
			org1.Name = "Organization-1"

			org2 := models.Organization{}
			org2.Name = "Organization-2"

			org3 := models.Organization{}
			org3.Name = "Organization-3"

			orgRepo.ListOrgsReturns([]models.Organization{org1, org2, org3}, nil)
		})

		It("tries to get the organizations", func() {
			runCommand()
			Expect(orgRepo.ListOrgsCallCount()).To(Equal(1))
			Expect(orgRepo.ListOrgsArgsForCall(0)).To(Equal(0))
		})

		It("lists orgs", func() {
			runCommand()

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting orgs as my-user"},
				[]string{"Organization-1"},
				[]string{"Organization-2"},
				[]string{"Organization-3"},
			))
		})
	})

	It("tells the user when no orgs were found", func() {
		orgRepo.ListOrgsReturns([]models.Organization{}, nil)
		runCommand()

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Getting orgs as my-user"},
			[]string{"No orgs found"},
		))
	})
})
