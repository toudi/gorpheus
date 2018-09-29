package gorpheus

import (
	"errors"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus/migration"
)

const (
	DirectionUp   = iota
	DirectionDown = iota
)

type Migrations []migration.MigrationI

type Collection struct {
	Versions Migrations
}

var NoSuchVersionErr = errors.New("No such migration")

func Collection_init() *Collection {
	return &Collection{
		Versions: make(Migrations, 0),
	}
}

func (c *Collection) Register(m migration.MigrationI) {
	log.Debugf("Registering migration: %+v", m)

	m.SetDependencies()
	c.Versions = append(c.Versions, m)
}

func (c *Collection) Sort() {
	sort.Sort(c.Versions)
}
