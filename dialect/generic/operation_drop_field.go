package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericDropField(df *operation.DropField, buffer *strings.Builder) error {
	buffer.WriteString("ALTER TABLE ")
	buffer.WriteString(df.TableName)
	buffer.WriteString(" DROP ")
	buffer.WriteString(df.Field)

	return nil
}
