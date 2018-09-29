package migration

import (
	"fmt"
)

const (
	TypeSQL  = iota
	TypeFizz = iota
)

var TypeMappings = map[string]uint8{
	".sql":  TypeSQL,
	".fizz": TypeFizz,
}

type Migration struct {
	Version   string
	Namespace string
	Depends   []string
}

func (m Migration) Revision() string {
	return fmt.Sprintf("%s/%s", m.Namespace, m.Version)
}

func (m Migration) GetNamespace() string {
	return m.Namespace
}

func (m *Migration) Dependencies() []string {
	return m.Depends
}

func (m *Migration) SetDependencies() error {
	return nil
}

type MigrationI interface {
	UpScript() (string, uint8, error)
	DownScript() (string, uint8, error)
	SetDependencies() error
	Revision() string
	GetNamespace() string
	Dependencies() []string
}
