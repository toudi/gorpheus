package postgresql

import (
	"database/sql"

	"github.com/toudi/gorpheus/v2/interfaces"
)

const tableExistsQuery = "SELECT COUNT(1) FROM information_schema.tables WHERE table_name = 'migrations_history'"

type pgHistoryManager struct{}

func (pg *pgDialect) HistoryManager() interfaces.MigrationsHistory {
	return &pgHistoryManager{}
}

func (pghm *pgHistoryManager) HistoryTableExists(transaction *sql.Tx) (bool, error) {
	var err error

	var count int

	res := transaction.QueryRow(tableExistsQuery)
	err = res.Scan(&count)

	return count == 1, err
}

func (pghm *pgHistoryManager) LogAppliedMigration(transaction *sql.Tx, namespace, migration string) error {
	_, err := transaction.Exec("INSERT INTO migrations_history (namespace, migration, applied) VALUES ($1, $2, now())", namespace, migration)
	return err
}

func (pghm *pgHistoryManager) UnlogAppliedMigration(transaction *sql.Tx, namespace, migration string) error {
	_, err := transaction.Exec("DELETE FROM migrations_history WHERE namespace = $1 AND migration = $2", namespace, migration)
	return err
}
