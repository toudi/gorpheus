package sqlite3

import (
	"github.com/toudi/gorpheus/v2/internal/ddl/constraint"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"

	"github.com/samber/lo"
)

func (s3 *sqlite3Dialect) unpackDropConstaint(dc *operation.DropConstraint) ([]operation.Operation, error) {
	createTempTableOps, err := s3.createTempTable(
		dc.Table,
		func(columnName string) bool { return true },
	)
	if err != nil {
		return nil, err
	}
	sideEffects, err := createTempTableOps.createTempTable.CreateTable.SideEffects()
	if err != nil {
		return nil, err
	}
	createTempTableOps.createTempTable.CreateTable.Constraints = lo.Reject(
		createTempTableOps.createTempTable.CreateTable.Constraints,
		func(c *constraint.Constraint, _ int) bool { return c.Name == dc.Constraint },
	)

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
	queue = append(queue, sideEffects...)

	return s3.operationsWithPotentialIndex(queue, dc.Table)
}
