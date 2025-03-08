package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericDropTable(dt *operation.DropTable, buffer *strings.Builder) error {
	buffer.WriteString("DROP TABLE ")
	buffer.WriteString(dt.TableName)
	return nil
}
