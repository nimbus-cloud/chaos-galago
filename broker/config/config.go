package config

import (
	"chaos-galago/broker/utils"
	"encoding/json"
)

type Config struct {
	CatalogPath              string  `json:"catalog_path"`
	DefaultProbability       float64 `json:"default_probability"`
	DefaultFrequency         int     `json:"default_frequency"`
	DatabaseConnectionString string
}

var (
	currentConfiguration Config
)

func LoadConfig(path string) (*Config, error) {
	bytes, err := utils.ReadFile(path)
	if err != nil {
		return &currentConfiguration, err
	}

	err = json.Unmarshal(bytes, &currentConfiguration)
	if err != nil {
		return &currentConfiguration, err
	}
	return &currentConfiguration, nil
}

func GetConfig() *Config {
	return &currentConfiguration
}
