package gorpheus

import (
	"errors"
	"fmt"
	"sort"

	"github.com/gobuffalo/fizz"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus/v1/migration"
)

type Migrations []migration.MigrationI

type Collection struct {
	Versions       Migrations
	FizzTranslator fizz.Translator
	SQLBindType    int
}

var ErrNoSuchVersion = errors.New("no such migration")

func Collection_init() *Collection {
	return &Collection{
		Versions:    make(Migrations, 0),
		SQLBindType: sqlx.UNKNOWN,
	}
}

func (c *Collection) SetTranslator(translator fizz.Translator) {
	fmt.Printf("set translator to %v", translator)
	c.FizzTranslator = translator
}

func (c *Collection) TranslatedSQL(sql string) (string, error) {
	return fizz.NewBubbler(c.FizzTranslator).Bubble(sql)
}

func (c *Collection) Register(m migration.MigrationI) {
	log.Debugf("Registering migration: %+v", m)

	m.SetDependencies()
	c.Versions = append(c.Versions, m)
}

func (c *Collection) Sort() {
	sort.Sort(c.Versions)
}
