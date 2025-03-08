package sqlite3

import (
	"slices"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func (s3 *sqlite3Dialect) Unpack(o operation.Operation) ([]operation.Operation, error) {
	// sqlite has a quite specific method for dealing with alter fields - we need to
	// convert them into several sub operations
	if o.AlterField != nil {
		return s3.unpackAlterField(o.AlterField)
	}

	if o.AddConstraint != nil {
		return s3.unpackAddConstaint(o.AddConstraint)
	}

	if o.DropConstraint != nil {
		return s3.unpackDropConstaint(o.DropConstraint)
	}

	// otherwise, it's a regular operation so just return it.
	return []operation.Operation{o}, nil
}

func (s3 *sqlite3Dialect) operationsWithPotentialIndex(queue []operation.Operation, srcTable string) ([]operation.Operation, error) {
	// var err error
	// we also have to unpack the operations
	// we do not want to adjust table version as this is handled by the unpacked alter field itself
	// what we're doing is simply replacing a single operation into several ones.
	// queue, err = operation.UnpackOperations(queue, false)
	// if err != nil {
	// 	return nil, err
	// }

	// ok now comes a slightly odd part which I've learned the hard way. basically, when there's a create table
	// statement, it is quite possible to define indexes. that part isn't actually weird. what's weird is -
	// sqlite3 names indexes uniquely *per schema* rather than per table. what this means is if we're
	// creating the temporary table and then we want to also create the index - we cannot do that. let's explain this
	// as we go along.
	var queueContainsCreateIndexOps bool = false
	for _, enqueued := range queue {
		if enqueued.CreateIndex != nil {
			// we have found the create index. let's change the table name to avoid 'the index already exists'.
			// what's the most important is that the create index statement must be executed *after* the original
			// table is removed. we will deal with this later.
			enqueued.CreateIndex.Table = srcTable
			queueContainsCreateIndexOps = true
		}
	}

	// now for dealing with indexes - we have to shift them as the last operations in the queue which is fairly trivial
	if queueContainsCreateIndexOps {
		slices.SortStableFunc(queue, func(a, b operation.Operation) int {
			if a.CreateIndex != nil {
				// the function should return positive value when a is less than b
				// but since we want to push create index to the very end, we need to
				// reverse it
				return 1
			}
			if b.CreateIndex != nil {
				return -1
			}
			return 0
		})
	}

	return queue, nil
}
