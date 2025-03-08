package postgresql

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/toudi/gorpheus/v2/interfaces"
	"github.com/toudi/gorpheus/v2/internal/ddl/column"
	"github.com/toudi/gorpheus/v2/internal/ddl/constraint"
)

var (
	ErrUnsupportedColumnType = errors.New("unknown column type")
	ErrUnknownPolicy         = errors.New("unknown constraint policy")
)

type pgDialect struct {
	connection *sql.DB
	dsn        string
}

func pgForeignKeyPolicy(policy string, buffer *strings.Builder) error {
	policyConstraint, exists := constraint.PolicyFromAlias[policy]
	if !exists {
		return ErrUnknownPolicy
	}
	switch policyConstraint {
	case constraint.CASCADE:
		buffer.WriteString("CASCADE")
	case constraint.NO_ACTION:
		buffer.WriteString("NO ACTION")
	case constraint.SET_DEFAULT:
		buffer.WriteString("SET DEFAULT")
	case constraint.SET_NULL:
		buffer.WriteString("SET NULL")
	}
	return nil
}

func New(_ interfaces.DBState, dsn string) interfaces.Dialect {
	return &pgDialect{
		dsn: dsn,
	}
}

func (pd *pgDialect) columnTypeDDL(c *column.Column, buffer *strings.Builder) error {
	switch c.Type {
	case column.TypeBoolean:
		buffer.WriteString("BOOLEAN")
	case column.TypeDate:
		fallthrough
	case column.TypeTime:
		fallthrough
	case column.TypeTimestamp:
		return pd.columnTypeTimeBasedDDL(c, buffer)
	case column.TypeSmallInt:
		buffer.WriteString("SMALLINT")
	case column.TypeInt:
		buffer.WriteString("INTEGER")
	case column.TypeBigInt:
		buffer.WriteString("BIGINT")
	case column.TypeDecimal:
		buffer.WriteString("NUMERIC")
		// according to docs, it is still valid to create a numeric column without precision or scale.
		// if precision is given, then it must be positive, whereas the scale can be negative.
		// for details, please refer to:
		// https://www.postgresql.org/docs/current/datatype-numeric.html#DATATYPE-NUMERIC-DECIMAL
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
	case column.TypeSmallSerial:
		buffer.WriteString("SMALLSERIAL")
	case column.TypeSerial:
		buffer.WriteString("SERIAL")
	case column.TypeBigSerial:
		buffer.WriteString("BIGSERIAL")
	case column.TypeString:
		buffer.WriteString("VARCHAR (")
		buffer.WriteString(strconv.Itoa(c.MaxLength))
		buffer.WriteString(")")
	case column.TypeText:
		buffer.WriteString("TEXT")
	default:
		return errors.Join(ErrUnsupportedColumnType, errors.New(c.TypeAlias))
	}

	return nil
}
