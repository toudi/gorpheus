package fixtures

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func (ff *FixtureFile) Save(output string) error {
	ext := filepath.Ext(output)
	switch ext {
	case ".json":
		return ff.saveAsJSON(output)
	case ".yaml":
		return ff.saveAsYAML(output)
	default:
		return fmt.Errorf("unsupported format: %s", ext)
	}
}

func (ff *FixtureFile) saveAsJSON(output string) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()
	var encoder = json.NewEncoder(file)
	return encoder.Encode(ff)
}

func (ff *FixtureFile) saveAsYAML(output string) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()
	var encoder = yaml.NewEncoder(file)
	return encoder.Encode(ff)
}
