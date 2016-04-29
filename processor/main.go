package main

import (
	"database/sql"
	"fmt"
	"github.com/FidelityInternational/chaos-galago/shared/utils"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/FidelityInternational/chaos-galago/processor/utils"
	"os"
	"strconv"
	"time"
)

var (
	dbConnectionString string
	err                error
	config             *cfclient.Config
)

func init() {
	dbConnectionString, err = sharedUtils.GetDBConnectionDetails()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	config = utils.LoadCFConfig()
	fmt.Println("Config loaded:")
	fmt.Println("ApiAddress: ", config.ApiAddress)
	fmt.Println("LoginAddress: ", config.LoginAddress)
	fmt.Println("Username: ", config.Username)
	fmt.Println("SkipSslValidation: ", config.SkipSslValidation)
}

func logError(err error) bool {
	if err != nil {
		fmt.Println("An error has occured")
		fmt.Println(err.Error())
		return true
	}
	return false
}

func main() {
	cfClient := cfclient.NewClient(config)
	ticker := time.NewTicker(1 * time.Minute)

	processServices(cfClient)
	for range ticker.C {
		processServices(cfClient)
	}
}

func processServices(cfClient *cfclient.Client) {
	db, err := sql.Open("mysql", dbConnectionString)
	defer db.Close()
	if logError(err) {
		return
	}

	services := utils.GetBoundApps(db)

	for _, service := range services {
		if utils.ShouldProcess(service.Frequency, service.LastProcessed) {
			fmt.Printf("Processing chaos for %s\n", service.AppID)
			err = utils.UpdateLastProcessed(db, service.AppID, utils.TimeNow())
			if logError(err) {
				continue
			}
			if utils.ShouldRun(service.Probability) {
				fmt.Printf("Running chaos for %s\n", service.AppID)
				appInstances := cfClient.GetAppInstances(service.AppID)
				if utils.IsAppHealthy(appInstances) {
					fmt.Printf("App %s is Healthy\n", service.AppID)
					chaosInstance := strconv.Itoa(utils.PickAppInstance(appInstances))
					fmt.Printf("About to kill app instance: %s at index: %s\n", service.AppID, chaosInstance)
					cfClient.KillAppInstance(service.AppID, chaosInstance)
					err = utils.UpdateLastProcessed(db, service.AppID, utils.TimeNow())
					logError(err)
				} else {
					fmt.Printf("App %s is unhealthy, skipping\n", service.AppID)
				}
			} else {
				fmt.Printf("Not running chaos for %s\n", service.AppID)
				err = utils.UpdateLastProcessed(db, service.AppID, utils.TimeNow())
				logError(err)
			}
		} else {
			fmt.Printf("Skipping processing chaos for %s\n", service.AppID)
		}
	}
}
