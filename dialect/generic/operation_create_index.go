package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericCreateIndex(ci *operation.CreateIndex, buffer *strings.Builder) error {
	buffer.WriteString("CREATE ")
	if ci.Index.Unique {
		buffer.WriteString("UNIQUE ")
	}
	buffer.WriteString("INDEX ")
	buffer.WriteString(ci.Index.GetName())
	buffer.WriteString(" ON ")
	buffer.WriteString(ci.Table)
	if ci.Index.Columns != nil {
		buffer.WriteString(" (")
		buffer.WriteString(strings.Join(ci.Index.Columns, ", "))
		buffer.WriteString(")")
	} else if ci.Index.Definition != "" {
		buffer.WriteString(ci.Index.Definition)
	}
	return nil
}
