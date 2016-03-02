package model

type CFServices struct {
	CFService []CFService `json:"user-provided"`
}

type CFService struct {
	Credentials CFCredentials `json:"credentials"`
	Name        string        `json:"name"`
}

type CFCredentials struct {
	Domain            string `json:"domain"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	SkipSslValidation bool   `json:"skipsslvalidation"`
}
