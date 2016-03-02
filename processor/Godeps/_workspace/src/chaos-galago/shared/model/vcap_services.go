package shared_model

type VCAPServices struct {
	UserProvided []UserProvidedServices `json:"user-provided"`
}

type UserProvidedServices struct {
	Credentials UPSCredentials `json:"credentials"`
	Name        string         `json:"name"`
}

type UPSCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}
