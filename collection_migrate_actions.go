package gorpheus

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/toudi/gorpheus/v1/migration"
)

var operationDesc = map[uint8]string{
	DirectionUp:   "applying",
	DirectionDown: "unapplying",
}

type migrationsArray []migration.MigrationI

func (ma migrationsArray) Contains(_migration migration.MigrationI) bool {
	for _, m := range ma {
		if m.Revision() == _migration.Revision() {
			return true
		}
	}
	return false
}

func (c *Collection) prepareMigrationsToApply(namespace string, currentVersionNo int, targetVersionNo int, dst *migrationsArray) error {
	fmt.Printf("prepareMigrationsToApply(%s, %d, %d)\n", namespace, currentVersionNo, targetVersionNo)
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
		versionNumber, err = _migration.VersionNumber()
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
			fmt.Printf("dependencies for %s: [", _migration.Revision())
			for _, dependency := range dependencies {
				fmt.Printf("%s, ", dependency.Revision())
			}
			fmt.Printf("]\n")
			for _, dependency := range dependencies {
				if dst.Contains(dependency) {
					continue
				}
				fmt.Printf("dependencies loop\n")
				dependencyNamespace := dependency.GetNamespace()
				dependencyMeta := c.metadata[dependencyNamespace]
				dependencyVersionNo, err := dependency.VersionNumber()
				if err != nil {
					return fmt.Errorf("could not get dependency version number: %v", err)
				}
				if dependencyNamespace == namespace && dependencyVersionNo == targetVersionNo {
					continue
				}
				c.prepareMigrationsToApply(dependencyNamespace, dependencyMeta.current, dependencyVersionNo, dst)
			}
			fmt.Printf("namespace = %s; currentVersionNo = %d; targetVersionNo = %d\n", namespace, currentVersionNo, targetVersionNo)
			if delta == 1 {
				*dst = append(*dst, _migration)
			}
		}
		if currentVersionNo == targetVersionNo {
			break
		}
		currentVersionNo += delta
		fmt.Printf("dst after append: %+v\n", dst)
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
		log.Fatalf("unable to connect to the database: %v", err)
	}

	fmt.Printf("making sure that migrations table exist\n")
	if err = c.ensureMigrationsTableExists(db); err != nil {
		log.Fatalf("cannot create migrations table: %v", err)
	}

	if params.DropRevisionsTable {
		return c.dropMigrationsTable(db)
	}

	// check which versions are applied in the database and index the existing collection
	c.index(db)

	var direction = DirectionUp
	if params.Namespace != "" {
		namespaceMeta := c.metadata[params.Namespace]
		if params.Zero || params.VersionNo < namespaceMeta.current {
			direction = DirectionDown
		}
	}

	migrationsToApply := make(migrationsArray, 0)

	fmt.Printf("params: %+v\n", params)

	for namespace, metadata := range c.metadata {
		if params.Namespace != "" && namespace != params.Namespace {
			continue
		}
		currentVersionNo = metadata.current
		targetVersionNo = metadata.mostRecent
		if params.Namespace != "" {
			if params.Zero {
				targetVersionNo = -1
			} else if params.VersionNo > 0 {
				targetVersionNo = params.VersionNo
			}
		}
		fmt.Printf("call to prepareMigrationsToApply(%s, %d, %d)\n", namespace, currentVersionNo, targetVersionNo)
		c.prepareMigrationsToApply(namespace, currentVersionNo, targetVersionNo, &migrationsToApply)
	}

	fmt.Printf("migrations to apply: \n")
	for _, m := range migrationsToApply {
		fmt.Printf("-> %s\n", m.Revision())
	}

	description = operationDesc[uint8(direction)]

	err = Atomic(db, func(tx *sqlx.Tx) error {
		for _, _migration := range migrationsToApply {
			if err = c.performMigration(tx, _migration, direction, params.Fake, 0); err != nil {
				return fmt.Errorf("could not %s %s: %v", description[:len(description)-3], _migration.Revision(), err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("could not perform migrations: %v\n", err)
	}

	return err
}

func (c *Collection) FindMigration(version string) (int, error) {
	for idx, m := range c.Versions {
		if strings.HasPrefix(m.Revision(), version) {
			return idx, nil
		}
	}
	return -1, ErrNoSuchVersion
}

func (c *Collection) FindMigrationWithNamespaceAndRevision(namespace string, revision int) (int, error) {
	// if the namespace was not given then simply migrate to the latest available revision
	if namespace == "" {
		return len(c.Versions) - 1, nil
	}
	// this is a mode where we look for the specific revision within namespace
	if revision > 0 {
		for idx, m := range c.Versions {
			revisionName := m.Revision()

			if strings.HasPrefix(revisionName, namespace) {
				// extract just the numeric part
				underscoreIdx := strings.Index(revisionName, "_")
				slashIdx := strings.Index(revisionName, "/")

				versionNumberString := revisionName[slashIdx+1 : underscoreIdx]
				versionNumber, err := strconv.Atoi(versionNumberString)
				if err != nil {
					return -1, fmt.Errorf("unable to parse the numeric revision: %v", err)
				}
				if versionNumber == revision {
					return idx, nil
				}
			}
		}
	}

	return -1, ErrNoSuchVersion
}

func (c *Collection) performMigration(tx *sqlx.Tx, _migration migration.MigrationI, direction int, fake bool, indent int) error {
	var script string
	var scriptType uint8
	var err error
	var performAction bool = true

	fmt.Printf("%s performMigration(%s, %d)\n", strings.Repeat(" ", indent), _migration.Revision(), direction)

	_, exists := c.applied[_migration.Revision()]

	// basically, if we're migrating upwards and the migration was already applied then let's
	// skip it.
	// also, if we're migrating downwards and the migrations was not applied then there's no point
	// in unapplying it.
	if (direction == DirectionUp && exists) || (!exists && direction == DirectionDown) {
		fmt.Printf("%s no-op\n", strings.Repeat(" ", indent))
		performAction = false
	}

	if performAction {
		fmt.Printf("%s %s %s .. ", strings.Repeat(" ", indent), operationDesc[uint8(direction)], _migration.Revision())

		if direction == DirectionUp {
			script, scriptType, err = _migration.UpScript()
		} else {
			script, scriptType, err = _migration.DownScript()
		}

		if err != nil {
			return fmt.Errorf("unable to obtain migration script: %v", err)
		}

		if scriptType == migration.TypeFizz {
			if script, err = c.TranslatedSQL(script); err != nil {
				return fmt.Errorf("cannot translate migration script: %v", err)
			}
		}
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
			fmt.Printf("[ FAKED ]\n")
		} else {
			fmt.Printf("[ OK ]\n")
		}

		if direction == DirectionUp {
			err = c.InsertRevision(tx, _migration.Revision())
		} else {
			err = c.RemoveRevision(tx, _migration.Revision())
		}
	} else {
		fmt.Printf("[ ERROR ]\n")
	}
	return err
}
