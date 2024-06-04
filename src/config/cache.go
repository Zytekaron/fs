package config

import (
	"errors"
	"fmt"
)

type Cache struct {
	// MaxEntry is the max number of items to cache.
	MaxCount int `yaml:"max_count"`
	// MaxEntry is the max size of an entry, in bytes.
	MaxEntry int64 `yaml:"max_entry"`
}

func (c *Cache) Validate() error {
	if c.MaxCount < 0 {
		return errors.New("cache.max_count cannot be negative")
	}
	if c.MaxEntry < 0 {
		return errors.New("cache.max_entry cannot be negative")
	}
	return nil
}

func (c *Cache) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data struct {
		MaxCount int    `yaml:"max_count"`
		MaxEntry string `yaml:"max_entry"`
	}

	err := unmarshal(&data)
	if err != nil {
		return err
	}

	maxEntry, err := parseFileSize(data.MaxEntry)
	if err != nil {
		return fmt.Errorf("error parsing file size '%s': %w", data.MaxEntry, err)
	}

	c.MaxCount = data.MaxCount
	c.MaxEntry = maxEntry
	return nil
}
