package planner

import (
	"errors"
	"strconv"

	"github.com/toudi/gorpheus/v2/internal/migration"

	"github.com/samber/lo"
)

var (
	ErrInvalidVersion = errors.New("unknown namespace revision")
)

// obligatory: Dijikstra probably hates me.
func (p *Planner) PreparePlan(namespace string, version string) ([]*migration.Migration, int, error) {
	var direction = DIRECTION_FORWARD

	var start, dest *migration.Migration
	// let's initialize the defaults:
	start = p.migrations[0]
	dest = p.migrations[len(p.migrations)-1]

	if namespace != "" {
		// if the namespace is not empty then we want to migrate only a single
		// namespace.
		start, _ = p.findMigration(namespace, 1)

		startWithinNamespace := start
		// let's initialize the default to last migration within the namespace
		dest, _, _ = lo.FindLastIndexOf(p.migrations, func(m *migration.Migration) bool {
			return m.Revision.Namespace == namespace
		})
		if version != "" {
			// we want to migrate towards a specific version but at this point
			// we still don't know if it's forwards or backwards
			currentVersion, currentVersionFound := p.currentVersion[namespace]
			if currentVersionFound {
				start, _ = p.findMigration(namespace, currentVersion)
			}
			// if the current version does not exist then it doesn't matter because
			// we will still migrate towards `dest`
			var versionNum int
			var err error

			// version is specified. we have to defice if this is a
			// forwards or backwards migration.
			if version == "zero" || version == "0" {
				// this is kind of an edge case, in that, the user wants to
				// completely revert migrations.
				dest = startWithinNamespace
				direction = DIRECTION_BACKWARD
				versionNum = 0
			} else {
				versionNum, err = strconv.Atoi(version)
				if err != nil {
					return nil, direction, err
				}
				if currentVersionFound && versionNum < currentVersion {
					// let's do a small trick. since the calculate path function appends
					// dest node to the path, let's increment the "expected" version
					// which will mean that the algorithm will not select `dest` to be
					// unapplied.
					versionNum += 1
				}
				var found bool
				if dest, found = p.findMigration(namespace, versionNum); !found {
					return nil, direction, ErrInvalidVersion
				}
			}
			if versionNum <= currentVersion {
				direction = DIRECTION_BACKWARD
			}
		}
	}

	var plan = &Plan{
		visited: make(map[migration.Revision]bool),
	}

	err := p.pathBetween(plan, start, dest, direction)

	return plan.migrations, direction, err
}

func (p *Planner) pathBetween(plan *Plan, start, dest *migration.Migration, direction int) (err error) {
	if direction == DIRECTION_FORWARD {
		for _, d := range start.Dependencies {
			// we have to calculate the path between first migration of the dependency
			// and the dependency itself.
			startOfDependency, _ := p.findMigration(d.Namespace, 1)
			dependency, _ := p.findMigration(d.Namespace, d.Version)
			if err := p.pathBetween(plan, startOfDependency, dependency, direction); err != nil {
				return err
			}
		}
		if !p.isApplied(start) {
			plan.add(start)
		}
	} else {
		// let's see if start migration is a dependency to some other migrations.
		for _, d := range p.findDependenciesOf(start.Revision) {
			// now, for each of such dependencies we have to migrate from the current
			// version of namespace to dependency's version
			dependencyNamespace := d.Revision.Namespace
			if dependencyNamespace == start.Revision.Namespace && d.Revision.Version > start.Revision.Version {
				// this means that this migration is a dependency to the next one whithin the same migration.
				// this isn't a problem as we're migrating downwards anyway
				continue
			}
			currentVersionInNamespace := p.currentVersion[dependencyNamespace]
			if currentVersionInNamespace == 0 {
				// that is not a problem - it simply means that the dependency's namespace
				// is not migrated.
				continue
			}
			currentMigrationInNamespace, _ := p.findMigration(dependencyNamespace, currentVersionInNamespace)
			if err := p.pathBetween(
				plan, currentMigrationInNamespace, d, direction,
			); err != nil {
				return err
			}
		}
		if p.isApplied(start) {
			plan.add(start)
		}
	}
	if start.Revision == dest.Revision {
		return nil
	}
	if direction == DIRECTION_FORWARD {
		if !p.isApplied(start) {
			plan.add(start)
		}
		var found bool
		start, found = p.findMigration(start.Revision.Namespace, start.Revision.Version+1)
		if !found {
			return nil
		}
	} else {
		for _, d := range start.Dependencies {
			dependency, _ := p.findMigration(d.Namespace, d.Version)
			if p.isApplied(dependency) {
				plan.add(dependency)
			}
		}

		var found bool
		start, found = p.findMigration(start.Revision.Namespace, start.Revision.Version-1)
		if !found {
			return nil
		}
	}

	return p.pathBetween(plan, start, dest, direction)
}
