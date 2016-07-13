package sharedUtils_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/FidelityInternational/chaos-galago/shared/model"
	"github.com/FidelityInternational/chaos-galago/shared/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("#ReadServiceInstances", func() {
	var (
		serviceInstancesMap map[string]sharedModel.ServiceInstance
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

			serviceInstancesMap, err = sharedUtils.ReadServiceInstances(db)
			Expect(err).To(BeNil())
			Expect(serviceInstancesMap).To(HaveLen(2))
			Expect(serviceInstancesMap["1"]).To(Equal(sharedModel.ServiceInstance{ID: "1", DashboardURL: "example.com/1", PlanID: "1", Probability: 0.2, Frequency: 5}))
			Expect(serviceInstancesMap["2"]).To(Equal(sharedModel.ServiceInstance{ID: "2", DashboardURL: "example.com/2", PlanID: "2", Probability: 0.4, Frequency: 10}))
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

				serviceInstancesMap, err = sharedUtils.ReadServiceInstances(db)
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

				serviceInstancesMap, err = sharedUtils.ReadServiceInstances(db)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("sql: expected 3 destination arguments in Scan, not 5"))
			})
		})

		Context("when the database return an error", func() {
			It("Returns an error", func() {
				db, mock, err := sqlmock.New()
				if err != nil {
					fmt.Printf("\nan error '%s' was not expected when opening a stub database connection\n", err)
					os.Exit(1)
				}
				defer db.Close()

				mock.ExpectQuery("^SELECT (.+) FROM service_instances$").WillReturnError(fmt.Errorf("An error was raised: %s", "Database Error"))

				serviceInstancesMap, err = sharedUtils.ReadServiceInstances(db)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("An error was raised: Database Error"))
			})
		})

		Context("when the rows return an error", func() {
			It("Returns an error", func() {
				db, mock, err := sqlmock.New()
				if err != nil {
					fmt.Printf("\nan error '%s' was not expected when opening a stub database connection\n", err)
					os.Exit(1)
				}
				defer db.Close()

				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"}).
					AddRow("1", "example.com/1", "1", 0.2, 5).
					AddRow("2", "example.com/2", "2", 0.4, 10).
					RowError(1, fmt.Errorf("An error was raised: %s", "Row Error"))

				mock.ExpectQuery("^SELECT (.+) FROM service_instances$").WillReturnRows(rows)

				serviceInstancesMap, err = sharedUtils.ReadServiceInstances(db)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("An error was raised: Row Error"))
			})
		})
	})
})

var _ = Describe("#ReadServiceBindings", func() {
	var (
		serviceBindingsMap map[string]sharedModel.ServiceBinding
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

			serviceBindingsMap, err = sharedUtils.ReadServiceBindings(db)
			Expect(err).To(BeNil())
			Expect(serviceBindingsMap).To(HaveLen(2))
			Expect(serviceBindingsMap["1"]).To(Equal(sharedModel.ServiceBinding{ID: "1", AppID: "1", ServicePlanID: "1", ServiceInstanceID: "1", LastProcessed: "2014-11-12T10:31:20Z"}))
			Expect(serviceBindingsMap["2"]).To(Equal(sharedModel.ServiceBinding{ID: "2", AppID: "2", ServicePlanID: "2", ServiceInstanceID: "2", LastProcessed: "2014-11-12T10:34:20Z"}))
		})
	})

	Context("When the database schema is incorrect", func() {
		Context("because the database is missing a field", func() {
			It("Returns an error", func() {
				db, mock, err := sqlmock.New()
				if err != nil {
					fmt.Printf("\nan error '%s' was not expected when opening a stub database connection\n", err)
					os.Exit(1)
				}
				defer db.Close()

				rows := sqlmock.NewRows([]string{"id", "appID", "servicePlanID", "serviceInstanceID"}).
					AddRow("1", "1", "1", "1").
					AddRow("2", "2", "2", "2")

				mock.ExpectQuery("^SELECT (.+) FROM service_bindings$").WillReturnRows(rows)

				serviceBindingsMap, err = sharedUtils.ReadServiceBindings(db)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("sql: expected 4 destination arguments in Scan, not 5"))
			})
		})

		Context("when the database return an error", func() {
			It("Returns an error", func() {
				db, mock, err := sqlmock.New()
				if err != nil {
					fmt.Printf("\nan error '%s' was not expected when opening a stub database connection\n", err)
					os.Exit(1)
				}
				defer db.Close()

				mock.ExpectQuery("^SELECT (.+) FROM service_bindings$").WillReturnError(fmt.Errorf("An error was raised: %s", "Database Error"))

				serviceBindingsMap, err = sharedUtils.ReadServiceBindings(db)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("An error was raised: Database Error"))
			})
		})

		Context("when the rows return an error", func() {
			It("Returns an error", func() {
				db, mock, err := sqlmock.New()
				if err != nil {
					fmt.Printf("\nan error '%s' was not expected when opening a stub database connection\n", err)
					os.Exit(1)
				}
				defer db.Close()

				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"}).
					AddRow("1", "1", "1", "1", "2014-11-12T10:31:20Z").
					AddRow("2", "2", "2", "2", "2014-11-12T10:34:20Z").
					RowError(1, fmt.Errorf("An error was raised: %s", "Row Error"))

				mock.ExpectQuery("^SELECT (.+) FROM service_bindings$").WillReturnRows(rows)

				serviceBindingsMap, err = sharedUtils.ReadServiceBindings(db)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("An error was raised: Row Error"))
			})
		})
	})
})

var _ = Describe("GetDBConnectionDetails", func() {
	Context("when a chaos-galago-db service does not exist", func() {
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
  "name": "other-service",
  "syslog_drain_url": "",
  "tags": []
 }
]
}`
		})

		JustBeforeEach(func() {
			os.Setenv("VCAP_SERVICES", vcapServicesJSON)
			os.Setenv("VCAP_APPLICATION", "{}")
		})

		AfterEach(func() {
			os.Unsetenv("VCAP_SERVICES")
		})
		It("returns an error", func() {
			_, err := sharedUtils.GetDBConnectionDetails()
			Expect(err).To(MatchError("no service with name chaos-galago-db"))
		})
	})

	Context("when a chaos-galago-db service exists", func() {
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
		})

		JustBeforeEach(func() {
			os.Setenv("VCAP_SERVICES", vcapServicesJSON)
			os.Setenv("VCAP_APPLICATION", "{}")
		})

		AfterEach(func() {
			os.Unsetenv("VCAP_SERVICES")
		})

		It("Returns the database connection string", func() {
			dbConnString, err := sharedUtils.GetDBConnectionDetails()
			Expect(err).To(BeNil())
			Expect(dbConnString).To(Equal("test_user:test_password@tcp(test_host:test_port)/test_database"))
		})

		Context("When unmarshaling a managed database connection", func() {
			BeforeEach(func() {
				vcapServicesJSON = `
		 {
		   "p-mysql": [
		    {
		     "credentials": {
		      "hostname": "test_host",
		      "jdbcUrl": "jdbc:mysql:/test_host:3306/test_database?user=test_user\u0026password=test_password",
		      "name": "test_database",
		      "password": "test_password",
		      "port": 3306,
		      "uri": "mysql://test_user:test_password@test_host:3306/test_database?reconnect=true",
		      "username": "test_user"
		     },
		     "label": "p-mysql",
		     "name": "chaos-galago-db",
		     "plan": "512mb",
		     "provider": null,
		     "syslog_drain_url": null,
		     "tags": [
		      "mysql"
		     ]
		    }
		   ]
		 }`
				os.Setenv("VCAP_SERVICES", vcapServicesJSON)

			})

			It("returns the database connection string", func() {
				dbConnString, err := sharedUtils.GetDBConnectionDetails()
				Expect(err).To(BeNil())
				Expect(dbConnString).To(Equal("test_user:test_password@tcp(test_host:3306)/test_database"))
			})
		})

		Context("When unmarshalling raises an error", func() {
			BeforeEach(func() {
				vcapServicesJSON = `{
  "user-provided": [
   {
    "credenti
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
			It("Returns an error", func() {
				_, err := sharedUtils.GetDBConnectionDetails()
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(MatchRegexp("invalid character"))
			})
		})
	})
})
