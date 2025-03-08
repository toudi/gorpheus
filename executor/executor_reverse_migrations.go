package executor

import (
	"errors"

	"github.com/toudi/gorpheus/v2/internal/migration"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func (e *Executor) reverseMigrations(input []*migration.Migration) ([]*migration.Migration, error) {
	// the length of output will be known so let's use this to our advantage.
	reversedMigrations := make([]*migration.Migration, 0, len(input))

	// we don't need to reverse the slice in place as the planner already prepared
	// appropriate plan for us - we simply have to reverse the operations themselves.
	for _, migration := range input {
		reversed, err := e.reverseMigration(migration)
		if err != nil {
			return nil, err
		}
		reversedMigrations = append(reversedMigrations, reversed)
	}

	return reversedMigrations, nil
}

func (e *Executor) reverseMigration(input *migration.Migration) (*migration.Migration, error) {
	reversedOperations, err := e.reverseOperations(input.Operations)
	if err != nil {
		return nil, err
	}
	return &migration.Migration{
		Revision: migration.Revision{
			Namespace: input.Revision.Namespace,
			Name:      input.Revision.Name,
			Version:   input.Revision.Version,
		},
		Table:      input.Table,
		Operations: reversedOperations,
	}, nil
}

func (e *Executor) reverseOperations(input []operation.Operation) ([]operation.Operation, error) {
	reversedOperations := make([]operation.Operation, 0, len(input))

	var index = len(input) - 1

	for index >= 0 {
		reversed, err := e.reverseOperation(input[index])
		if err != nil {
			return nil, err
		}
		reversedOperations = append(reversedOperations, reversed)
		index -= 1
	}

	return reversedOperations, nil
}

func (e *Executor) reverseOperation(input operation.Operation) (operation.Operation, error) {
	var reversed = operation.Operation{}
	var err error

	if input.AddField != nil {
		reversed.DropField = &operation.DropField{
			TableName: input.AddField.TableName,
			Field:     input.AddField.Column.Name,
		}
	} else if input.CreateTable != nil {
		reversed.DropTable = &operation.DropTable{TableName: input.CreateTable.TableName}
	} else if input.AlterField != nil {
		// this one is a bit trickier. we need to look into the db state, see what's the column version and
		// generate alter to the version smaller than this one.
		reversed.AlterField = &operation.AlterField{
			TableName:    input.AlterField.TableName,
			FieldName:    input.AlterField.FieldName,
			FieldVersion: input.AlterField.FieldVersion,
		}

		if err = e.dbState.Delta(reversed.AlterField, reversed.AlterField.FieldVersion, reversed.AlterField.FieldVersion-1); err != nil {
			return reversed, err
		}

	} else if input.CreateIndex != nil {
		reversed.DropIndex = &operation.DropIndex{
			Table:     input.CreateIndex.Table,
			IndexName: input.CreateIndex.Index.Name,
		}
	} else if input.RunStatement != nil {
		reversed.RunStatement = &operation.RunStatement{
			DDL: input.RunStatement.ReverseDDL,
		}
		if input.RunStatement.ForwardsOnly {
			return reversed, errors.New("this operation cannot be reversed")
		}
	} else if input.AdjustTableVersion != nil {
		reversed.AdjustTableVersion = &operation.AdjustTableVersion{
			TableName: input.AdjustTableVersion.TableName,
			Delta:     -1,
		}
	} else if input.AddConstraint != nil {
		reversed.DropConstraint = &operation.DropConstraint{
			Table:      input.AddConstraint.Table,
			Constraint: input.AddConstraint.Constraint.Name,
		}
	} else {
		err = errors.New("unable to reverse operation")
	}

	return reversed, err
}
