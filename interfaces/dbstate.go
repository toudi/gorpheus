package interfaces

import "github.com/toudi/gorpheus/v2/internal/ddl/column"

type DBState interface {
	TableColumns(tableName string) ([]*column.Column, error)
}
