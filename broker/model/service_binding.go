package model

type CreateServiceBindingResponse struct {
	Credentials interface{} `json:"credentials"`
}

type Credential struct {
	Probability float64 `json:"probability"`
	Frequency   int     `json:"frequency"`
}
