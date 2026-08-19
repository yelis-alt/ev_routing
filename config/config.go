package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	OpenRouteService struct {
		Request struct {
			URL string `yaml:"url"`
		} `yaml:"request"`
	} `yaml:"openrouteservice"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
