package gorpheus

import (
	"github.com/jmoiron/sqlx"
)

type Callable func(tx *sqlx.Tx) error

func DBConnection(p *MigrationParams) (*sqlx.DB, error) {
	return nil, nil
}

func Atomic(db *sqlx.DB, callable Callable) error {
	transaction := db.MustBegin()
	err := callable(transaction)
	if err != nil {
		transaction.Rollback()
		return err
	}
	return transaction.Commit()
}
