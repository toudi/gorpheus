package executor

import (
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/toudi/gorpheus/v2/interfaces"
	"github.com/toudi/gorpheus/v2/internal/planner"

	"github.com/toudi/gorpheus/v2/internal/migration"

	"github.com/samber/lo"
)

func (e *Executor) Migrate(namespace string, version string) error {
	var buffer strings.Builder
	var err error

	plan, direction, err := e.PreparePlan(namespace, version)

	if len(plan) > 0 {
		// if there is a plan, we have to apply all migrations leading up till the
		// starting point in order to reflect the changes on the dbState
		_, startingPoint, found := lo.FindIndexOf(e.migrations, func(m *migration.Migration) bool {
			return m.Revision.Namespace == plan[0].Revision.Namespace && m.Revision.Name == plan[0].Revision.Name
		})
		if !found {
			return errors.New("starting point not found in graph")
		}

		for _, _migration := range e.migrations[:startingPoint+1] {
			// we do not want to fake bump table version unless the migration was already applied
			// this is because when we have dependencies, the migration might be applied afterwards
			// and the starting point is somewhere else rather than 0
			if !e.history.IsMigrationApplied(_migration) {
				continue
			}
			for _, _operation := range _migration.Operations {
				if _operation.AdjustTableVersion != nil {
					e.dbState.AdjustMigratedTableVersion(_operation.AdjustTableVersion)
				}
			}
		}
	}

	if err != nil {
		return err
	}

	var action = "applying"
	if direction == planner.DIRECTION_BACKWARD {
		action = "unapplying"
	}

	for _, migration := range plan {
		e.logger.Log(interfaces.LogMessage{
			Level:   slog.LevelInfo,
			Message: action,
			Payload: []any{
				slog.String("namespace", migration.Revision.Namespace),
				slog.String("migration", migration.Revision.Name),
			},
		})
		// we iterate through each migrations, prepare transaction and execute it
		if err = e.dialect.Transaction(func(transaction *sql.Tx) error {
			for _, operation := range migration.Operations {
				if operation.AdjustTableVersion != nil {
					e.dbState.AdjustMigratedTableVersion(operation.AdjustTableVersion)
					continue
				}

				// if this is a "data" migration then we need to execute the code
				if operation.RunCode != nil {
					if direction == planner.DIRECTION_FORWARD {
						err = operation.RunCode.Apply(transaction)
					} else {
						err = operation.RunCode.Unapply(transaction)
					}

					if err != nil {
						return err
					}
					continue
				}

				// right now this is only for sqlite3 which does not support full range of
				// alter table / add constraint and so on. for these type of operations,
				// a single "source" operation is converted into a sequence of operations
				// any other engine simply returns the same operation as a slice.
				queue, err := e.dialect.Unpack(operation)
				if err != nil {
					return err
				}

				for _, o := range queue {
					buffer.Reset()
					if err := e.dialect.OperationDDL(o, &buffer); err != nil {
						return err
					}
					if buffer.Len() > 0 {
						e.logger.Log(interfaces.LogMessage{
							Level:   slog.LevelDebug,
							Message: "executing",
							Payload: []any{slog.String("buffer", buffer.String())},
						})
						if _, err := transaction.Exec(buffer.String()); err != nil {
							return err
						}
					}
				}

				e.dbState.AdjustVersions(operation)

				buffer.Reset()
			}

			if err = e.history.HandleAppliedMigration(
				transaction,
				migration.Revision.Namespace,
				migration.Revision.Name,
				direction,
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}

		// everything was successful so we can reset the buffer to save some
		// memory.
		buffer.Reset()
	}

	return nil
}
