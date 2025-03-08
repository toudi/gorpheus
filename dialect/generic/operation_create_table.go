package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericCreateTable(
	ct *operation.CreateTable,
	columnDDL func(c *column.Column, buffer *strings.Builder) error,
	buffer *strings.Builder,
) error {
	buffer.WriteString("CREATE TABLE ")
	buffer.WriteString(ct.TableName)
	buffer.WriteString(" (\n")

	// we can iterate over the defined columns
	numColumns := len(ct.Columns)
	var lastItem = numColumns - 1
	if len(ct.Constraints) > 0 {
		lastItem = numColumns
	}
	for colNo, column := range ct.Columns {
		buffer.WriteString("\t")
		if err := columnDDL(column, buffer); err != nil {
			return err
		}
		if colNo < lastItem {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}

	for i, constraint := range ct.Constraints {
		buffer.WriteString("CONSTRAINT ")
		buffer.WriteString(constraint.Name)

		if constraint.ForeignKey != nil {
			buffer.WriteString(" FOREIGN KEY (")
			buffer.WriteString(strings.Join(constraint.ForeignKey.SrcColumns, ", "))
			buffer.WriteString(") REFERENCES ")
			buffer.WriteString(constraint.ForeignKey.Table)
			buffer.WriteString(" (")
			buffer.WriteString(strings.Join(constraint.ForeignKey.Columns, ", "))
			buffer.WriteString(")")

			if constraint.ForeignKey.OnUpdate != nil {
				buffer.WriteString(" ON UPDATE ")
				if err := ConstraintAction(constraint.ForeignKey.OnUpdate, buffer); err != nil {
					return err
				}
			}
			if constraint.ForeignKey.OnDelete != nil {
				buffer.WriteString(" ON DELETE ")
				if err := ConstraintAction(constraint.ForeignKey.OnDelete, buffer); err != nil {
					return err
				}
			}
		}
		if constraint.Unique != nil {
			buffer.WriteString(" UNIQUE (")
			buffer.WriteString(strings.Join(constraint.Unique.Columns, ", "))
			buffer.WriteString(")")
		}
		if i < len(ct.Constraints)-1 {
			buffer.WriteString(",\n")
		}
	}

	buffer.WriteString(")")
	return nil
}
