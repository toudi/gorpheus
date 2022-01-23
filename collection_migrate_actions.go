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

func (c *Collection) Migrate(params *MigrationParams) error {
	var description string

	var err error
	var db *sqlx.DB
	var breakpoint string

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

	migrationsToApply := make([]migration.MigrationI, 0)
	for namespace, metadata := range c.metadata {
		if params.Namespace != "" && namespace != params.Namespace {
			continue
		}
		migrationsToApply = append(migrationsToApply, c.Versions[metadata.positionIndex[metadata.mostRecent]])
	}

	if params.Namespace != "" {
		namespaceMeta, exists := c.metadata[params.Namespace]
		if !exists {
			return fmt.Errorf("unknown namespace: %s", params.Namespace)
		}
		if params.VersionNo == -1 {
			params.VersionNo = namespaceMeta.mostRecent - 1
		} else if params.VersionNo == 0 {
			if !params.Zero {
				params.VersionNo = namespaceMeta.mostRecent
			} else {
				breakpoint = c.Versions[namespaceMeta.positionIndex[1]].Revision()
			}
		}
		var migrationIndex int

		if breakpoint == "" {
			migrationIndex, exists = namespaceMeta.positionIndex[params.VersionNo]
			if !exists {
				return fmt.Errorf("unknown version number: %d", params.VersionNo)
			}
			breakpoint = c.Versions[migrationIndex].Revision()
		}
		fmt.Printf("current = %v\n", namespaceMeta.current)
		if params.VersionNo < c.metadata[params.Namespace].mostRecent && namespaceMeta.current != -1 {
			direction = DirectionDown
			if !params.Zero {
				breakpoint = c.Versions[migrationIndex+1].Revision()
			}
		}
	}

	description = operationDesc[uint8(direction)]

	err = Atomic(db, func(tx *sqlx.Tx) error {
		for _, _migration := range migrationsToApply {
			if _, err = c.performMigration(tx, _migration, direction, breakpoint, 0); err != nil {
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

func (c *Collection) performMigration(tx *sqlx.Tx, _migration migration.MigrationI, direction int, breakpoint string, indent int) (bool, error) {
	var script string
	var scriptType uint8
	var err error
	var breakpointReached bool
	var performAction bool = true

	fmt.Printf("%s performMigration(%s, %d, %s)\n", strings.Repeat(" ", indent), _migration.Revision(), direction, breakpoint)

	dependenciesArray, err := c.GetDependencies(_migration)
	// fmt.Printf("result of getDependencies(%s): %v, %v\n", _migration.Revision(), dependenciesArray, err)
	if err != nil {
		return false, fmt.Errorf("could not parse dependencies of %s: %v", _migration.Revision(), err)
	}
	// if the direction is upwards then we apply the dependencies prior to applying the actual
	// migration
	if direction == DirectionUp {
		for _, dependency := range dependenciesArray {
			if _, exists := c.applied[dependency.Revision()]; !exists {
				// fmt.Printf("call c.performMigration(%+v)\n", dependency)
				if breakpointReached, err = c.performMigration(tx, dependency, direction, breakpoint, indent+1); breakpointReached || err != nil {
					fmt.Printf("%s returning since breakpointReached=%v; err=%v\n", strings.Repeat(" ", indent), breakpointReached, err)
					return breakpointReached, err
				}
			}
		}
	}

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
		fmt.Printf("%s %s %s\n", strings.Repeat(" ", indent), operationDesc[uint8(direction)], _migration.Revision())

		if direction == DirectionUp {
			script, scriptType, err = _migration.UpScript()
		} else {
			script, scriptType, err = _migration.DownScript()
		}

		if err != nil {
			return false, fmt.Errorf("unable to obtain migration script: %v", err)
		}

		if scriptType == migration.TypeFizz {
			if script, err = c.TranslatedSQL(script); err != nil {
				return false, fmt.Errorf("cannot translate migration script: %v", err)
			}
		}
		if scriptType == migration.TypeGo {
			err = _migration.Up(tx)
		} else {
			_, err = tx.Exec(script)
		}
	}

	if err == nil {
		breakpointReached = _migration.Revision() == breakpoint
		fmt.Printf("%s breakpointReached=%v\n", strings.Repeat(" ", indent), breakpointReached)
		if direction == DirectionUp {
			err = c.InsertRevision(tx, _migration.Revision())
		} else {
			err = c.RemoveRevision(tx, _migration.Revision())
			if err == nil {
				if !breakpointReached {
					// https://stackoverflow.com/questions/28058278/how-do-i-reverse-a-slice-in-go
					for i, j := 0, len(dependenciesArray)-1; i < j; i, j = i+1, j-1 {
						dependenciesArray[i], dependenciesArray[j] = dependenciesArray[j], dependenciesArray[i]
					}
					fmt.Printf("%s dependencies for %s => %v\n", strings.Repeat(" ", indent), _migration.Revision(), dependenciesArray)
					for _, dependency := range dependenciesArray {
						fmt.Printf("%s call c.performMigration(%+v)\n", strings.Repeat(" ", indent), dependency)
						if breakpointReached, err = c.performMigration(tx, dependency, direction, breakpoint, indent+1); breakpointReached || err != nil {
							fmt.Printf("%s returning since breakpointReached=%v; err=%v\n", strings.Repeat(" ", indent), breakpointReached, err)
							return breakpointReached, err
						}
					}
				}
			}
		}
	}
	return breakpointReached, err
}
