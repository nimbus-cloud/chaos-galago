package model

type VCAPApplication struct {
	ApplicationName string   `json:"application_name"`
	ApplicationURIs []string `json:"application_uris"`
	SpaceID         string   `json:"space_id"`
}
