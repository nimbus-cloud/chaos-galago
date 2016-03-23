package sharedModel

// ServiceBinding struct
type ServiceBinding struct {
	ID                string `json:"id"`
	AppID             string `json:"app_guid"`
	ServicePlanID     string `json:"plan_id"`
	ServiceInstanceID string `json:"service_instance_id"`
	LastProcessed     string `json:"LastProcessed"`
}
