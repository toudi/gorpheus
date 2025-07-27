package planner

import (
	"slices"

	"github.com/toudi/gorpheus/v2/internal/migration"
)

type DFSSorter struct {
	sortedMigrations []*migration.Migration // used when recursing with DFS to avoid nasty syntax of *foo = append(*foo, bar)
	alreadyVisited   map[string]bool
	migrationsMap    map[string]*migration.Migration
}

func (s *DFSSorter) visitMigrationNode(migration *migration.Migration) {
	for _, dependency := range migration.Dependencies {
		if !s.alreadyVisited[dependency.String()] {
			s.visitMigrationNode(s.migrationsMap[dependency.String()])
		}
	}
	if !s.alreadyVisited[migration.Revision.String()] {
		s.sortedMigrations = append(s.sortedMigrations, migration)
		s.alreadyVisited[migration.Revision.String()] = true
	}
}

func sortMigrations(migrations []*migration.Migration) []*migration.Migration {
	sorter := &DFSSorter{
		// prepare a lookup map (we will need it during the DFS traversal)

		migrationsMap:  make(map[string]*migration.Migration),
		alreadyVisited: make(map[string]bool),
	}

	for currIdx, migration := range migrations {
		// side effect of using DFS is that we now have to pre-emptively add dependencies between migrations of the same
		// namespace. Otherwise we could end up with a slightly odd order.
		// Oh well :D
		if currIdx > 0 {
			previous := migrations[currIdx-1]
			if migration.Revision.Namespace == previous.Revision.Namespace && !slices.Contains(migration.RawDependencies, previous.Revision.String()) {
				migration.Dependencies = append(migration.Dependencies, previous.Revision)
			}
		}
		sorter.migrationsMap[migration.Revision.String()] = migration
	}

	// for a long while, I was using slices.StableSortFunc and referring to dependencies slice. However,
	// it turns out that some of the nodes were not always comapared against. I've googled for this and
	// it turns out that quicksort is not recommented way of dealing with dependent objects. Oh well.
	for _, migration := range migrations {
		sorter.visitMigrationNode(migration)
	}

	return sorter.sortedMigrations
}
