package firebird

import (
	"errors"
	"fmt"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"

	"github.com/samber/lo"
)

var (
	ErrFirebirdDoesNotSupportStoredGeneratedColumns = errors.New("FirebirdSQL does not support stored generated columns")
	ErrUnsupportedTypeForAutoIncrement              = errors.New("Unsupported type for auto-increment property")
)

func (fb *fbDialect) ColumnDDL(c *column.Column, buffer *strings.Builder) error {
	if err := fb.validateColumn(c); err != nil {
		return err
	}
	buffer.WriteString(c.Name)
	buffer.WriteString(" ")
	if err := fb.columnTypeDDL(c, buffer); err != nil {
		return err
	}
	// buffer.WriteString(" ")
	if c.Null == nil && c.Generated == nil {
		// computed column does not support NOT NULL constraint.
		c.Null = lo.ToPtr(true)
	}
	if c.PrimaryKey != nil && *c.PrimaryKey {
		buffer.WriteString(" PRIMARY KEY")
	}
	if c.Generated != nil {
		// firebird does support generated always, however the documentation
		// specifies that they are only generated on insert which is most likely not what
		// you want.
		buffer.WriteString("COMPUTED BY (")
		buffer.WriteString(c.Generated.Definition)
		buffer.WriteString(")")
	}

	// in firebird SQL, a column is nullable by default, unless specified otherwise.
	if c.Null != nil && !*c.Null {
		buffer.WriteString(" NOT NULL")
	}
	if c.Default != nil {
		buffer.WriteString(fmt.Sprintf(" DEFAULT %v", c.Default))
	}
	return nil
}

func (fb *fbDialect) validateColumn(c *column.Column) error {
	if c.AutoIncrement != nil && *c.AutoIncrement {
		// let's switch auto-increment integer fields to serial fields.
		switch c.Type {
		case column.TypeSmallSerial:
			c.Type = column.TypeSmallInt
		case column.TypeSerial:
			c.Type = column.TypeInt
		case column.TypeBigSerial:
			c.Type = column.TypeBigInt
		case column.TypeSmallInt:
		case column.TypeInt:
		case column.TypeBigInt:
			// these types are already valid types for auto-increment so we
			// don't have to do anything.
		default:
			return ErrUnsupportedTypeForAutoIncrement
		}
	}

	return nil
}
