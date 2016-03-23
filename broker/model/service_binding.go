package model

// CreateServiceBindingResponse struct
type CreateServiceBindingResponse struct {
	Credentials interface{} `json:"credentials"`
}

// Credential struct
type Credential struct {
	Probability float64 `json:"probability"`
	Frequency   int     `json:"frequency"`
}
