package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/interfaces"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func OperationsDDL(operations []operation.Operation, dialect interfaces.Dialect, buffer *strings.Builder) error {
	lastOperationIdx := len(operations) - 1

	for index, _operation := range operations {
		if err := dialect.OperationDDL(_operation, buffer); err != nil {
			return err
		}

		if index != lastOperationIdx {
			buffer.WriteString(";\n")
		}
	}

	return nil
}
