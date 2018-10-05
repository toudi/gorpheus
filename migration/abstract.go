package migration

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	TypeSQL  = iota
	TypeFizz = iota
	TypeGo   = iota
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

func (m *Migration) Up(tx *sqlx.Tx) error {
	return nil
}

func (m *Migration) Down(tx *sqlx.Tx) error {
	return nil
}

type MigrationI interface {
	UpScript() (string, uint8, error)
	DownScript() (string, uint8, error)
	Up(tx *sqlx.Tx) error
	Down(tx *sqlx.Tx) error
	SetDependencies() error
	Revision() string
	GetNamespace() string
	Dependencies() []string
}
