package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericDropView(
	dv *operation.DropView,
	buffer *strings.Builder,
) error {
	if dv.Materialized {
		buffer.WriteString("DROP MATERIALIZED VIEW IF EXISTS ")
	} else {
		buffer.WriteString("DROP VIEW IF EXISTS ")
	}
	buffer.WriteString(dv.Name)
	return nil
}
