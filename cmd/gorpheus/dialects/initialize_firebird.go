//go:build firebird || all

package dialects

import (
	"github.com/toudi/gorpheus/v2/dialect"
	"github.com/toudi/gorpheus/v2/dialect/firebird"
)

func init() {
	dialect.Register("firebird", firebird.New)
}
