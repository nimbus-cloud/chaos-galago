package model

type Service struct {
	Name        string        `json:"name"`
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Bindable    bool          `json:"bindable"`
	Plans       []ServicePlan `json:"plans"`
	Metadata    interface{}   `json:"metadata, omitempty"`
}
