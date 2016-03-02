package shared_utils_test

import (
	"chaos-galago/shared/model"
	"chaos-galago/shared/utils"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("#ReadServiceInstances", func() {
	var (
		serviceInstancesMap map[string]shared_model.ServiceInstance
	)

	Context("When the database schema is correct", func() {
		It("returns serviceInstancesMap with records", func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				fmt.Printf("\nan error '%s' was not expected when opening a stub database connection\n", err)
				os.Exit(1)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"}).
				AddRow("1", "example.com/1", "1", 0.2, 5).
				AddRow("2", "example.com/2", "2", 0.4, 10)

			mock.ExpectQuery("^SELECT (.+) FROM service_instances$").WillReturnRows(rows)

			serviceInstancesMap, err = shared_utils.ReadServiceInstances(db)
			Expect(err).To(BeNil())
			Expect(serviceInstancesMap).To(HaveLen(2))
			Expect(serviceInstancesMap["1"]).To(Equal(shared_model.ServiceInstance{ID: "1", DashboardURL: "example.com/1", PlanID: "1", Probability: 0.2, Frequency: 5}))
			Expect(serviceInstancesMap["2"]).To(Equal(shared_model.ServiceInstance{ID: "2", DashboardURL: "example.com/2", PlanID: "2", Probability: 0.4, Frequency: 10}))
		})
	})

	Context("When the database schema is incorrect", func() {
		Context("because the database has additional fields", func() {
			It("Returns an error", func() {
				db, mock, err := sqlmock.New()
				if err != nil {
					fmt.Printf("\nan error '%s' was not expected when opening a stub database connection\n", err)
					os.Exit(1)
				}
				defer db.Close()

				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency", "invalid"}).
					AddRow("1", "example.com/1", "1", 0.2, 5, "test").
					AddRow("2", "example.com/2", "2", 0.2, 5, "test")

				mock.ExpectQuery("^SELECT (.+) FROM service_instances$").WillReturnRows(rows)

				serviceInstancesMap, err = shared_utils.ReadServiceInstances(db)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("sql: expected 6 destination arguments in Scan, not 5"))
			})
		})

		Context("because the database is missing a field", func() {
			It("Returns an error", func() {
				db, mock, err := sqlmock.New()
				if err != nil {
					fmt.Printf("\nan error '%s' was not expected when opening a stub database connection\n", err)
					os.Exit(1)
				}
				defer db.Close()

				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID"}).
					AddRow("1", "example.com/1", "1").
					AddRow("2", "example.com/2", "2")

				mock.ExpectQuery("^SELECT (.+) FROM service_instances$").WillReturnRows(rows)

				serviceInstancesMap, err = shared_utils.ReadServiceInstances(db)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("sql: expected 3 destination arguments in Scan, not 5"))
			})
		})
	})
})

var _ = Describe("#ReadServiceBindings", func() {
	var (
		serviceBindingsMap map[string]shared_model.ServiceBinding
	)

	Context("When the database schema is correct", func() {
		It("returns serviceBindingsMap with records", func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				fmt.Printf("\nan error '%s' was not expected when opening a stub database connection\n", err)
				os.Exit(1)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"id", "appID", "servicePlanID", "serviceInstanceID", "lastProcessed"}).
				AddRow("1", "1", "1", "1", "2014-11-12T10:31:20Z").
				AddRow("2", "2", "2", "2", "2014-11-12T10:34:20Z")

			mock.ExpectQuery("^SELECT (.+) FROM service_bindings$").WillReturnRows(rows)

			serviceBindingsMap, err = shared_utils.ReadServiceBindings(db)
			Expect(err).To(BeNil())
			Expect(serviceBindingsMap).To(HaveLen(2))
			Expect(serviceBindingsMap["1"]).To(Equal(shared_model.ServiceBinding{ID: "1", AppID: "1", ServicePlanID: "1", ServiceInstanceID: "1", LastProcessed: "2014-11-12T10:31:20Z"}))
			Expect(serviceBindingsMap["2"]).To(Equal(shared_model.ServiceBinding{ID: "2", AppID: "2", ServicePlanID: "2", ServiceInstanceID: "2", LastProcessed: "2014-11-12T10:34:20Z"}))
		})
	})
})

var _ = Describe("GetDBConnectionDetails", func() {
	var vcapServicesJSON string

	BeforeEach(func() {
		vcapServicesJSON = `{
  "user-provided": [
   {
    "credentials": {
    	"username":"test_user",
    	"password":"test_password",
    	"host":"test_host",
    	"port":"test_port",
    	"database":"test_database"
    },
    "label": "user-provided",
    "name": "chaos-galago-db",
    "syslog_drain_url": "",
    "tags": []
   }
  ]
 }`
		os.Setenv("VCAP_SERVICES", vcapServicesJSON)
	})

	AfterEach(func() {
		os.Unsetenv("VCAP_SERVICES")
	})

	It("Returns the database connection string", func() {
		dbConnString, err := shared_utils.GetDBConnectionDetails()
		Expect(err).To(BeNil())
		Expect(dbConnString).To(Equal("test_user:test_password@tcp(test_host:test_port)/test_database"))
	})
})
