package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericRunDDL(rd *operation.RunStatement, buffer *strings.Builder) error {
	buffer.WriteString(rd.DDL)
	return nil
}
