package model

// AppInstances struct
type AppInstances struct {
	Instances map[string]AppInstance
}

// AppInstance struct
type AppInstance struct {
	State string `json:"state"`
}
