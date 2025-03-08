package firebird

import (
	"strings"

	"github.com/toudi/gorpheus/v2/dialect/generic"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func (fb *fbDialect) OperationDDL(o operation.Operation, buffer *strings.Builder) error {
	if o.AlterField != nil {
		return fb.AlterFieldDDL(o.AlterField, buffer)
	}
	return generic.OperationDDL(o, fb.ColumnDDL, buffer)
}

func (fb fbDialect) Unpack(o operation.Operation) ([]operation.Operation, error) {
	// firebird does not need operation unpacking
	return []operation.Operation{o}, nil
}
