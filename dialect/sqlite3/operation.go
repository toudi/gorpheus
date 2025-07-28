package sqlite3

import (
	"errors"
	"strings"

	"github.com/toudi/gorpheus/v2/dialect/generic"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

var (
	ErrMaterializedViewNotSupported = errors.New("materialized view not supported on sqlite3")
)

func (s3 *sqlite3Dialect) OperationDDL(o operation.Operation, buffer *strings.Builder) error {
	if o.CreateView != nil && o.CreateView.Materialized {
		return ErrMaterializedViewNotSupported
	}
	return generic.OperationDDL(o, s3.ColumnDDL, buffer)
}
