package operation

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/index"
)

type CreateIndex struct {
	Table string       `yaml:"table"`
	Index *index.Index `yaml:",inline"`
}

func (ci *CreateIndex) DDL(buffer *strings.Builder) error {
	buffer.WriteString("CREATE ")
	if ci.Index.Unique {
		buffer.WriteString("UNIQUE ")
	}
	buffer.WriteString("INDEX IF NOT EXISTS ")
	buffer.WriteString(ci.Index.GetName())
	buffer.WriteString(" ON ")
	buffer.WriteString(ci.Table)
	if ci.Index.Columns != nil {
		buffer.WriteString(" (")
		buffer.WriteString(strings.Join(ci.Index.Columns, ", "))
		buffer.WriteString(")")
	} else if ci.Index.Definition != "" {
		buffer.WriteString(ci.Index.Definition)
	}
	buffer.WriteString(";")
	return nil
}
