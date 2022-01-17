package gorpheus

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/toudi/gorpheus/v1/migration"
)

func (c *Collection) Migrate(params *MigrationParams) error {
	var operationDesc = map[uint8]string{
		DirectionUp:   "applying",
		DirectionDown: "unapplying",
	}

	var err error
	var revision string
	var db *sqlx.DB

	sort.Sort(c.Versions)

	// index of currently applied migration within revisions table
	var currentIdx int = -1
	// index of target migration (i.e. the one that user wants to migrate to) within revisions table
	var targetIdx int
	// how many revisions should we apply - that's just for the internal forloop
	var numRevisionsToApply int
	var directionMode = DirectionUp
	// how much should the index change after each loop. that's +1 for migrating up, -1 for rolling back
	var delta = 1

	if db, err = c.connectToDb(params); err != nil {
		log.Fatalf("unable to connect to the database: %v", err)
	}

	fmt.Printf("making sure that migrations table exist\n")
	if err = c.ensureMigrationsTableExists(db); err != nil {
		log.Fatalf("cannot create migrations table: %v", err)
	}

	// let's start by looking up the indices and based on that detect whether we are migrating or
	// rolling back.

	if revision, err = c.retrieveCurrentRevision(db); err != nil {
		fmt.Printf("cannot detect current revision: %v", err)
		// we don't really have to do anything since go assigns zero by default
	}
	if revision != "" {
		if currentIdx, err = c.FindMigration(revision); err != nil {
			log.Fatalf("cannot calculate index of current revision in revisions table: %v", err)
		}
	}
	if targetIdx, err = c.FindMigrationWithNamespaceAndRevision(params.Namespace, params.Revision); err != nil {
		log.Fatalf("cannot find target revision: %v", err)
	}
	if params.Zero {
		targetIdx = -1
	}

	if currentIdx == targetIdx {
		fmt.Printf("no migrations to apply.\n")
		return nil
	}

	if targetIdx < currentIdx {
		directionMode = DirectionDown
		targetIdx += 1
		delta = -1
	} else {
		currentIdx += 1
	}

	numRevisionsToApply = (currentIdx - targetIdx)
	if numRevisionsToApply < 0 {
		numRevisionsToApply = -numRevisionsToApply
	}
	numRevisionsToApply += 1

	fmt.Printf("currentIdx=%d, targetIdx=%d\n", currentIdx, targetIdx)
	fmt.Printf("%s %d migration(s)\n", operationDesc[uint8(directionMode)], numRevisionsToApply)

	for i := 0; i < numRevisionsToApply; i++ {
		err = nil
		fmt.Printf("%s %s .. ", operationDesc[uint8(directionMode)], c.Versions[currentIdx].Revision())

		if err = Atomic(db, func(tx *sqlx.Tx) error {
			return c.performMigration(tx, currentIdx, directionMode)
		}); err != nil {
			log.Fatalf("unable to perform operation: %v", err)
			break
		}
		fmt.Printf("[OK]\n")

		currentIdx += delta
	}

	if params.Zero {
		err = c.dropMigrationsTable(db)
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

func (c *Collection) performMigration(tx *sqlx.Tx, idx int, direction int) error {
	var script string
	var scriptType uint8
	var err error

	_migration := c.Versions[idx]

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
	if scriptType == migration.TypeGo {
		err = _migration.Up(tx)
	} else {
		_, err = tx.Exec(script)
	}
	if err == nil {
		if direction == DirectionUp {
			return c.InsertRevision(tx, _migration.Revision())
		}
		return c.RemoveRevision(tx, _migration.Revision())
	}
	return err
}
