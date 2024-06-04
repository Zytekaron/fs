package config

import (
	"errors"
	"fmt"
	"os"
)

type Files struct {
	// Directory is the base directory to serve files from.
	Directory string `yaml:"directory"`
}

func (f *Files) Validate() error {
	file, err := os.Stat(f.Directory)
	if err != nil {
		return fmt.Errorf("error statting files.directory: %w", err)
	}
	if !file.IsDir() {
		return errors.New("files.directory must be a directory")
	}
	return nil
}
