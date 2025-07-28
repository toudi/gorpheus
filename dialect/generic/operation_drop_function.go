package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericDropFunction(
	df *operation.DropFunction,
	buffer *strings.Builder,
) error {
	buffer.WriteString("DROP FUNCTION IF EXISTS ")
	buffer.WriteString(df.Name)
	return nil
}
