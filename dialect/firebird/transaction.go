package firebird

import (
	"database/sql"
	"errors"

	_ "github.com/nakagami/firebirdsql"
	"github.com/xo/dburl"
)

func (fb *fbDialect) Transaction(handler func(transaction *sql.Tx) error) error {
	var err error
	var transaction *sql.Tx

	if err = fb.connect(); err != nil {
		return err
	}

	transaction, err = fb.connection.Begin()
	if err != nil {
		return err
	}
	if err = handler(transaction); err != nil {
		rollbackErr := transaction.Rollback()
		return errors.Join(err, rollbackErr)
	}
	return transaction.Commit()
}

func (fb *fbDialect) connect() error {
	var err error
	if fb.connection == nil {
		fb.connection, err = dburl.Open(fb.dsn)
		if err != nil {
			return err
		}
		if err = fb.connection.Ping(); err != nil {
			return err
		}
	}
	return nil
}
