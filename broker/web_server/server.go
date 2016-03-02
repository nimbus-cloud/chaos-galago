package web_server

import (
	"fmt"
	"net/http"
	"os"

	"chaos-galago/broker/Godeps/_workspace/src/github.com/gorilla/mux"

	"chaos-galago/broker/Godeps/_workspace/src/chaos-galago/shared/utils"
	_ "chaos-galago/broker/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"chaos-galago/broker/config"
	"chaos-galago/broker/utils"
	"database/sql"
)

var (
	conf               = config.GetConfig()
	db                 *sql.DB
	err                error
	dbConnectionString string
)

type Server struct {
	controller *Controller
}

func CreateServer() (*Server, error) {
	dbConnectionString, err = shared_utils.GetDBConnectionDetails()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db, err = sql.Open("mysql", dbConnectionString)
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

	controller := CreateController(db)

	return &Server{
		controller: controller,
	}, nil
}

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
		fmt.Printf("ListenAndServe:", err)
	}
}
