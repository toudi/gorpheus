package gorpheus

import (
	"errors"
	"fmt"
	"sort"

	"github.com/gobuffalo/fizz"
	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus/v1/migration"
)

type Migrations []migration.MigrationI

type Collection struct {
	Versions       Migrations
	FizzTranslator fizz.Translator
	metadata       map[string]NamespaceMeta
	applied        map[string]bool
}

var ErrNoSuchVersion = errors.New("no such migration")

func Collection_init() *Collection {
	return &Collection{
		Versions: make(Migrations, 0),
		metadata: make(map[string]NamespaceMeta),
		applied:  make(map[string]bool),
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
