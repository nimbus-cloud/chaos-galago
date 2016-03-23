package model

// CFServices struct
type CFServices struct {
	CFService []CFService `json:"user-provided"`
}

// CFService struct
type CFService struct {
	Credentials CFCredentials `json:"credentials"`
	Name        string        `json:"name"`
}

// CFCredentials struct
type CFCredentials struct {
	Domain            string `json:"domain"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	SkipSslValidation bool   `json:"skipsslvalidation"`
}
