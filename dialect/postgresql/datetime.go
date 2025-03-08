package postgresql

import (
	"strconv"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"
)

func (pd *pgDialect) columnTypeTimeBasedDDL(c *column.Column, buffer *strings.Builder) error {
	if c.Type == column.TypeTime {
		buffer.WriteString("TIME")
		precisionInfo(c, buffer)
	}
	if c.Type == column.TypeDate {
		buffer.WriteString("DATE")
	}
	if c.Type == column.TypeTimestamp {
		buffer.WriteString("TIMESTAMP")
		precisionInfo(c, buffer)
	}
	return nil
}

func precisionInfo(c *column.Column, buffer *strings.Builder) {
	if c.Precision > 0 {
		buffer.WriteString(" (")
		buffer.WriteString(strconv.Itoa(int(c.Precision)))
		buffer.WriteString(")")
	}
	if c.WithTimeZone {
		buffer.WriteString(" WITH TIME ZONE")
	}
}
