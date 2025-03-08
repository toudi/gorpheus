package sqlite3

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xo/dburl"
)

// Transaction is a nice wrapper for dealing with code that should be executed
// within a transaction block.
func (s3 *sqlite3Dialect) Transaction(handler func(transaction *sql.Tx) error) error {
	var err error
	var transaction *sql.Tx

	if err = s3.connect(); err != nil {
		return err
	}

	transaction, err = s3.connection.Begin()
	if err != nil {
		return err
	}
	if err = handler(transaction); err != nil {
		rollbackErr := transaction.Rollback()
		return errors.Join(err, rollbackErr)
	}
	return transaction.Commit()
}

func (s3 *sqlite3Dialect) connect() error {
	var err error
	if s3.connection == nil {
		s3.connection, err = dburl.Open(s3.dsn)
		if err != nil {
			return err
		}
	}
	return nil
}
