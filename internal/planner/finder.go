package planner

import (
	"github.com/toudi/gorpheus/v2/internal/migration"
)

func (p *Planner) isApplied(m *migration.Migration) bool {
	if currentVersion, exists := p.currentVersion[m.Revision.Namespace]; !exists {
		// if we cannot find the namespace in our current revisions then the migration
		// cannot possibly be applied
		return false
	} else {
		// otherwise, all we have to check is whether the number is bigger than the current
		// one applied.
		return currentVersion >= m.Revision.Version
	}
}

func (p *Planner) findMigration(namespace string, version int) (*migration.Migration, bool) {
	if idx, exists := p.lookupIndex[migration.LookupID(namespace, version)]; exists {
		return p.migrations[idx], true
	}

	return p.migrations[0], false
}

func (p *Planner) findDependenciesOf(id migration.Revision) []*migration.Migration {
	var dependencies []*migration.Migration

	for _, m := range p.migrations {
		for _, d := range m.Dependencies {
			if d == id {
				dependencies = append(dependencies, m)
				break
			}
		}
	}

	return dependencies
}
