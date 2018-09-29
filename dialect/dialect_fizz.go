package dialect

import (
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
)

func FizzDecode(script string) (string, error) {
	bubbler := fizz.NewBubbler(translators.NewSQLite("sqlite://example.sqlite3"))
	return bubbler.Bubble(script)
}
