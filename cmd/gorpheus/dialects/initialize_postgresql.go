//go:build postgresql || all

package dialects

import (
	"github.com/toudi/gorpheus/v2/dialect"
	"github.com/toudi/gorpheus/v2/dialect/postgresql"
)

func init() {
	dialect.Register("postgres", postgresql.New)
}
