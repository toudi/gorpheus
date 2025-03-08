package interfaces

import (
	"database/sql"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

type Dialect interface {
	ColumnDDL(column *column.Column, buffer *strings.Builder) error
	OperationDDL(operation operation.Operation, buffer *strings.Builder) error
	// if the engine determines that an operation requires 1 .. n parts to be completed - this method will return them.
	Unpack(operation operation.Operation) ([]operation.Operation, error)
	HistoryManager() MigrationsHistory
	// transaction executor
	Transaction(handler func(transaction *sql.Tx) error) error
}

type DialectConstructor func(dbState DBState, dsn string) Dialect
