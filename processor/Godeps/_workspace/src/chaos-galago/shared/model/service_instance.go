package shared_model

type ServiceInstance struct {
	ID           string  `json:"id"`
	DashboardURL string  `json:"dashboard_url"`
	PlanID       string  `json:"plan_id"`
	Probability  float64 `json:"probability"`
	Frequency    int     `json:"frequency"`
}
