package web_server_test

import (
	"bytes"
	"chaos-galago/broker/Godeps/_workspace/src/chaos-galago/shared/model"
	"chaos-galago/broker/Godeps/_workspace/src/github.com/DATA-DOG/go-sqlmock"
	"chaos-galago/broker/Godeps/_workspace/src/github.com/gorilla/mux"
	. "chaos-galago/broker/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "chaos-galago/broker/Godeps/_workspace/src/github.com/onsi/gomega"
	webs "chaos-galago/broker/web_server"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

func Router(controller *webs.Controller) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/v2/service_instances/{service_instance_guid}", controller.CreateServiceInstance).Methods("PUT")
	r.HandleFunc("/v2/service_instances/{service_instance_guid}", controller.GetServiceInstance).Methods("GET")
	r.HandleFunc("/v2/service_instances/{service_instance_guid}", controller.RemoveServiceInstance).Methods("DELETE")
	r.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", controller.Bind).Methods("PUT")
	r.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", controller.UnBind).Methods("DELETE")
	r.HandleFunc("/dashboard/{service_instance_guid}", controller.GetDashboard).Methods("GET")
	r.HandleFunc("/dashboard/{service_instance_guid}", controller.UpdateServiceInstance).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web_server/resources/")))

	return r
}

func init() {
	var controller *webs.Controller
	http.Handle("/", Router(controller))
}

var _ = Describe("Contoller", func() {
	var (
		db   *sql.DB
		mock sqlmock.Sqlmock
		err  error
	)

	BeforeEach(func() {
		db, mock, err = sqlmock.New()
		if err != nil {
			fmt.Printf("\nan error '%s' was not expected when opening a stub database connection\n", err)
			os.Exit(1)
		}
	})

	AfterEach(func() {
		defer db.Close()
	})

	Describe("#CreateController", func() {

		BeforeEach(func() {
		})

		It("Creates a controller", func() {
			controller := webs.CreateController(db)
			Expect(controller).To(BeAssignableToTypeOf(&webs.Controller{}))
		})
	})

	Describe("#DeleteAssociatedBindings", func() {
		var (
			controller *webs.Controller
		)

		BeforeEach(func() {
			controller = webs.CreateController(db)
		})

		Context("When the binding exists", func() {
			It("returns nil", func() {
				mock.ExpectExec("DELETE FROM service_bindings WHERE serviceInstanceID=").WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))
				Expect(controller.DeleteAssociatedBindings("1")).To(BeNil())
			})
		})

		Context("When the binding does not exist", func() {
			It("returns nil", func() {
				mock.ExpectExec("DELETE FROM service_bindings WHERE serviceInstanceID=").WithArgs("does_not_exist").WillReturnResult(sqlmock.NewResult(1, 0))
				Expect(controller.DeleteAssociatedBindings("does_not_exist")).To(BeNil())
			})
		})
	})

	Describe("#RemoveServiceInstance", func() {
		var (
			controller   *webs.Controller
			req          *http.Request
			mockRecorder *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			controller = webs.CreateController(db)
			mockRecorder = httptest.NewRecorder()
		})

		Context("When the service instance exists", func() {
			BeforeEach(func() {
				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"}).
					AddRow("1", "https://example.com/dashboard/1", "1", 0.2, 5)
				mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("1").WillReturnRows(rows)
				mock.ExpectExec("DELETE FROM service_bindings WHERE serviceInstanceID=").WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("DELETE FROM service_instances WHERE id=").WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))
				req, _ = http.NewRequest("DELETE", "http://example.com/v2/service_instances/1", nil)
				Router(controller).ServeHTTP(mockRecorder, req)
			})

			It("Returns a 200", func() {
				Expect(mockRecorder.Code).To(Equal(200))
				Expect(mockRecorder.Body.String()).To(Equal("{}"))
			})
		})

		Context("When the service instance does not exist", func() {
			BeforeEach(func() {
				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"})
				mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("2").WillReturnRows(rows)
				req, _ = http.NewRequest("DELETE", "http://example.com/v2/service_instances/2", nil)
				Router(controller).ServeHTTP(mockRecorder, req)
			})

			It("Returns a 410", func() {
				Expect(mockRecorder.Code).To(Equal(410))
			})
		})
	})

	Describe("#Unbind", func() {
		var (
			controller   *webs.Controller
			req          *http.Request
			mockRecorder *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			controller = webs.CreateController(db)
			mockRecorder = httptest.NewRecorder()
		})

		Context("When the service instance exists", func() {
			BeforeEach(func() {
				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"}).
					AddRow("1", "https://example.com/dashboard/1", "1", 0.2, 5)
				mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("1").WillReturnRows(rows)
				mock.ExpectExec("DELETE FROM service_bindings WHERE id=").WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))
				req, _ = http.NewRequest("DELETE", "http://example.com/v2/service_instances/1/service_bindings/1", nil)
				Router(controller).ServeHTTP(mockRecorder, req)
			})

			It("Returns a 200", func() {
				Expect(mockRecorder.Code).To(Equal(200))
				Expect(mockRecorder.Body.String()).To(Equal("{}"))
			})
		})

		Context("When the service instance does not exist", func() {
			BeforeEach(func() {
				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"})
				mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("2").WillReturnRows(rows)
				req, _ = http.NewRequest("DELETE", "http://example.com/v2/service_instances/2/service_bindings/2", nil)
				Router(controller).ServeHTTP(mockRecorder, req)
			})

			It("Returns a 410", func() {
				Expect(mockRecorder.Code).To(Equal(410))
			})
		})
	})

	Describe("#GetDashboard", func() {
		var (
			response     string
			controller   *webs.Controller
			req          *http.Request
			mockRecorder *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			controller = webs.CreateController(db)
			mockRecorder = httptest.NewRecorder()
			response = `<html>
	<head>
		<link rel="stylesheet" href="/css/bootstrap.min.css">
		<link rel="stylesheet" href="/css/bootstrap-theme.min.css">
		<script src="/js/jquery-1.11.3.min.js"></script>
		<script src="/js/bootstrap.min.js"></script>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Dashboard</title>
	</head>
	<body>
		<div class="container">
			<h1>Change Service Instance Config</h1>
			<form action="/dashboard/1" method="POST">
				<fieldset class="form-group">
					<label for "probability">Probability</label>
					<input type="number" step="0.01" min="0" max="1" class="form-control" id="probability" name="probability" placeholder="0.2">
				</fieldset>
				<fieldset class="form-group">
					<label for "frequency">Frequency</label>
					<input type="number" min="1" max="60" class="form-control" id="frequency" name="frequency" placeholder="5">
				</fieldset>
				<div class="form-group row">
					<button type="submit" class="btn btn-primary">Submit</button>
				</div>
			</form>
		</div>
	</body>
</html>
`
		})

		Context("When the service instance exists", func() {
			BeforeEach(func() {
				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"}).
					AddRow("1", "https://example.com/dashboard/1", "1", 0.2, 5)
				mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("1").WillReturnRows(rows)
				req, _ = http.NewRequest("GET", "http://example.com/dashboard/1", nil)
				Router(controller).ServeHTTP(mockRecorder, req)
			})
			It("returns the form", func() {
				Expect(mockRecorder.Code).To(Equal(200))
				Expect(mockRecorder.Body.String()).To(Equal(response))
			})
		})

		Context("When the service instance does not exist", func() {
			BeforeEach(func() {
				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"})
				mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("2").WillReturnRows(rows)
				req, _ = http.NewRequest("GET", "http://example.com/v2/service_instances/2/service_bindings/2", nil)
				Router(controller).ServeHTTP(mockRecorder, req)
				It("returns the form", func() {
					Expect(mockRecorder.Code).To(Equal(400))
					Expect(mockRecorder.Body.String()).To(Equal(""))
				})
			})
		})
	})

	Describe("#UpdateServiceInstance", func() {
		var (
			response     string
			controller   *webs.Controller
			req          *http.Request
			mockRecorder *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			controller = webs.CreateController(db)
			mockRecorder = httptest.NewRecorder()
			rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"}).
				AddRow("1", "https://example.com/dashboard/1", "1", 0.2, 5)
			mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("1").WillReturnRows(rows)
		})

		Context("When the service instance exists", func() {
			Context("When probability is invalid", func() {
				BeforeEach(func() {
					response = `<html>
	<head>
		<link rel="stylesheet" href="/css/bootstrap.min.css">
		<link rel="stylesheet" href="/css/bootstrap-theme.min.css">
		<script src="/js/jquery-1.11.3.min.js"></script>
		<script src="/js/bootstrap.min.js"></script>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Dashboard</title>
	</head>
	<body>
		<div class="container">
			<h1>Invalid Configuration Request</h1>
			<p>Probability must be between 0 and 1</p>
			<p>Frequency must be between 1 and 60</p>
		</div>
	</body>
</html>`
					req, _ = http.NewRequest("POST", "http://example.com/dashboard/1", strings.NewReader("probability=3&frequency=5"))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
					Router(controller).ServeHTTP(mockRecorder, req)
				})
				It("returns an error page", func() {
					Expect(mockRecorder.Code).To(Equal(400))
					Expect(mockRecorder.Body.String()).To(Equal(response))
				})
			})

			Context("When frequency is invalid", func() {
				BeforeEach(func() {
					response = `<html>
	<head>
		<link rel="stylesheet" href="/css/bootstrap.min.css">
		<link rel="stylesheet" href="/css/bootstrap-theme.min.css">
		<script src="/js/jquery-1.11.3.min.js"></script>
		<script src="/js/bootstrap.min.js"></script>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Dashboard</title>
	</head>
	<body>
		<div class="container">
			<h1>Invalid Configuration Request</h1>
			<p>Probability must be between 0 and 1</p>
			<p>Frequency must be between 1 and 60</p>
		</div>
	</body>
</html>`
					req, _ = http.NewRequest("POST", "http://example.com/dashboard/1", strings.NewReader("probability=0.1&frequency=0"))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
					Router(controller).ServeHTTP(mockRecorder, req)
				})
				It("returns an error page", func() {
					Expect(mockRecorder.Code).To(Equal(400))
					Expect(mockRecorder.Body.String()).To(Equal(response))
				})
			})

			Context("When probability and frequency are valid", func() {
				BeforeEach(func() {
					mock.ExpectExec("UPDATE service_instances.*").WithArgs(0.4, 10, "1").WillReturnResult(sqlmock.NewResult(1, 1))
					response = `<html>
	<head>
		<link rel="stylesheet" href="/css/bootstrap.min.css">
		<link rel="stylesheet" href="/css/bootstrap-theme.min.css">
		<script src="/js/jquery-1.11.3.min.js"></script>
		<script src="/js/bootstrap.min.js"></script>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Dashboard</title>
	</head>
	<body>
		<div class="container">
			<h1>New Service Instance Configuration</h1>
			<p>Probability: 0.4</p>
			<p>Frequency: 10</p>
		</div>
	</body>
</html>`
					req, _ = http.NewRequest("POST", "http://example.com/dashboard/1", strings.NewReader("probability=0.4&frequency=10"))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
					Router(controller).ServeHTTP(mockRecorder, req)
				})

				It("updates the service isntance", func() {
					Expect(mockRecorder.Code).To(Equal(202))
					Expect(mockRecorder.Body.String()).To(Equal(response))
				})
			})
		})

		Context("When the service instance does not exist", func() {
			BeforeEach(func() {
				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"})
				mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("2").WillReturnRows(rows)
				req, _ = http.NewRequest("POST", "http://example.com/v2/service_instances/2/service_bindings/2", nil)
				Router(controller).ServeHTTP(mockRecorder, req)
				It("returns the form", func() {
					Expect(mockRecorder.Code).To(Equal(400))
					Expect(mockRecorder.Body.String()).To(Equal(""))
				})
			})
		})
	})

	Describe("#GetServiceInstance", func() {
		var (
			controller   *webs.Controller
			req          *http.Request
			mockRecorder *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			controller = webs.CreateController(db)
			mockRecorder = httptest.NewRecorder()
		})

		Context("When the service instance exists", func() {
			BeforeEach(func() {
				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"}).
					AddRow("1", "https://example.com/dashboard/1", "1", 0.2, 5)
				mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("1").WillReturnRows(rows)
				req, _ = http.NewRequest("GET", "http://example.com/v2/service_instances/1", nil)
				Router(controller).ServeHTTP(mockRecorder, req)
			})

			It("returns dashboard URL, probability and frequency", func() {
				Expect(mockRecorder.Code).To(Equal(200))
				Expect(mockRecorder.Body.String()).To(Equal(`{"dashboard_url":"https://example.com/dashboard/1","probability":0.2,"frequency":5}`))
			})
		})

		Context("When the service instance does not exist", func() {
			BeforeEach(func() {
				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"})
				mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("2").WillReturnRows(rows)
				req, _ = http.NewRequest("GET", "http://example.com/v2/service_instances/2", nil)
				Router(controller).ServeHTTP(mockRecorder, req)
			})

			It("returns a 404 not found", func() {
				Expect(mockRecorder.Code).To(Equal(404))
			})
		})
	})

	Describe("#CreateServiceInstance", func() {
		var (
			controller   *webs.Controller
			req          *http.Request
			mockRecorder *httptest.ResponseRecorder
			instance     shared_model.ServiceInstance
		)

		BeforeEach(func() {
			vcapApplicationJSON := `{"application_name": "test", "application_uris": ["example.com"]}`
			os.Setenv("VCAP_APPLICATION", vcapApplicationJSON)
			os.Setenv("PROBABILITY", "0.2")
			os.Setenv("FREQUENCY", "5")
			reqJSON := `{
  "organization_guid": "org-guid-here",
  "plan_id":           "plan-guid-here",
  "service_id":        "service-guid-here",
  "space_guid":        "space-guid-here"
 }`
			req, _ = http.NewRequest("PUT", "http://example.com/v2/service_instances/test", bytes.NewReader([]byte(reqJSON)))
			mockRecorder = httptest.NewRecorder()
			probability := 0.2
			frequency := 5
			planID := "default"
			instanceID := "test"
			applicationURI := "example.com"
			dashboardURL := fmt.Sprintf("https://%s/dashboard/%s", applicationURI, instanceID)

			instance.DashboardURL = dashboardURL
			instance.ID = instanceID
			instance.PlanID = planID
			instance.Probability = probability
			instance.Frequency = frequency
			mock.ExpectExec("INSERT INTO service_instances").WithArgs(instanceID, dashboardURL, planID, probability, frequency).WillReturnResult(sqlmock.NewResult(1, 1))
			controller = webs.CreateController(db)
			Router(controller).ServeHTTP(mockRecorder, req)
		})

		AfterEach(func() {
			os.Unsetenv("VCAP_APPLICATION")
			os.Unsetenv("PROBABILITY")
			os.Unsetenv("FREQUENCY")
		})

		It("Adds a instance and returns dashboard URL, probability and frequency", func() {
			Expect(mockRecorder.Code).To(Equal(201))
			Expect(mockRecorder.Body.String()).To(Equal(`{"dashboard_url":"https://example.com/dashboard/test","probability":0.2,"frequency":5}`))
		})
	})

	Describe("#Bind", func() {
		var (
			controller   *webs.Controller
			req          *http.Request
			mockRecorder *httptest.ResponseRecorder
			binding      shared_model.ServiceBinding
		)

		BeforeEach(func() {
			reqJSON := `{
  "plan_id":      "plan-guid-here",
  "service_id":   "service-guid-here",
  "app_guid":     "app-guid-here"
 }`
			req, _ = http.NewRequest("PUT", "http://example.com/v2/service_instances/test/service_bindings/1", bytes.NewReader([]byte(reqJSON)))
			mockRecorder = httptest.NewRecorder()
		})

		Context("When the service instance exists", func() {
			BeforeEach(func() {
				planID := "1"
				instanceID := "test"
				bindingID := "1"
				appID := "app-guid-here"

				binding.ID = bindingID
				binding.ServicePlanID = planID
				binding.ServiceInstanceID = instanceID
				binding.AppID = appID

				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"}).
					AddRow("test", "https://example.com/dashboard/1", "1", 0.2, 5)
				mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("test").WillReturnRows(rows)
				mock.ExpectExec("INSERT INTO service_bindings").WithArgs(bindingID, appID, planID, instanceID, "").WillReturnResult(sqlmock.NewResult(1, 1))
				controller = webs.CreateController(db)
				Router(controller).ServeHTTP(mockRecorder, req)
			})

			It("Adds a binding and returns credentials", func() {
				Expect(mockRecorder.Code).To(Equal(201))
				Expect(mockRecorder.Body.String()).To(Equal(`{"credentials":{"probability":0.2,"frequency":5}}`))
			})
		})

		Context("When the service instance does not exist", func() {
			BeforeEach(func() {
				rows := sqlmock.NewRows([]string{"id", "dashboardURL", "planID", "probability", "frequency"})
				mock.ExpectQuery("^SELECT (.+) FROM service_instances WHERE id=").WithArgs("test").WillReturnRows(rows)
				controller = webs.CreateController(db)
				Router(controller).ServeHTTP(mockRecorder, req)
			})

			It("Adds a binding and returns credentials", func() {
				Expect(mockRecorder.Code).To(Equal(404))
			})
		})
	})

	Describe("#Catalog", func() {
		var (
			controller   *webs.Controller
			req          *http.Request
			mockRecorder *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			controller = webs.CreateController(db)
			req, _ = http.NewRequest("GET", "http://example.com/v2/catalog", nil)
			mockRecorder = httptest.NewRecorder()
		})

		Context("When catalog path is set", func() {
			AfterEach(func() {
				os.Unsetenv("CATALOG_PATH")
			})

			Context("and the catalog file can be found", func() {
				Context("and the catalog file is valid", func() {
					BeforeEach(func() {
						os.Setenv("CATALOG_PATH", "fixtures/valid")
					})

					It("Returns the catalog", func() {
						controller.Catalog(mockRecorder, req)
						Expect(mockRecorder.Code).To(Equal(200))
						Expect(mockRecorder.Body.String()).To(ContainSubstring(`"id":"chaos-galago"`))
					})
				})

				Context("and the catalog file is invalid", func() {
					BeforeEach(func() {
						os.Setenv("CATALOG_PATH", "fixtures/invalid")
					})

					It("Returns an error 500", func() {
						controller.Catalog(mockRecorder, req)
						Expect(mockRecorder.Code).To(Equal(500))
					})
				})
			})

			Context("and the catalog file cannot be found", func() {
				BeforeEach(func() {
					os.Setenv("CATALOG_PATH", "fixtures/not_found")
				})

				It("Returns an error 500", func() {
					controller.Catalog(mockRecorder, req)
					Expect(mockRecorder.Code).To(Equal(500))
				})
			})
		})

		Context("When catalog path is not set", func() {
			It("Returns an error 500", func() {
				controller.Catalog(mockRecorder, req)
				Expect(mockRecorder.Code).To(Equal(500))
			})
		})
	})
})
