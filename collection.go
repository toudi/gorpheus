package gorpheus

import (
	"errors"
	"sort"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus/migration"
)

type Migrations []migration.MigrationI

type Collection struct {
	Versions Migrations
}

var NoSuchVersionErr = errors.New("No such migration")

// type Collection []MigrationI

func Collection_init() *Collection {
	return &Collection{
		Versions: make(Migrations, 0),
	}
}

func (c *Collection) Register(m migration.MigrationI) {
	log.Debugf("Registering migration: %+v", m)

	m.Parse()
	c.Versions = append(c.Versions, m)
}

func (c *Collection) Sort() {
	sort.Sort(c.Versions)
}

func (c *Collection) MigrateUp(tx *sqlx.Tx) error {
	// migrate towards latest version
	// return c.MigrateUpTo(tx, "")
	return nil
}

func (c *Collection) MigrateUpTo(tx *sqlx.Tx, version string) error {
	// find the index of target version
	// var targetIdx int
	// var err error

	// if version == "" {
	// 	targetIdx = len(c.versions)
	// } else {
	// 	targetIdx, err = c.FindMigration(Id_from_string(version))
	// 	if err != nil {
	// 		log.Errorf("Cannot find target migration")
	// 		return err
	// 	}
	// }

	// for idx, m := range c.versions {
	// 	if idx <= targetIdx {
	// 		log.Debugf("Applying %s", m.ID().ToString())
	// 		if err := m.UpFn(tx); err != nil {
	// 			log.Errorf("Error executing migration %s", m.ID().ToString())
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

func (c *Collection) FindMigration(version *Id) (int, error) {
	// log.Debugf("Looking for %v", version)
	// for idx, m := range c.versions {
	// 	id := m.ID()
	// 	log.Debugf("Comparing %+v, %+v", version, m.ID())
	// 	if id == *version || (id.Namespace == version.Namespace && strings.HasPrefix(id.Version, version.Version)) {
	// 		return idx, nil
	// 	}
	// }

	return -1, NoSuchVersionErr
}

// func Register(m MigrationI) error {
// 	log.Debugf("Registered migration: %s", m.ID().ToString())
// 	Migrations = append(Migrations, m)
// 	return nil
// }

// func MigrationNames(c Collection) []string {
// 	var out = []string{}

// 	for _, val := range c {
// 		out = append(out, val.ID().ToString())
// 	}

// 	return out
// }

// func SortMigrations() error {
// 	sort.Sort(Migrations)
// 	return nil
// }
