package generic

import (
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

func genericCreateFunction(
	cf *operation.CreateFunction,
	buffer *strings.Builder,
) error {
	buffer.WriteString("CREATE OR REPLACE FUNCTION ")
	buffer.WriteString(cf.Definition)
	return nil
}
