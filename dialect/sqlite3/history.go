package sqlite3

import (
	"database/sql"

	"github.com/toudi/gorpheus/v2/interfaces"
)

const tableExistsQuery = "SELECT COUNT(1) FROM sqlite_master WHERE type = 'table' AND name = 'migrations_history'"
const logMigratinQuery = "INSERT INTO migrations_history (namespace, migration, applied) VALUES (:namespace, :migration, unixepoch())"
const unlogMigratinQuery = "DELETE FROM migrations_history WHERE namespace = :namespace AND migration = :migration"

type s3HistoryManager struct{}

func (s3 *sqlite3Dialect) HistoryManager() interfaces.MigrationsHistory {
	return &s3HistoryManager{}
}

func (s3hm *s3HistoryManager) HistoryTableExists(transaction *sql.Tx) (bool, error) {
	var err error

	var count int

	res := transaction.QueryRow(tableExistsQuery)
	err = res.Scan(&count)

	return count == 1, err
}

func (s3hm *s3HistoryManager) LogAppliedMigration(transaction *sql.Tx, namespace, migration string) error {
	_, err := transaction.Exec(
		logMigratinQuery,
		sql.Named("namespace", namespace),
		sql.Named("migration", migration),
	)
	return err
}

func (s3hm *s3HistoryManager) UnlogAppliedMigration(transaction *sql.Tx, namespace, migration string) error {
	_, err := transaction.Exec(
		unlogMigratinQuery,
		sql.Named("namespace", namespace),
		sql.Named("migration", migration),
	)
	return err
}
