package webServer

import (
	"fmt"
	"net/http"
	"os"

	"chaos-galago/broker/Godeps/_workspace/src/github.com/gorilla/mux"

	"chaos-galago/broker/Godeps/_workspace/src/chaos-galago/shared/utils"
	// mysql driver
	_ "chaos-galago/broker/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"chaos-galago/broker/config"
	"chaos-galago/broker/utils"
	"database/sql"
)

var (
	db                 *sql.DB
	err                error
	dbConnectionString string
)

// Server struct
type Server struct {
	controller *Controller
}

// DBConn - database opening interface
type DBConn func(driverName string, connectionString string) (*sql.DB, error)

// ControllerCreator - controller creation function
type ControllerCreator func(db *sql.DB, conf *config.Config) *Controller

// CreateServer - creates a server
func CreateServer(dbConnFunc DBConn, controllerCreator ControllerCreator) (*Server, error) {
	dbConnectionString, err = sharedUtils.GetDBConnectionDetails()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db, err = dbConnFunc("mysql", dbConnectionString)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = utils.SetupInstanceDB(db)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = utils.SetupBindingDB(db)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	conf := config.GetConfig()
	controller := controllerCreator(db, conf)

	return &Server{
		controller: controller,
	}, nil
}

// Start - starts the web server
func (s *Server) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/v2/catalog", s.controller.Catalog).Methods("GET")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", s.controller.GetServiceInstance).Methods("GET")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", s.controller.CreateServiceInstance).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", s.controller.RemoveServiceInstance).Methods("DELETE")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", s.controller.Bind).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", s.controller.UnBind).Methods("DELETE")
	router.HandleFunc("/dashboard/{service_instance_guid}", s.controller.GetDashboard).Methods("GET")
	router.HandleFunc("/dashboard/{service_instance_guid}", s.controller.UpdateServiceInstance).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web_server/resources/")))

	http.Handle("/", router)

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
