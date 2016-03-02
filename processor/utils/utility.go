package utils

import (
	"chaos-galago/processor/Godeps/_workspace/src/chaos-galago/shared/model"
	"chaos-galago/processor/Godeps/_workspace/src/chaos-galago/shared/utils"
	"chaos-galago/processor/Godeps/_workspace/src/github.com/cloudfoundry-community/go-cfclient"
	"chaos-galago/processor/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func ShouldRun(probability float64) bool {
	return rand.Float64() <= probability
}

func IsAppHealthy(appInstances map[string]cfclient.AppInstance) bool {
	for _, instance := range appInstances {
		if instance.State != "RUNNING" {
			return false
		}
	}
	return true
}

func PickAppInstance(appInstances map[string]cfclient.AppInstance) int {
	for k := range appInstances {
		value, _ := strconv.Atoi(k)
		return value
	}
	return 0
}

func ShouldProcess(frequency int, lastProcessed string) bool {
	if lastProcessed == "" {
		return true
	}
	layout := "2006-01-02T15:04:05Z"
	timeStamp, err := time.Parse(layout, lastProcessed)
	if err != nil {
		return false
	}
	duration := time.Duration(-frequency) * time.Minute
	if timeStamp.Before(time.Now().Add(duration)) {
		return true
	}
	return false
}

func GetBoundApps(db *sql.DB) []model.Service {
	var services []model.Service
	serviceInstances, err := shared_utils.ReadServiceInstances(db)
	if err != nil {
		return services
	}
	serviceBindings, err := shared_utils.ReadServiceBindings(db)
	if err != nil {
		return services
	}
OUTER:
	for _, binding := range serviceBindings {
		serviceInstance := serviceInstances[binding.ServiceInstanceID]
		appID := binding.AppID
		probability := serviceInstance.Probability
		frequency := serviceInstance.Frequency
		if serviceInstance == (shared_model.ServiceInstance{}) || appID == "" || probability == 0 || frequency == 0 {
			continue OUTER
		}
		services = append(services, model.Service{AppID: appID, LastProcessed: binding.LastProcessed, Probability: probability, Frequency: frequency})
	}
	return services
}

func TimeNow() string {
	layout := "2006-01-02T15:04:05Z"
	return time.Now().Format(layout)
}

func UpdateLastProcessed(db *sql.DB, appID string, lastProcessed string) error {
	_, err := db.Exec("UPDATE service_bindings SET lastProcessed=? WHERE appID=?", lastProcessed, appID)
	if err != nil {
		return err
	}
	return nil
}

func LoadCFConfig() *cfclient.Config {
	var (
		cfServices model.CFServices
	)

	vcapServicesEnv := os.Getenv("VCAP_SERVICES")
	err := json.Unmarshal([]byte(vcapServicesEnv), &cfServices)
	if err != nil {
		return &cfclient.Config{}
	}

	for _, cfService := range cfServices.CFService {
		if cfService.Name == "cf-service" {
			return &cfclient.Config{
				ApiAddress:        fmt.Sprintf("https://api.%s", cfService.Credentials.Domain),
				LoginAddress:      fmt.Sprintf("https://login.%s", cfService.Credentials.Domain),
				Username:          cfService.Credentials.Username,
				Password:          cfService.Credentials.Password,
				SkipSslValidation: cfService.Credentials.SkipSslValidation,
			}
		}
	}

	return &cfclient.Config{}
}
