package firebird

import (
	"fmt"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

// we need to track how many changes we made in order to put the comma.
type columnModifierDDL struct {
	changeCounter int
	columnName    string
}

func (fd *fbDialect) AlterFieldDDL(af *operation.AlterField, buffer *strings.Builder) error {
	var err error

	buffer.WriteString("ALTER TABLE ")
	buffer.WriteString(af.TableName)
	buffer.WriteString(" ")

	modifier := &columnModifierDDL{columnName: af.FieldName}

	// if af.Column is populated then it means that the properties of the type has changed.
	if af.Column != nil {
		if err = modifier.append(buffer, func(buffer *strings.Builder) error {
			buffer.WriteString("TYPE ")
			return fd.columnTypeDDL(af.Column, buffer)
		}); err != nil {
			return err
		}
	}

	// otherwise, it is possible that the null constraint was altered:
	if af.Delta.Null != nil {
		if err = modifier.append(buffer, func(buffer *strings.Builder) error {
			if *af.Delta.Null {
				buffer.WriteString("DROP ")
			} else {
				buffer.WriteString("SET ")
			}
			buffer.WriteString("NOT NULL ")
			return nil
		}); err != nil {
			return err
		}
	}
	// or that the default value was set:
	if af.Delta.Default != nil {
		if err = modifier.append(buffer, func(buffer *strings.Builder) error {
			buffer.WriteString("SET DEFAULT ")
			buffer.WriteString(fmt.Sprintf("%v", af.Delta.Default))
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func (m *columnModifierDDL) append(buffer *strings.Builder, modificationDelta func(buffer *strings.Builder) error) error {
	if m.changeCounter > 0 {
		buffer.WriteString(", ")
	}
	buffer.WriteString("ALTER COLUMN ")
	buffer.WriteString(m.columnName)
	buffer.WriteString(" ")

	err := modificationDelta(buffer)
	m.changeCounter += 1
	return err
}
