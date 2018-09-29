package dialect

import (
	"bytes"

	"github.com/toudi/gorpheus"
)

type DialectFactory func(body bytes.Buffer) gorpheus.MigrationI
