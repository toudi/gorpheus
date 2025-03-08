package postgresql

import (
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (pg *pgDialect) Transaction(handler func(transaction *sql.Tx) error) error {
	var err error
	var transaction *sql.Tx

	if err = pg.connect(); err != nil {
		return err
	}

	transaction, err = pg.connection.Begin()
	if err != nil {
		return err
	}
	if err = handler(transaction); err != nil {
		rollbackErr := transaction.Rollback()
		return errors.Join(err, rollbackErr)
	}
	return transaction.Commit()
}

func (pg *pgDialect) connect() error {
	var err error
	if pg.connection == nil {
		pg.connection, err = sql.Open("pgx", pg.dsn)
		if err != nil {
			return err
		}
		if err = pg.connection.Ping(); err != nil {
			return err
		}
	}
	return nil
}
