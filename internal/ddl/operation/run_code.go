package operation

import "database/sql"

type RunCode struct {
	Apply   func(tx *sql.Tx) error
	Unapply func(tx *sql.Tx) error
}
