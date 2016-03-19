# chaos-galago

[![Build Status](https://travis-ci.org/FidelityInternational/chaos-galago.svg?branch=master)](https://travis-ci.org/FidelityInternational/chaos-galago)

### Overview

`chaos-galago` is a project designed to assist in proving applications running on Cloud Foundry are cloud resilient. It is deployed as a service broker and therfore available on the marketplace as a bindable service.

It works by aggressively stopping application instances based on the configured probabability and frequency of the bound service instance.

### Useage

```
cf create-service chaos-galago default {service_instance_name}
cf bind-service {app_name} {service_instance_name}
```

A service instance will be created with the default probability (0.2) and frequency (5) set.
Many applications can be bound to a single service-instance if you desire them to use the same probabilty and frequency.

All frequencies are set in minutes and can be any value between 1 and 60.
All probabilities are set as a float and must be between 0 and 1.

Both probability and frequency can be reconfigured via the dashboard. Note: This is an unsecured endpoint.

To get the dashboard url:

```
cf service {service_instance_name}
```

Example output:

```
Service instance: {service_instance_name}
Service: chaos-galago
Plan: default
Description: Provides the ability to cause chaos on bound application instances
Documentation url: https://github.com/FidelityInternational/chaos-galago
Dashboard: https://chaos-galago-broker.example.com/dashboard/{service_instance_guid}

Last Operation
Status: create succeeded
Message:
```

### Deployment

Clone this project to $GOPATH/src.

```
cd $GOPATH/src
git clone https://github.com/FidelityInternational/chaos-galago.git
```

A deploy scipt is included in order to make deploying chaos-galago as simple as possible.

What the script does:
* Creates an Org and Space to deploy chaos-galago too
* Sets up required user-provided-services - these are configuration only!
* Deploys the chaos-galago applications
* Creates the service broker
* Enables service access to the server
* Zero downtime upgrades of the broker if it is already deployed

What the script does **NOT do:
* Create a database - chaos galago requires a SQL database, the deploy script will configure the application to connect to a previously existing database.

Required Variables:

| Variable             | Required | Description                                                                                                                                                                                                                                      |
|----------------------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| CF_DEPLOY_USERNAME   | Yes      | A Cloud Foundry user with rights to create orgs/ spaces and service-brokers, this user will only be used to deploy chaos-galago                                                                                                                  |
| CF_DELOY_PASSWORD    | Yes      | The password for CF_DEPLOY_USERNAME                                                                                                                                                                                                              |
| CF_SYSTEM_DOMAIN     | Yes      | The system domain of Cloud Foundry, the deployment script will need to connect to https://api.CF_SYSTEM_DOMAIN and chaos-galago will need to connect to both https://api.CF_SYSTEM_DOMAIN and https://login.CF_SYSTEM_DOMAIN                     |
| DB_NAME              | Yes      | The database that chaos-galago will use                                                                                                                                                                                                          |
| DB_HOST              | Yes      | The IP address or hostname of the database                                                                                                                                                                                                       |
| DB_PORT              | Yes      | The port used to connect to the database                                                                                                                                                                                                         |
| DB_USERNAME          | Yes      | A database user with rights to create and update tables                                                                                                                                                                                          |
| DB_PASSWORD          | Yes      | The password for DB_USERNAME                                                                                                                                                                                                                     |
| CF_USERNAME          | Yes      | The Cloud Foundry user that chaos-galago will use to connect to the CF API, it must have at least space developer rights over any application that is bound. An admin level user is recommended to allow seamless interaction across all spaces. |
| CF_PASSWORD          | Yes      | The password for CF_USERNAME                                                                                                                                                                                                                     |
| CF_SKIPSSLVALIDATION | Yes      | Must be set to true or false, if you need to use --skip-ssl-validation when you login to the CF CLI then this will need to be set to true to allow chaos-galago to connect to the CF API.                                                        |

Example deploy:

```
CF_SYS_DOMAIN='system_domain.example.com' \
CF_DEPLOY_USERNAME='an_admin_user' \
CF_DEPLOY_PASSWORD='an_admin_user_password' \
DB_NAME='database_name' \
DB_HOST='database_hostname_or_ip' \
DB_PORT='database_port' \
DB_USERNAME='a_database_user' \
DB_PASSWORD='a_database_user_password' \
CF_SKIPSSLVALIDATION=true
CF_USERNAME='a_cf_user' \
CF_PASSWORD='a_cf_user_password' \
deploy.sh
```

### Monitor/Test

To monitor and test that chaos-galago is functioning as expected you can use the `chaos-galago-smoke-tests` project from `https://github.com/FidelityInternational/chaos-galago-smoke-tests`.
