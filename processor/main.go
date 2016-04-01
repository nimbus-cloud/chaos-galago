package main

import (
	"database/sql"
	"fmt"
	"github.com/FidelityInternational/chaos-galago/processor/Godeps/_workspace/src/chaos-galago/shared/utils"
	"github.com/FidelityInternational/chaos-galago/processor/Godeps/_workspace/src/github.com/cloudfoundry-community/go-cfclient"
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
	fmt.Println("\nConfig loaded:")
	fmt.Println("ApiAddress: ", config.ApiAddress)
	fmt.Println("LoginAddress: ", config.LoginAddress)
	fmt.Println("Username: ", config.Username)
	fmt.Println("SkipSslValidation: ", config.SkipSslValidation)
}

func freakOut(err error) bool {
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
	for _ = range ticker.C {
		processServices(cfClient)
	}
}

func processServices(cfClient *cfclient.Client) {
	db, err := sql.Open("mysql", dbConnectionString)
	if freakOut(err) {
		db.Close()
		return
	}
	services := utils.GetBoundApps(db)
	if len(services) == 0 {
		db.Close()
		return
	}

	for _, service := range services {
		if utils.ShouldProcess(service.Frequency, service.LastProcessed) {
			fmt.Printf("\nProcessing chaos for %s", service.AppID)
			err = utils.UpdateLastProcessed(db, service.AppID, utils.TimeNow())
			if freakOut(err) {
				continue
			}
			if utils.ShouldRun(service.Probability) {
				fmt.Printf("\nRunning chaos for %s", service.AppID)
				appInstances := cfClient.GetAppInstances(service.AppID)
				if utils.IsAppHealthy(appInstances) {
					fmt.Printf("\nApp %s is Healthy\n", service.AppID)
					chaosInstance := strconv.Itoa(utils.PickAppInstance(appInstances))
					fmt.Printf("\nAbout to kill app instance: %s at index: %s", service.AppID, chaosInstance)
					cfClient.KillAppInstance(service.AppID, chaosInstance)
					err = utils.UpdateLastProcessed(db, service.AppID, utils.TimeNow())
					if freakOut(err) {
						continue
					}
				} else {
					fmt.Printf("\nApp %s is unhealthy, skipping\n", service.AppID)
					continue
				}
			} else {
				fmt.Printf("\nNot running chaos for %s", service.AppID)
				err = utils.UpdateLastProcessed(db, service.AppID, utils.TimeNow())
				if freakOut(err) {
					continue
				}
			}
		} else {
			fmt.Printf("\nSkipping processing chaos for %s", service.AppID)
			continue
		}
	}
	db.Close()
}
