package postgresql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"

	"github.com/samber/lo"
)

var (
	ErrPostgreSQLDoesNotSupportVirtualColumns = errors.New("PostgreSQL does not support non-stored generated columns")
	ErrUnsupportedTypeForAutoIncrement        = errors.New("Unsupported type for auto-increment property")
)

func (pd *pgDialect) ColumnDDL(c *column.Column, buffer *strings.Builder) error {
	if err := pd.validateColumn(c); err != nil {
		return err
	}
	buffer.WriteString(c.Name)
	buffer.WriteString(" ")
	if err := pd.columnTypeDDL(c, buffer); err != nil {
		return err
	}
	buffer.WriteString(" ")
	if c.Null == nil {
		c.Null = lo.ToPtr(false)
	}
	if c.PrimaryKey != nil && *c.PrimaryKey {
		buffer.WriteString("PRIMARY KEY ")
	}
	if c.Generated != nil {
		buffer.WriteString(" ")
		buffer.WriteString("GENERATED ALWAYS AS (")
		buffer.WriteString(c.Generated.Definition)
		buffer.WriteString(") STORED")
	}
	if !*c.Null {
		buffer.WriteString(" NOT")
	}
	buffer.WriteString(" NULL")
	if c.Default != nil {
		buffer.WriteString(fmt.Sprintf(" DEFAULT %v", c.Default))
	}
	return nil
}

func (pd *pgDialect) validateColumn(c *column.Column) error {
	if err := c.Validate(); err != nil {
		return err
	}

	if c.Generated != nil && !c.Generated.Stored {
		return ErrPostgreSQLDoesNotSupportVirtualColumns
	}

	if c.AutoIncrement != nil && *c.AutoIncrement {
		// let's switch auto-increment integer fields to serial fields.
		switch c.Type {
		case column.TypeSmallInt:
			c.Type = column.TypeSerial
		case column.TypeInt:
			c.Type = column.TypeSerial
		case column.TypeBigInt:
			c.Type = column.TypeBigSerial
		case column.TypeSmallSerial:
		case column.TypeSerial:
		case column.TypeBigSerial:
			// if the column is already of type serial, then we don't need to do anything.
		default:
			return ErrUnsupportedTypeForAutoIncrement
		}
	}

	return nil
}
