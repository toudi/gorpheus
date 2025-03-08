package generic

import (
	"database/sql"
	"log/slog"
	"strings"

	"github.com/toudi/gorpheus/v2/interfaces"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
	"github.com/toudi/gorpheus/v2/internal/migration"
	"github.com/toudi/gorpheus/v2/internal/planner"
	"github.com/toudi/gorpheus/v2/internal/utils"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"
)

var migrationsHistoryOperations = []operation.Operation{
	{
		CreateTable: &operation.CreateTable{
			TableName: "migrations_history",
			Columns: []*column.Column{
				{
					Name:      "namespace",
					Type:      column.TypeString,
					MaxLength: 64,
					Index:     true,
				},
				{
					Name:      "migration",
					Type:      column.TypeString,
					MaxLength: 128,
				},
				{
					Name: "applied",
					Type: column.TypeTimestamp,
				},
			},
		},
	},
}

type MigrationsHistoryManager struct {
	dialect               interfaces.Dialect
	dialectHistoryManager interfaces.MigrationsHistory
	applied               []interfaces.MigrationsHistoryRow
	appliedSet            map[string]bool
	logger                interfaces.Logger
}

func NewMigrationsHistoryManager(dialectSpecific interfaces.Dialect, logger interfaces.Logger) (*MigrationsHistoryManager, error) {
	mhm := &MigrationsHistoryManager{
		dialect:               dialectSpecific,
		dialectHistoryManager: dialectSpecific.HistoryManager(),
		appliedSet:            make(map[string]bool),
		logger:                logger,
	}
	err := dialectSpecific.Transaction(func(transaction *sql.Tx) error {
		historyTableExists, err := mhm.dialectHistoryManager.HistoryTableExists(transaction)
		if err != nil {
			return err
		}
		if historyTableExists {
			rows, err := transaction.Query("SELECT namespace, migration, applied FROM migrations_history ORDER BY applied")
			if err != nil {
				return err
			} else {
				for rows.Next() {
					var row interfaces.MigrationsHistoryRow
					if err := rows.Scan(&row.Namespace, &row.Migration, &row.AppliedAt); err != nil {
						return err
					}
					mhm.applied = append(mhm.applied, row)
					mhm.appliedSet[row.Namespace+"/"+row.Migration] = true
				}
			}
		} else {
			var buffer strings.Builder
			if err = OperationsDDL(migrationsHistoryOperations, dialectSpecific, &buffer); err != nil {
				return err
			} else {
				if _, err = transaction.Exec(buffer.String()); err != nil {
					return err
				}
			}
		}

		return err
	})

	return mhm, err
}

func (mhm *MigrationsHistoryManager) HandleAppliedMigration(transaction *sql.Tx, namespace, migration string, direction int) error {
	if direction == planner.DIRECTION_FORWARD {
		mhm.logger.Log(
			interfaces.LogMessage{
				Level:   slog.LevelDebug,
				Message: "log applied migration",
				Payload: []any{
					slog.String("id", utils.MigrationID(namespace, migration)),
				},
			},
		)
		return mhm.dialectHistoryManager.LogAppliedMigration(transaction, namespace, migration)
	}
	mhm.logger.Log(
		interfaces.LogMessage{
			Level:   slog.LevelDebug,
			Message: "unlog applied migration",
			Payload: []any{
				slog.String("id", utils.MigrationID(namespace, migration)),
			},
		},
	)
	return mhm.dialectHistoryManager.UnlogAppliedMigration(
		transaction, namespace, migration,
	)
}

func (mhm *MigrationsHistoryManager) AppliedMigrations() []interfaces.MigrationsHistoryRow {
	return mhm.applied
}

func (mhm *MigrationsHistoryManager) IsMigrationApplied(_migration *migration.Migration) bool {
	return mhm.appliedSet[_migration.Revision.String()]
}

func (mhm *MigrationsHistoryManager) IsApplied(namespace string, name string) bool {
	return mhm.appliedSet[namespace+"/"+name]
}
