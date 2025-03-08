package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericAddField(
	af *operation.AddField,
	columnDDL func(c *column.Column, buffer *strings.Builder) error,
	buffer *strings.Builder,
) error {
	buffer.WriteString("ALTER TABLE ")
	buffer.WriteString(af.TableName)
	buffer.WriteString(" ADD ")
	if err := columnDDL(af.Column, buffer); err != nil {
		return err
	}
	if af.Column.ForeignKey != nil {

	}
	return nil
}
