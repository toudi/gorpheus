package executor

import (
	"io"

	"github.com/toudi/gorpheus/v2/internal/migration"

	"gopkg.in/yaml.v3"
)

var tableVersion map[string]int

func (e *Executor) Parse(source io.Reader) (*migration.Migration, error) {
	if tableVersion == nil {
		tableVersion = make(map[string]int)
	}
	var m migration.Migration

	decoder := yaml.NewDecoder(source)
	err := decoder.Decode(&m)

	return &m, err
}
