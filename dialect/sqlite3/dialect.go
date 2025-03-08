package sqlite3

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

type sqlite3Dialect struct {
	dbState    interfaces.DBState
	dsn        string
	connection *sql.DB
}

func New(dbState interfaces.DBState, dsn string) interfaces.Dialect {
	return &sqlite3Dialect{
		dbState: dbState,
		dsn:     dsn,
	}
}

func (s3 *sqlite3Dialect) columnTypeDDL(c *column.Column, buffer *strings.Builder) error {
	switch c.Type {
	case column.TypeBoolean:
		buffer.WriteString("BOOLEAN")
	case column.TypeSmallInt:
		buffer.WriteString("SMALLINT")
	case column.TypeInt:
		buffer.WriteString("INTEGER")
	case column.TypeBigInt:
		buffer.WriteString("BIGINT")
	case column.TypeDecimal:
		buffer.WriteString("NUMERIC")
		if c.MaxDigits > 0 {
			buffer.WriteString(" (")
			buffer.WriteString(strconv.Itoa(c.MaxDigits))
			if c.DecimalPlaces != 0 {
				buffer.WriteString(", ")
				buffer.WriteString(strconv.Itoa(c.DecimalPlaces))
			}
			buffer.WriteString(")")
		}
	case column.TypeFloat:
		buffer.WriteString("REAL")
	case column.TypeDouble:
		buffer.WriteString("DOUBLE PRECISION")
	case column.TypeString:
		buffer.WriteString("VARCHAR (")
		buffer.WriteString(strconv.Itoa(c.MaxLength))
		buffer.WriteString(")")
	case column.TypeText:
		buffer.WriteString("TEXT")
	case column.TypeDate:
		buffer.WriteString("DATE")
	case column.TypeTime:
		buffer.WriteString("TIME")
	case column.TypeTimestamp:
		buffer.WriteString("DATETIME")
	default:
		return errors.Join(ErrUnsupportedColumnType, errors.New(c.TypeAlias))
	}

	return nil
}
