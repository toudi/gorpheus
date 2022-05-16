package gorpheus

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/toudi/gorpheus/migration"
)

var operationDesc = map[uint8]string{
	DirectionUp:   "applying",
	DirectionDown: "unapplying",
}

type migrationsArray []migration.MigrationI

func (ma migrationsArray) Contains(_migration migration.MigrationI) bool {
	for _, m := range ma {
		if m.GetVersion() == _migration.GetVersion() {
			return true
		}
	}
	return false
}

func (ma migrationsArray) ToString() string {
	output := "["
	for _, m := range ma {
		output += m.GetVersion() + ", "
	}
	output += "]"

	return output
}

func (c *Collection) prepareMigrationsToApply(namespace string, currentVersionNo int, targetVersionNo int, dst *migrationsArray) error {
	c.Log(LoggerDebug, LogLevelDebug, "prepareMigrationsToApply(%s, %d, %d)\n", namespace, currentVersionNo, targetVersionNo)

	var err error
	var delta int = 1

	var versionNumber int

	if currentVersionNo == targetVersionNo {
		return nil
	}

	if currentVersionNo == -1 {
		currentVersionNo = 1
	}

	var _migration migration.MigrationI

	if targetVersionNo < currentVersionNo {
		delta = -1
	}

	namespaceMeta := c.metadata[namespace]

	for {
		_migration = c.Versions[namespaceMeta.positionIndex[currentVersionNo]]
		versionNumber, err = _migration.GetVersionNumber()
		if err != nil {
			return fmt.Errorf("could not get migration version number: %v", err)
		}
		if delta == -1 && versionNumber == targetVersionNo {
			break
		}

		// let's check if this migration is already on the list to apply. if so - we can safely skip.
		if !dst.Contains(_migration) {
			if delta == -1 {
				*dst = append(*dst, _migration)
			}
			dependencies, err := c.GetDependencies(_migration)
			if err != nil {
				return fmt.Errorf("could not get dependencies fro the migration: %v", err)
			}
			c.Log(LogLevelDebug, LogLevelDebug, "dependencies for %s: [", _migration.GetVersion())

			for _, dependency := range dependencies {
				c.Log(LogLevelDebug, LogLevelDebug, "%s, ", dependency.GetVersion())
			}
			c.Log(LogLevelDebug, LogLevelDebug, "]\n")

			for _, dependency := range dependencies {
				if dst.Contains(dependency) {
					continue
				}

				dependencyNamespace, err := dependency.GetNamespace()
				if err != nil {
					return fmt.Errorf("could not parse namespace of %s: %v", dependency.GetVersion(), err)
				}
				dependencyMeta := c.metadata[dependencyNamespace]
				dependencyVersionNo, err := dependency.GetVersionNumber()
				if err != nil {
					return fmt.Errorf("could not get dependency version number: %v", err)
				}
				if dependencyNamespace == namespace && dependencyVersionNo == targetVersionNo {
					continue
				}
				c.prepareMigrationsToApply(dependencyNamespace, dependencyMeta.current, dependencyVersionNo, dst)
			}
			c.Log(LogLevelDebug, LogLevelDebug, "namespace = %s; currentVersionNo = %d; targetVersionNo = %d\n", namespace, currentVersionNo, targetVersionNo)

			if delta == 1 {
				*dst = append(*dst, _migration)
			}
		}
		if currentVersionNo == targetVersionNo {
			break
		}
		currentVersionNo += delta
	}

	return err
}

func (c *Collection) Migrate(params *MigrationParams) error {
	var description string

	var err error
	var db *sqlx.DB
	var currentVersionNo, targetVersionNo int
	// var _migration migration.MigrationI

	if db, err = c.connectToDb(params); err != nil {
		c.Log(LoggerGorpheus, LogLevelError, "unable to connect to the database: %v", err)
		return err
	}

	c.Log(LoggerDebug, LogLevelDebug, "making sure that migrations table exist\n")

	if err = c.ensureMigrationsTableExists(db); err != nil {
		c.Log(LoggerGorpheus, LogLevelError, "cannot create migrations table: %v", err)
		return err
	}

	// check which versions are applied in the database and index the existing collection
	if err = c.index(db); err != nil {
		return fmt.Errorf("could not index migrations: %v", err)
	}

	var direction = DirectionUp
	if params.Vacuum {
		direction = DirectionDown
	}
	if params.Namespace != "" {
		namespaceMeta := c.metadata[params.Namespace]
		if params.Zero || params.VersionNo < namespaceMeta.current {
			direction = DirectionDown
		}
	}

	migrationsToApply := make(migrationsArray, 0)

	c.Log(LoggerDebug, LogLevelDebug, "params: %+v", params)

	for namespace, metadata := range c.metadata {
		if params.Namespace != "" && namespace != params.Namespace {
			continue
		}
		currentVersionNo = metadata.current
		targetVersionNo = metadata.mostRecent
		if params.Vacuum {
			targetVersionNo = -1
		}
		if params.Namespace != "" {
			if params.Zero {
				targetVersionNo = -1
			} else if params.VersionNo > 0 {
				targetVersionNo = params.VersionNo
			}
		}
		c.Log(LoggerDebug, LogLevelDebug, "call to prepareMigrationsToApply(%s, %d, %d)\n", namespace, currentVersionNo, targetVersionNo)

		c.prepareMigrationsToApply(namespace, currentVersionNo, targetVersionNo, &migrationsToApply)
	}

	c.Log(LoggerDebug, LogLevelDebug, "migrations to apply: \n")

	for _, m := range migrationsToApply {
		c.Log(LoggerDebug, LogLevelDebug, "-> %s\n", m.GetVersion())
	}

	description = operationDesc[uint8(direction)]

	err = Atomic(db, func(tx *sqlx.Tx) error {
		for _, _migration := range migrationsToApply {
			if err = c.performMigration(tx, _migration, direction, params.Fake, 0); err != nil {
				return fmt.Errorf("could not %s %s: %v", description[:len(description)-3], _migration.GetVersion(), err)
			}
		}
		return nil
	})

	// cleanup and close all the readers
	for _, m := range c.Versions {
		c.Log(LoggerDebug, LogLevelDebug, "closing %s\n", m.GetVersion())

		if err = m.Close(); err != nil {
			return fmt.Errorf("could not close migration: %v", err)
		}
	}

	if params.Vacuum {
		err = c.dropMigrationsTable(db)
	}

	if err != nil {
		c.Log(LoggerGorpheus, LogLevelError, "could not perform migrations: %v\n", err)
	}

	return err
}

func (c *Collection) performMigration(tx *sqlx.Tx, _migration migration.MigrationI, direction int, fake bool, indent int) error {
	var script string
	var scriptType uint8
	var err error
	var performAction bool = true

	c.Log(LoggerDebug, LogLevelDebug, "%s performMigration(%s, %d)\n", strings.Repeat(" ", indent), _migration.GetVersion(), direction)

	_, exists := c.applied[_migration.GetVersion()]

	// basically, if we're migrating upwards and the migration was already applied then let's
	// skip it.
	// also, if we're migrating downwards and the migrations was not applied then there's no point
	// in unapplying it.
	if (direction == DirectionUp && exists) || (!exists && direction == DirectionDown) {
		c.Log(LoggerDebug, LogLevelDebug, "%s no-op\n", strings.Repeat(" ", indent))
		performAction = false
	}

	if performAction {
		c.Log(LoggerGorpheus, LogLevelInfo, "%s %s %s .. ", strings.Repeat(" ", indent), operationDesc[uint8(direction)], _migration.GetVersion())

		if direction == DirectionUp {
			script, scriptType, err = _migration.UpScript()
		} else {
			script, scriptType, err = _migration.DownScript()
		}

		if err != nil {
			return fmt.Errorf("unable to obtain migration script: %v", err)
		}

		if script == "" && _migration.GetType() != migration.TypeGo {
			return fmt.Errorf("your migration is not a go migration yet it did not return any script. If you want to fake the migration, please use the fake parameter")
		}

		c.Log(LoggerSQL, LogLevelInfo, "script prior to translation: %s\n", script)

		if scriptType == migration.TypeFizz {
			if script, err = c.TranslatedSQL(script); err != nil {
				return fmt.Errorf("cannot translate migration script: %v", err)
			}
		}

		c.Log(LoggerSQL, LogLevelInfo, "script after translation: %s\n", script)

		if !fake {
			if scriptType == migration.TypeGo {
				err = _migration.Up(tx)
			} else {
				_, err = tx.Exec(script)
			}
		}
	}

	if err == nil {
		if fake {
			c.Log(LoggerGorpheus, LogLevelInfo, "[ FAKED ]\n")
		} else {
			c.Log(LoggerGorpheus, LogLevelInfo, "[ OK ]\n")
		}

		if direction == DirectionUp {
			err = c.InsertRevision(tx, _migration.GetVersion())
		} else {
			err = c.RemoveRevision(tx, _migration.GetVersion())
		}
	} else {
		c.Log(LoggerGorpheus, LogLevelInfo, "[ ERROR ]\n")
	}
	return err
}
