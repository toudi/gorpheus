package executor

import (
	"github.com/toudi/gorpheus/v2/internal/ddl/column"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
	"github.com/toudi/gorpheus/v2/internal/migration"
)

func (e *Executor) PrepareFinalOperations(m *migration.Migration) error {
	var tableName string
	var err error

	for _, operation := range m.Operations {
		if operation.CreateTable != nil {
			tableName = operation.CreateTable.TableName
			// the input syntax is a map between field name and column
			// therefore we have to construct the columns array.
			if err := e.parseFieldsMap(operation.CreateTable.FieldsMap, func(column *column.Column, rawColumnData map[string]interface{}) {
				operation.CreateTable.Columns = append(
					operation.CreateTable.Columns, column,
				)
				// let's also register the initial column state
				e.dbState.AddColumnVersion(
					operation.CreateTable.TableName,
					column.Name,
					rawColumnData,
				)
			}); err != nil {
				return err
			}
			// we need to register the columns in the db state so that it has
			// the initial parameters of them.
		}
		if operation.AddField != nil {
			if operation.AddField.TableName == "" {
				operation.AddField.TableName = m.Table
			}
			tableName = operation.AddField.TableName
			// the input may either be a single field or a map between
			// field name and Column type
			if err := e.parseFieldsMap(operation.AddField.FieldsMap, func(column *column.Column, rawColumnData map[string]interface{}) {
				operation.AddField.Columns = append(
					operation.AddField.Columns, column,
				)
				e.dbState.AddColumnVersion(operation.AddField.TableName, column.Name, rawColumnData)
			}); err != nil {
				return err
			}
			// we can still have a single field:
			if operation.AddField.Field != nil {
				column, err := e.parseField(operation.AddField.FieldName, *operation.AddField.Field)
				if err != nil {
					return err
				}
				operation.AddField.Columns = append(
					operation.AddField.Columns, column,
				)
				e.dbState.AddColumnVersion(operation.AddField.TableName, operation.AddField.FieldName, *operation.AddField.Field)
			}
		}
		if operation.AlterField != nil {
			if operation.AlterField.TableName == "" {
				operation.AlterField.TableName = m.Table
			}
			// within the operation we also register the "current" version of the column so that
			// when the executor runs operations and is aware of the direction (either forwards or backwards)
			// the db state can yield appropriate delta.
			// this is because we're merely during the parsing stage and at this point we have no idea
			// what will be the migration plan.
			operation.AlterField.FieldVersion = e.dbState.GetColumnVersion(operation.AlterField.TableName,
				operation.AlterField.FieldName,
			)
			// we're doing the alter field. let's register this with the db state
			// so that it can yield the appropriate delta.
			e.dbState.AddColumnVersion(
				operation.AlterField.TableName,
				operation.AlterField.FieldName,
				operation.AlterField.FieldDelta,
			)
		}
		if operation.CreateIndex != nil {
			if operation.CreateIndex.Table == "" {
				operation.CreateIndex.Table = m.Table
			}
		}
		if operation.AddConstraint != nil {
			if operation.AddConstraint.Table == "" {
				operation.AddConstraint.Table = m.Table
			}
			if err = operation.AddConstraint.Constraint.Validate(); err != nil {
				return err
			}
		}
		if tableName != "" {
			e.dbState.IncrementTableVersion(tableName)
		}
	}

	if m.Operations, err = operation.UnpackOperations(m.Operations, true); err != nil {
		return err
	}

	for _, operation := range m.Operations {
		if err = operation.Validate(); err != nil {
			return err
		}
	}

	return nil
}
