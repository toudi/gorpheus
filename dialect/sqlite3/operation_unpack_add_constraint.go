package sqlite3

import (
	"github.com/toudi/gorpheus/v2/internal/ddl/constraint"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func (s3 *sqlite3Dialect) unpackAddConstaint(ac *operation.AddConstraint) ([]operation.Operation, error) {
	createTempTableOps, err := s3.createTempTable(
		ac.Table,
		func(columnName string) bool { return true },
	)
	if err != nil {
		return nil, err
	}
	createTempTableOps.createTempTable.CreateTable.Constraints = []*constraint.Constraint{ac.Constraint}

	var queue []operation.Operation

	// now we can continue with the algorith, by appending the following operations
	// to the post process queue:
	queue = append(
		queue,
		// 1. create new table
		createTempTableOps.createTempTable,
		// 2. populate temp table
		createTempTableOps.insertIntoSelectFrom,
		// 3. remove "old" table (which is actually the live one)
		createTempTableOps.dropOriginal,
		// 4. rename temp table into "old" table
		createTempTableOps.renameTempTable,
		// 5. hope for the best, but there isn't any operation that we can encode it with.
	)

	return s3.operationsWithPotentialIndex(queue, ac.Table)
}
