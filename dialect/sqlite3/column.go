package sqlite3

import (
	"errors"
	"fmt"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"

	"github.com/samber/lo"
)

var (
	ErrUnsupportedTypeForAutoIncrement = errors.New("Unsupported type for auto-increment property")
)

func (s3 *sqlite3Dialect) ColumnDDL(c *column.Column, buffer *strings.Builder) error {
	if err := s3.validateColumn(c); err != nil {
		return err
	}
	buffer.WriteString(c.Name)
	buffer.WriteString(" ")
	if err := s3.columnTypeDDL(c, buffer); err != nil {
		return err
	}
	buffer.WriteString(" ")
	if c.Null == nil && c.Generated == nil {
		c.Null = lo.ToPtr(false)
	}
	if c.Generated != nil {
		buffer.WriteString("GENERATED ALWAYS AS (")
		buffer.WriteString(c.Generated.Definition)
		buffer.WriteString(") ")
		if c.Generated.Stored {
			buffer.WriteString("STORED")
		} else {
			buffer.WriteString("VIRTUAL")
		}
	}
	if c.PrimaryKey != nil && *c.PrimaryKey {
		buffer.WriteString("PRIMARY KEY ")
	}
	if c.AutoIncrement != nil && *c.AutoIncrement {
		buffer.WriteString("AUTOINCREMENT ")
	}
	if c.Null != nil {
		if !*c.Null {
			buffer.WriteString("NOT ")
		}
		buffer.WriteString("NULL")
	}
	if c.Default != nil {
		buffer.WriteString(fmt.Sprintf(" DEFAULT %v", c.Default))
	}
	return nil
}

func (s3 *sqlite3Dialect) validateColumn(c *column.Column) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.AutoIncrement != nil && *c.AutoIncrement {
		// let's switch auto-increment integer fields to serial fields.
		switch c.Type {
		case column.TypeSmallInt:
		case column.TypeInt:
		case column.TypeBigInt:
		case column.TypeSmallSerial:
			c.Type = column.TypeSmallInt
		case column.TypeSerial:
			c.Type = column.TypeInt
		case column.TypeBigSerial:
			c.Type = column.TypeBigInt
		default:
			return ErrUnsupportedTypeForAutoIncrement
		}
	}

	return nil
}
