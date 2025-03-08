package sqlite3

import (
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func (s3 *sqlite3Dialect) unpackAlterField(af *operation.AlterField) ([]operation.Operation, error) {
	tempTableOps, err := s3.createTempTable(af.TableName, func(columnName string) bool {
		return columnName != af.FieldName
	})
	if err != nil {
		return nil, err
	}

	addNewField := &operation.AddField{
		TableName: tempTableOps.createTempTable.CreateTable.TableName,
		Column:    af.Column,
	}

	var queue []operation.Operation

	// now we can continue with the algorith, by appending the following operations
	// to the post process queue:
	queue = append(
		queue,
		// 1. create new table
		tempTableOps.createTempTable,
		// 2. add field to the newly created temp table
		operation.Operation{
			AddField: addNewField,
		},
		// 3. populate temp table
		tempTableOps.insertIntoSelectFrom,
		// 4. remove "old" table (which is actually the live one)
		tempTableOps.dropOriginal,
		// 5. rename temp table into "old" table
		tempTableOps.renameTempTable,
		// 6. hope for the best, but there isn't any operation that we can encode it with.
	)
	if tempTableOps.sideEffects != nil {
		queue = append(queue, tempTableOps.sideEffects...)
	}

	return s3.operationsWithPotentialIndex(queue, af.TableName)
}
