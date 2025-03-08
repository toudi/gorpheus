package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericDropIndex(di *operation.DropIndex, buffer *strings.Builder) error {
	buffer.WriteString("DROP INDEX ")
	buffer.WriteString(di.IndexName)
	return nil
}
