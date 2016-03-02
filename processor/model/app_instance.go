package model

type AppInstances struct {
	Instances map[string]AppInstance
}

type AppInstance struct {
	State string `json:"state"`
}
