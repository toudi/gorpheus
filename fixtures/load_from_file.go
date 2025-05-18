package fixtures

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadFromFile(filename string) (*FixtureFile, error) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".json":
		return loadFromJSON(filename)
	case ".yaml":
		return loadFromYAML(filename)
	default:
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}
}

func loadFromJSON(filename string) (*FixtureFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var ff FixtureFile
	err = json.NewDecoder(file).Decode(&ff)
	return &ff, err
}

func loadFromYAML(filename string) (*FixtureFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var ff FixtureFile
	err = yaml.NewDecoder(file).Decode(&ff)
	return &ff, err
}
