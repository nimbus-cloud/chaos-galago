package utils

import (
	"chaos-galago/broker/Godeps/_workspace/src/chaos-galago/shared/model"
	_ "chaos-galago/broker/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"chaos-galago/broker/Godeps/_workspace/src/github.com/gorilla/mux"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func RemoveGreenFromURI(URI string) string {
	return strings.Replace(URI, "-green", "", 1)
}

func GetVCAPApplicationVars(object interface{}) error {
	vcapApplication := os.Getenv("VCAP_APPLICATION")
	err := json.Unmarshal([]byte(vcapApplication), object)
	if err != nil {
		return err
	}
	return nil
}

func ReadAndUnmarshal(object interface{}, dir string, fileName string) error {
	path := dir + string(os.PathSeparator) + fileName

	bytes, err := ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, object)
	if err != nil {
		return err
	}

	return nil
}

func SetupInstanceDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS service_instances
	(
		id varchar(255),
		dashboardURL varchar(255),
		planID varchar(255),
		probability decimal(2,2),
		frequency int
	)`)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func SetupBindingDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS service_bindings
	(
		id varchar(255),
		appID varchar(255),
		servicePlanID varchar(255),
		serviceInstanceID varchar(255),
		lastProcessed varchar(255)
	)`)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func UpdateServiceInstance(db *sql.DB, serviceInstanceID string, probability float64, frequency int) error {
	_, err := db.Exec("UPDATE service_instances SET probability=?,frequency=? WHERE id=?", probability, frequency, serviceInstanceID)
	if err != nil {
		return err
	}
	return nil
}

func GetServiceInstance(db *sql.DB, serviceInstanceID string) (shared_model.ServiceInstance, error) {
	var (
		id, dashboardURL, planID string
		probability              float64
		frequency                int
	)

	row := db.QueryRow("SELECT * FROM service_instances WHERE id=?", serviceInstanceID)
	err := row.Scan(&id, &dashboardURL, &planID, &probability, &frequency)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return shared_model.ServiceInstance{}, nil
		}
		return shared_model.ServiceInstance{}, err
	}
	if id == "" {
		return shared_model.ServiceInstance{}, err
	}
	return shared_model.ServiceInstance{ID: id, DashboardURL: dashboardURL, PlanID: planID, Probability: probability, Frequency: frequency}, nil
}

func DeleteServiceInstanceBindings(db *sql.DB, serviceInstanceID string) error {
	_, err := db.Exec("DELETE FROM service_bindings WHERE serviceInstanceID=?", serviceInstanceID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteServiceBinding(db *sql.DB, serviceBindingID string) error {
	_, err := db.Exec("DELETE FROM service_bindings WHERE id=?", serviceBindingID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteServiceInstance(db *sql.DB, serviceInstance shared_model.ServiceInstance) error {
	_, err := db.Exec("DELETE FROM service_instances WHERE id=?", serviceInstance.ID)
	if err != nil {
		return err
	}
	return nil
}

func AddServiceBinding(db *sql.DB, serviceBinding shared_model.ServiceBinding) error {
	_, err := db.Exec("INSERT INTO service_bindings VALUES (?, ?, ?, ?, ?)", serviceBinding.ID, serviceBinding.AppID, serviceBinding.ServicePlanID, serviceBinding.ServiceInstanceID, serviceBinding.LastProcessed)
	if err != nil {
		return err
	}
	return nil
}

func AddServiceInstance(db *sql.DB, serviceInstance shared_model.ServiceInstance) error {
	_, err := db.Exec("INSERT INTO service_instances VALUES (?, ?, ?, ?, ?)", serviceInstance.ID, serviceInstance.DashboardURL, serviceInstance.PlanID, serviceInstance.Probability, serviceInstance.Frequency)
	if err != nil {
		return err
	}
	return nil
}

func MarshalAndRecord(object interface{}, dir string, fileName string) error {
	MkDir(dir)
	path := dir + string(os.PathSeparator) + fileName

	bytes, err := json.MarshalIndent(object, "", " ")
	if err != nil {
		return err
	}

	return WriteFile(path, bytes)
}

func WriteResponse(w http.ResponseWriter, code int, object interface{}) {
	var (
		data []byte
		err  error
	)

	if str, ok := object.(string); ok {
		data = []byte(str)
	} else {
		data, err = json.Marshal(object)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, string(data))
}

func ProvisionDataFromRequest(r *http.Request, object interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		return err
	}

	return nil
}

func ExtractVarsFromRequest(r *http.Request, varName string) string {
	return mux.Vars(r)[varName]
}

func ReadFile(path string) (content []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	content = bytes

	return
}

func WriteFile(path string, content []byte) error {
	err := ioutil.WriteFile(path, content, 0700)
	if err != nil {
		return err
	}

	return nil
}

func GetPath(paths []string) string {
	workDirectory, _ := os.Getwd()

	if len(paths) == 0 {
		return workDirectory
	}

	resultPath := workDirectory

	for _, path := range paths {
		resultPath += string(os.PathSeparator)
		resultPath += path
	}

	return resultPath
}

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func MkDir(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			return false
		}
	}

	return true
}
