package model

// Service struct
type Service struct {
	Probability   float64 `json:"probability"`
	Frequency     int     `json:"frequency"`
	AppID         string  `json:"app_guid"`
	LastProcessed string  `json:"LastProcessed"`
}
