package utils

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func Recast(input map[string]interface{}, dest interface{}) error {
	var buffer bytes.Buffer
	_ = yaml.NewEncoder(&buffer).Encode(input) // always returns nil
	return yaml.NewDecoder(&buffer).Decode(dest)
}
