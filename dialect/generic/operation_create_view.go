package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericCreateView(
	cv *operation.CreateView,
	buffer *strings.Builder,
) error {
	if !cv.Materialized {
		buffer.WriteString("CREATE OR REPLACE VIEW ")
		buffer.WriteString(cv.Name)
		buffer.WriteString(" AS\n")
		buffer.WriteString(cv.Definition)
		return nil
	}
	buffer.WriteString("CREATE MATERIALIZED VIEW IF NOT EXISTS ")
	buffer.WriteString(cv.Name)
	buffer.WriteString(" AS\n")
	buffer.WriteString(cv.Definition)
	return nil
}
