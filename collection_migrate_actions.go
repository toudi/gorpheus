package gorpheus

import (
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus/v1/migration"
)

func (c *Collection) MigrateUp(db *sqlx.DB) error {
	// migrate towards latest version
	return c.MigrateUpTo(db, "")
}

func (c *Collection) performMigrate(db *sqlx.DB, version string, direction uint8) error {
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

	var upScript, downScript string
	var upType, downType uint8

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
			tx, err := db.Beginx()
			log.Debugf("begin() err = %v", err)
			exists, err := c.Exists(tx, m.Revision())
			log.Debugf("exists() = %v err = %v", exists, err)
			err = tx.Commit()
			log.Debugf("commit() err = %v", err)
			if err != nil {
				log.WithError(err).Errorf("Cannot check if %s exists.", m.Revision())
				return err
			}
			if exists && direction == DirectionUp {
				log.Debugf("Skipping %s - already migrated", m.Revision())
				continue
			}
			log.Debugf("%s %s", operationDesc[direction], m.Revision())
			log.Infof("Revision begin()")
			tx, err = db.Beginx()
			if err != nil {
				log.WithError(err).Error("Cannot start transaction")
				return err
			}
			if direction == DirectionUp {
				// get UP script and it's type.
				if upScript, upType, err = m.UpScript(); err == nil {
					if upType == migration.TypeFizz {
						if upScript, err = c.TranslatedSQL(upScript); err != nil {
							break
						}
					}
					if upType == migration.TypeGo {
						err = m.Up(tx)
					} else {
						_, err = tx.Exec(upScript)
					}
					if err == nil {
						err = c.InsertRevision(tx, m.Revision())
					}
				}
			} else {
				if downScript, downType, err = m.DownScript(); err == nil {
					if downType == migration.TypeFizz {
						if downScript, err = c.TranslatedSQL(downScript); err != nil {
							break
						}
					}
					if downType == migration.TypeGo {
						if err = m.Down(tx); err != nil {
							break
						}
					} else {
						if _, err = tx.Exec(downScript); err != nil {
							break
						}
					}
					_, err = tx.Exec(downScript)
					if err == nil {
						err = c.RemoveRevision(tx, m.Revision())
					}
				}
			}
			if err != nil {
				log.WithError(err).Errorf("Error %s migration %s", operationDesc[direction], m.Revision())
				err = tx.Rollback()

				return err
			} else {
				log.Info("Revision commit()")
				err = tx.Commit()
				if err != nil {
					log.WithError(err).Error("Could not commit")
					return err
				}
			}
		} else {
			break
		}
	}

	return err

}
func (c *Collection) MigrateUpTo(db *sqlx.DB, version string) error {
	log.Debugf("] MigrateUpTo::%v", version)
	return c.performMigrate(db, version, DirectionUp)
}

func (c *Collection) MigrateDownTo(db *sqlx.DB, version string) error {
	log.Debugf("] MigrateDownTo::%v", version)
	return c.performMigrate(db, version, DirectionDown)
}

func (c *Collection) FindMigration(version string) (int, error) {
	for idx, m := range c.Versions {
		if strings.HasPrefix(m.Revision(), version) {
			return idx, nil
		}
	}
	return -1, ErrNoSuchVersion
}
