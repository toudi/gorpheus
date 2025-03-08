package firebird

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"
)

func (fd *fbDialect) columnTypeTimeBasedDDL(c *column.Column, buffer *strings.Builder) error {
	if c.Type == column.TypeTime {
		buffer.WriteString("TIME")
		timeZoneInfo(c, buffer)
	}
	if c.Type == column.TypeDate {
		buffer.WriteString("DATE")
	}
	if c.Type == column.TypeTimestamp {
		buffer.WriteString("TIMESTAMP")
		timeZoneInfo(c, buffer)
	}
	return nil
}

func timeZoneInfo(c *column.Column, buffer *strings.Builder) {
	if c.WithTimeZone {
		buffer.WriteString(" WITH TIME ZONE")
	}
}
