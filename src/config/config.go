package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server *Server `yaml:"server"`
	Files  *Files  `yaml:"files"`
	Cache  *Cache  `yaml:"cache"`
}

func (c *Config) Validate() error {
	err := c.Server.Validate()
	if err != nil {
		return err
	}
	err = c.Files.Validate()
	if err != nil {
		return err
	}
	return c.Cache.Validate()
}

func Load(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	var cfg *Config
	return cfg, yaml.NewDecoder(file).Decode(&cfg)
}
