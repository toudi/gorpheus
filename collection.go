package gorpheus

import (
	"errors"

	"github.com/gobuffalo/fizz"
	"github.com/toudi/gorpheus/migration"
)

type Migrations []migration.MigrationI

type Collection struct {
	Versions       Migrations
	FizzTranslator fizz.Translator
	metadata       map[string]NamespaceMeta
	applied        map[string]bool
	loggers        map[uint]Logger
}

var ErrNoSuchVersion = errors.New("no such migration")

func Collection_init() *Collection {
	return &Collection{
		Versions: make(Migrations, 0),
		metadata: make(map[string]NamespaceMeta),
		applied:  make(map[string]bool),
		loggers:  make(map[uint]Logger),
	}
}

func (c *Collection) SetTranslator(translator fizz.Translator) {
	c.Log(LoggerGorpheus, LogLevelDebug, "set translator to %v", translator)
	c.FizzTranslator = translator
}

func (c *Collection) TranslatedSQL(sql string) (string, error) {
	return fizz.NewBubbler(c.FizzTranslator).Bubble(sql)
}

func (c *Collection) Register(m migration.MigrationI) {
	c.Log(LoggerGorpheus, LogLevelDebug, "registering migration: %+v", m)

	m.SetDependencies()
	c.Versions = append(c.Versions, m)
}
