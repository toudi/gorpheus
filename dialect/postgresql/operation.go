package postgresql

import (
	"strings"

	"github.com/toudi/gorpheus/v2/dialect/generic"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func (pg *pgDialect) OperationDDL(o operation.Operation, buffer *strings.Builder) error {
	if o.AlterField != nil {
		return pg.AlterFieldDDL(o.AlterField, buffer)
	}
	if o.AddConstraint != nil {
		return pg.AddConstraintDDL(o.AddConstraint, buffer)
	}
	return generic.OperationDDL(o, pg.ColumnDDL, buffer)
}

func (pg *pgDialect) Unpack(o operation.Operation) ([]operation.Operation, error) {
	// postgresql does not need operation unpacking
	return []operation.Operation{o}, nil
}
