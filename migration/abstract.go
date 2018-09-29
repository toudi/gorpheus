package migration

import (
	"database/sql"
	"fmt"
)

const (
	TypeSQL  = iota
	TypeFizz = iota
	TypeGo   = iota
)

const (
	StorageMemory = iota
	StorageFile   = iota
)

var TypeMappings = map[string]uint8{
	".sql":  TypeSQL,
	".fizz": TypeFizz,
}

type Migration struct {
	Version   string
	Namespace string
	Depends   []string
	Type      uint8
	Storage   uint8
}

func (m Migration) Revision() string {
	return fmt.Sprintf("%s/%s", m.Namespace, m.Version)
}

func (m *Migration) Dependencies() []string {
	return m.Depends
}

type MigrationI interface {
	UpFn(tx *sql.Tx) error
	DownFn(tx *sql.Tx) error
	Parse() error
	Revision() string
	Dependencies() []string
}
