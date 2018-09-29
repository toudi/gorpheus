package gorpheus

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const UP = "UP"
const DOWN = "DOWN"

type Id struct {
	Namespace string
	Version   string
}

func Id_from_string(id string) *Id {
	out := &Id{
		Namespace: "default",
	}
	parts := strings.Split(id, "/")
	if len(parts) == 2 {
		out.Namespace = parts[0]
		out.Version = parts[1]
	} else {
		out.Version = parts[0]
	}

	return out
}

type BaseMigration struct {
	Id Id
}

func (i Id) ToString() string {
	return fmt.Sprintf("%s/%s", i.Namespace, i.Version)
}

type MigrationI interface {
	ID() Id
	Execute(tx *sqlx.Tx) error
	Depends() *Id
}
