package model

// CreateServiceInstanceResponse struct
type CreateServiceInstanceResponse struct {
	DashboardURL string  `json:"dashboard_url"`
	Probability  float64 `json:"probability"`
	Frequency    int     `json:"frequency"`
}
