package sqlite3

import (
	"strings"

	"github.com/toudi/gorpheus/v2/dialect/generic"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func (s3 *sqlite3Dialect) OperationDDL(o operation.Operation, buffer *strings.Builder) error {
	return generic.OperationDDL(o, s3.ColumnDDL, buffer)
}
