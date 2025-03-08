package generic

import (
	"fmt"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func OperationDDL(
	operation operation.Operation,
	columnDDL func(c *column.Column, buffer *strings.Builder) error,
	buffer *strings.Builder,
) error {
	if operation.CreateTable != nil {
		return genericCreateTable(operation.CreateTable, columnDDL, buffer)
	}
	if operation.DropTable != nil {
		return genericDropTable(operation.DropTable, buffer)
	}
	if operation.CreateIndex != nil {
		return genericCreateIndex(operation.CreateIndex, buffer)
	}
	if operation.DropIndex != nil {
		return genericDropIndex(operation.DropIndex, buffer)
	}
	if operation.AddField != nil {
		return genericAddField(operation.AddField, columnDDL, buffer)
	}
	if operation.DropField != nil {
		return genericDropField(operation.DropField, buffer)
	}
	if operation.RunStatement != nil {
		return genericRunDDL(operation.RunStatement, buffer)
	}
	if operation.DropConstraint != nil {
		return genericDropConstraint(operation.DropConstraint, buffer)
	}
	return fmt.Errorf("unhandled operation: %+v", operation)
}
