package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericDropConstraint(dc *operation.DropConstraint, buffer *strings.Builder) error {
	buffer.WriteString("ALTER TABLE ")
	buffer.WriteString(dc.Table)
	buffer.WriteString(" DROP CONSTRAINT ")
	buffer.WriteString(dc.Constraint)
	return nil
}
