package interfaces

import (
	"database/sql"
	"time"
)

type MigrationsHistoryRow struct {
	Namespace string    `db:"namespace"`
	Migration string    `db:"migration"`
	AppliedAt time.Time `db:"applied_at"`
}

type MigrationsHistory interface {
	HistoryTableExists(transaction *sql.Tx) (bool, error)
	LogAppliedMigration(transaction *sql.Tx, namespace, migration string) error
	UnlogAppliedMigration(transaction *sql.Tx, namespace, migration string) error
}
