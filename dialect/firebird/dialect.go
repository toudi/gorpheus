package firebird

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/toudi/gorpheus/v2/interfaces"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"
)

var (
	ErrUnsupportedColumnType = errors.New("unknown column type")
)

type fbDialect struct {
	connection *sql.DB
	dsn        string
}

func New(_ interfaces.DBState, dsn string) interfaces.Dialect {
	return &fbDialect{
		dsn: dsn,
	}
}

func (fd *fbDialect) columnTypeDDL(c *column.Column, buffer *strings.Builder) error {
	switch c.Type {
	case column.TypeBoolean:
		buffer.WriteString("BOOLEAN")
	case column.TypeDate:
		fallthrough
	case column.TypeTime:
		fallthrough
	case column.TypeTimestamp:
		return fd.columnTypeTimeBasedDDL(c, buffer)
	case column.TypeSmallInt:
		fallthrough
	case column.TypeInt:
		fallthrough
	case column.TypeBigInt:
		fallthrough
	case column.TypeInt128:
		return fd.integerDDL(c, buffer)
	case column.TypeDecimal:
		buffer.WriteString("DECIMAL")
		// according to docs, it is still valid to create a numeric column without precision or scale.
		// if precision is given, then it must be positive
		// if scale is given, then it must be positive.
		// for details, please refer to:
		// https://firebirdsql.org/file/documentation/html/en/refdocs/fblangref50/firebird-50-language-reference.html#fblangref50-datatypes-numeric
		if c.MaxDigits > 0 {
			buffer.WriteString(" (")
			buffer.WriteString(strconv.Itoa(c.MaxDigits))
			if c.DecimalPlaces > 0 {
				buffer.WriteString(", ")
				buffer.WriteString(strconv.Itoa(c.DecimalPlaces))
			}
			buffer.WriteString(")")
		}
	case column.TypeFloat:
		buffer.WriteString("FLOAT")
	case column.TypeDouble:
		buffer.WriteString("DOUBLE PRECISION")
	case column.TypeString:
		buffer.WriteString("VARCHAR (")
		buffer.WriteString(strconv.Itoa(c.MaxLength))
		buffer.WriteString(")")
	default:
		return errors.Join(ErrUnsupportedColumnType, errors.New(c.TypeAlias))
	}

	return nil
}
