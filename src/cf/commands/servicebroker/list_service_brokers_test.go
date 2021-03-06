package servicebroker_test

import (
	. "cf/commands/servicebroker"
	"cf/configuration"
	"cf/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
)

func callListServiceBrokers(args []string, serviceBrokerRepo *testapi.FakeServiceBrokerRepo) (ui *testterm.FakeUI) {
	ui = &testterm.FakeUI{}
	config := testconfig.NewRepositoryWithDefaults()
	ctxt := testcmd.NewContext("service-brokers", args)
	cmd := NewListServiceBrokers(ui, config, serviceBrokerRepo)
	testcmd.RunCommand(cmd, ctxt, &testreq.FakeReqFactory{})

	return
}

var _ = Describe("service-brokers command", func() {
	var (
		ui                  *testterm.FakeUI
		config              configuration.Repository
		cmd                 ListServiceBrokers
		repo                *testapi.FakeServiceBrokerRepo
		requirementsFactory *testreq.FakeReqFactory
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		config = testconfig.NewRepositoryWithDefaults()
		repo = &testapi.FakeServiceBrokerRepo{}
		cmd = NewListServiceBrokers(ui, config, repo)
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}
	})

	Describe("login requirements", func() {
		It("fails if the user is not logged in", func() {
			requirementsFactory.LoginSuccess = false
			ctxt := testcmd.NewContext("service-brokers", []string{})
			testcmd.RunCommand(cmd, ctxt, requirementsFactory)
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})
	})

	It("lists service brokers", func() {
		repo.ServiceBrokers = []models.ServiceBroker{models.ServiceBroker{
			Name: "service-broker-to-list-a",
			Guid: "service-broker-to-list-guid-a",
			Url:  "http://service-a-url.com",
		}, models.ServiceBroker{
			Name: "service-broker-to-list-b",
			Guid: "service-broker-to-list-guid-b",
			Url:  "http://service-b-url.com",
		}, models.ServiceBroker{
			Name: "service-broker-to-list-c",
			Guid: "service-broker-to-list-guid-c",
			Url:  "http://service-c-url.com",
		}}

		context := testcmd.NewContext("service-brokers", []string{})
		testcmd.RunCommand(cmd, context, requirementsFactory)

		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"Getting service brokers as", "my-user"},
			{"name", "url"},
			{"service-broker-to-list-a", "http://service-a-url.com"},
			{"service-broker-to-list-b", "http://service-b-url.com"},
			{"service-broker-to-list-c", "http://service-c-url.com"},
		})
	})

	It("says when no service brokers were found", func() {
		context := testcmd.NewContext("service-brokers", []string{})
		testcmd.RunCommand(cmd, context, requirementsFactory)

		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"Getting service brokers as", "my-user"},
			{"No service brokers found"},
		})
	})

	It("reports errors when listing service brokers", func() {
		repo.ListErr = true
		context := testcmd.NewContext("service-brokers", []string{})
		testcmd.RunCommand(cmd, context, requirementsFactory)

		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"Getting service brokers as ", "my-user"},
			{"FAILED"},
		})
	})
})
