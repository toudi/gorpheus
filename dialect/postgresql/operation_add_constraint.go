package postgresql

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func (pg *pgDialect) AddConstraintDDL(ac *operation.AddConstraint, buffer *strings.Builder) error {
	buffer.WriteString("ALTER TABLE ")
	buffer.WriteString(ac.Table)
	buffer.WriteString(" ADD CONSTRAINT ")
	buffer.WriteString(ac.Constraint.Name)
	buffer.WriteString(" ")
	if ac.Constraint.ForeignKey != nil {
		buffer.WriteString("FOREIGN KEY (")
		buffer.WriteString(strings.Join(ac.Constraint.ForeignKey.SrcColumns, ", "))
		buffer.WriteString(") REFERENCES ")
		buffer.WriteString(ac.Constraint.ForeignKey.Table)
		buffer.WriteString("(")
		buffer.WriteString(strings.Join(ac.Constraint.ForeignKey.Columns, ", "))
		buffer.WriteString(")")
	}

	if ac.Constraint.Unique != nil {
		buffer.WriteString("UNIQUE ")
		if ac.Constraint.Unique.Null != nil && *ac.Constraint.Unique.Null {
			buffer.WriteString("NULLS ")
		}
		if ac.Constraint.Unique.Distinct != nil {
			if !*ac.Constraint.Unique.Distinct {
				buffer.WriteString("NOT ")
			}
			buffer.WriteString("DISTINCT ")
		}
		buffer.WriteString("(")
		buffer.WriteString(strings.Join(ac.Constraint.Unique.Columns, ", "))
		buffer.WriteString(")")
	}

	return nil
}
