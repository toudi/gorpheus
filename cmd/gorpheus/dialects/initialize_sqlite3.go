//go:build sqlite3 || all

package dialects

import (
	"github.com/toudi/gorpheus/v2/dialect"
	"github.com/toudi/gorpheus/v2/dialect/sqlite3"
)

func init() {
	dialect.Register("sqlite3", sqlite3.New)
}
