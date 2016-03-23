package sharedModel

// VCAPServices struct
type VCAPServices struct {
	UserProvided []UserProvidedServices `json:"user-provided"`
}

// UserProvidedServices struct
type UserProvidedServices struct {
	Credentials UPSCredentials `json:"credentials"`
	Name        string         `json:"name"`
}

// UPSCredentials struct
type UPSCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}
