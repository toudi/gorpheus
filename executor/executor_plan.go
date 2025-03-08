package executor

import (
	"github.com/toudi/gorpheus/v2/internal/migration"
	"github.com/toudi/gorpheus/v2/internal/planner"
)

// This method is supposed to take the current db state, compare it to
// already applied migrations and come up with the plan to bring it to
// the target version
// if `namespace` is selected, the plan should only select migration path
// that would only apply to this single namespace.
// if `version` is selected, the plan should stop at `version` within `namespace`
// but it also must take dependencies into consideration
// if the version requested is located earlier in the dependency graph than
// the currently applied migrations, the plan should instead consist of rollback
// migrations
func (e *Executor) PreparePlan(namespace string, version string) ([]*migration.Migration, int, error) {
	_planner := planner.New(
		planner.WithMigrations(e.migrations),
		planner.WithApplied(
			e.history.AppliedMigrations(),
		),
	)

	plan, direction, err := _planner.PreparePlan(namespace, version)
	if err != nil {
		return nil, -1, err
	}

	if direction == planner.DIRECTION_BACKWARD {
		plan, err = e.reverseMigrations(plan)
		if err != nil {
			return nil, direction, err
		}
	} else {
		// we just have to take care of the alter field operations since we need to apply deltas for them
		for _, migration := range plan {
			for _, operation := range migration.Operations {
				if operation.AlterField != nil {
					if err = e.dbState.Delta(
						operation.AlterField,
						operation.AlterField.FieldVersion-1,
						operation.AlterField.FieldVersion,
					); err != nil {
						return nil, direction, err
					}
				}
			}
		}
	}

	return plan, direction, nil
}
