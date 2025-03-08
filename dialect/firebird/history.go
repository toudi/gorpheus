package firebird

import (
	"database/sql"

	"github.com/toudi/gorpheus/v2/interfaces"
)

const tableExistsQuery = "SELECT COUNT(1) FROM RDB$RELATIONS WHERE LOWER(RDB$RELATION_NAME) = 'migrations_history'"

type fbHistoryManager struct{}

func (fd *fbDialect) HistoryManager() interfaces.MigrationsHistory {
	return &fbHistoryManager{}
}

func (fbhm *fbHistoryManager) HistoryTableExists(transaction *sql.Tx) (bool, error) {
	var err error

	var count int

	res := transaction.QueryRow(tableExistsQuery)
	err = res.Scan(&count)

	return count == 1, err
}

func (fbhm *fbHistoryManager) LogAppliedMigration(transaction *sql.Tx, namespace, migration string) error {
	_, err := transaction.Exec("INSERT INTO migrations_history (namespace, migration, applied) VALUES (?, ?, CURRENT_TIMESTAMP)", namespace, migration)
	return err
}

func (fbhm *fbHistoryManager) UnlogAppliedMigration(transaction *sql.Tx, namespace, migration string) error {
	_, err := transaction.Exec("DELETE FROM migrations_history WHERE namespace = ? AND migration = ?", namespace, migration)
	return err
}
