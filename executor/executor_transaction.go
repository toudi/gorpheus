package executor

import "database/sql"

func (e *Executor) Transaction(handler func(tx *sql.Tx) error) error {
	return e.dialect.Transaction(handler)
}
