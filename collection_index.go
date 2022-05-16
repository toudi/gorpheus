package gorpheus

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/toudi/gorpheus/migration"
)

type NamespaceMeta struct {
	current       int
	mostRecent    int
	positionIndex map[int]int
}

const AppliedRevisionsQuery = "SELECT revision FROM gorpheus_revisions"

// index function creates indices of migrations so that when we want to
// apply them we will know at which position in the array they residue.
func (c *Collection) index(db *sqlx.DB) error {
	var err error
	var exists bool
	var namespace string
	var versionNo int
	var namespaceMeta NamespaceMeta

	var appliedRevision string

	result, err := db.Query(AppliedRevisionsQuery)
	for result.Next() {
		if err = result.Scan(&appliedRevision); err != nil {
			return fmt.Errorf("cannot parse applied revision: %v", err)
		}
		c.applied[appliedRevision] = true
	}

	for collectionIndex, migration := range c.Versions {
		if namespace, err = migration.GetNamespace(); err != nil {
			return fmt.Errorf("could not extract namespace for %v: %v", migration.GetVersion(), err)
		}
		if versionNo, err = migration.GetVersionNumber(); err != nil {
			return fmt.Errorf("could not get version number for %v: %v", migration.GetVersion(), err)
		}
		if namespaceMeta, exists = c.metadata[namespace]; !exists {
			namespaceMeta = NamespaceMeta{positionIndex: make(map[int]int), current: -1}
		}

		namespaceMeta.positionIndex[versionNo] = collectionIndex
		if versionNo > namespaceMeta.mostRecent {
			namespaceMeta.mostRecent = versionNo
		}
		if _, exists = c.applied[migration.GetVersion()]; exists && versionNo > namespaceMeta.current {
			namespaceMeta.current = versionNo
		}

		c.metadata[namespace] = namespaceMeta

	}

	return err
}

func (c *Collection) GetDependencies(_migration migration.MigrationI) ([]migration.MigrationI, error) {
	var versionNo int
	var err error
	var namespace string

	// fmt.Printf("getDependencies of %s\n", _migration.GetVersion())
	migrationDependencies := _migration.Dependencies()
	dependencies := make([]migration.MigrationI, 0)
	// fmt.Printf("dependencies: %+v\n", dependencies)

	for _, dependency := range migrationDependencies {
		namespace, versionNo, err = parseNamespaceAndVersionNo(dependency)
		// fmt.Printf("result of parseNamespaceAndVersionNo: %s, %d, %v\n", namespace, versionNo, err)
		if err != nil {
			return nil, fmt.Errorf("could not detect dependency: %v", err)
		}

		dependencies = append(dependencies, c.Versions[c.metadata[namespace].positionIndex[versionNo]])
		// fmt.Printf("dependecies array: %+v\n", dependencies)
	}

	versionNo, err = _migration.GetVersionNumber()

	if err != nil {
		return nil, fmt.Errorf("could not detect version number of current migration: %v", err)
	}

	if versionNo > 1 {
		namespace, err = _migration.GetNamespace()
		if err != nil {
			return nil, fmt.Errorf("could not parse namespace from %s: %v", _migration.GetVersion(), err)
		}
		dependencies = append(dependencies, c.Versions[c.metadata[namespace].positionIndex[versionNo-1]])
	}

	return dependencies, nil
}

func parseNamespaceAndVersionNo(revision string) (string, int, error) {
	underscoreIdx := strings.Index(revision, "_")
	slashIdx := strings.Index(revision, "/")

	namespace := revision[:slashIdx]
	versionNumberString := revision[slashIdx+1 : underscoreIdx]
	versionNumber, err := strconv.Atoi(versionNumberString)
	if err != nil {
		return "", -1, fmt.Errorf("unable to parse the numeric revision: %v", err)
	}

	return namespace, versionNumber, nil
}
