package gorpheus

import (
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus/dialect"
	"github.com/toudi/gorpheus/migration"
)

func (c *Collection) MigrateUp(tx *sqlx.Tx) error {
	// migrate towards latest version
	return c.MigrateUpTo(tx, "")
}

func (c *Collection) performMigrate(tx *sqlx.Tx, version string, direction uint8) error {
	sort.Sort(c.Versions)

	// find the index of target version
	var targetIdx int
	var err error
	var operationDesc = map[uint8]string{
		DirectionUp:   "applying",
		DirectionDown: "unapplying",
	}
	var breakLoop = false
	var revIdx int
	var m migration.MigrationI

	if version == "" {
		targetIdx = len(c.Versions)
	} else {
		targetIdx, err = c.FindMigration(version)
		if err != nil {
			log.Errorf("Cannot find target migration")
			return err
		}
	}

	for idx := range c.Versions {
		if direction == DirectionDown {
			revIdx = len(c.Versions) - idx - 1
			m = c.Versions[revIdx]
			breakLoop = (revIdx <= targetIdx)
		} else {
			m = c.Versions[idx]
			breakLoop = (idx >= targetIdx)
		}

		if !breakLoop {
			exists, err := c.Exists(tx, m.Revision())
			if err != nil {
				log.WithError(err).Errorf("Cannot check if %s exists.", m.Revision())
				return err
			}
			if exists == true && direction == DirectionUp {
				log.Debugf("Skipping %s - already migrated", m.Revision())
				continue
			}
			log.Debugf("%s %s", operationDesc[direction], m.Revision())
			if direction == DirectionUp {
				// get UP script and it's type.
				upScript, upType, err := m.UpScript()
				if upType == migration.TypeFizz {
					upScript, err = dialect.FizzDecode(upScript)
				}
				_, err = tx.Exec(upScript)
				if err == nil {
					err = c.InsertRevision(tx, m.Revision())
				}
			} else {
				downScript, downType, err := m.DownScript()
				if downType == migration.TypeFizz {
					downScript, err = dialect.FizzDecode(downScript)
				}
				_, err = tx.Exec(downScript)
				if err == nil {
					err = c.RemoveRevision(tx, m.Revision())
				}
			}
			if err != nil {
				log.WithError(err).Errorf("Error %s migration %s", operationDesc[direction], m.Revision())
				return err
			}
		} else {
			break
		}
	}

	return nil

}
func (c *Collection) MigrateUpTo(tx *sqlx.Tx, version string) error {
	log.Debugf("] MigrateUpTo::%v", version)
	return c.performMigrate(tx, version, DirectionUp)
}

func (c *Collection) MigrateDownTo(tx *sqlx.Tx, version string) error {
	log.Debugf("] MigrateDownTo::%v", version)
	return c.performMigrate(tx, version, DirectionDown)
}

func (c *Collection) FindMigration(version string) (int, error) {
	for idx, m := range c.Versions {
		if strings.HasPrefix(m.Revision(), version) {
			return idx, nil
		}
	}
	return -1, NoSuchVersionErr
}
