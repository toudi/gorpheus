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

	if len(p.migrations) == 0 {
		return nil, direction, nil
	}

	var start, dest *migration.Migration
	// let's initialize the defaults:
	start = p.migrations[0]
	dest = p.migrations[len(p.migrations)-1]

	if namespace != "" {
		// if the namespace is not empty then we want to migrate only a single
		// namespace.
		start, _ = p.findMigration(namespace, 1)

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
				// in order to help the algorithm, we're creating a fake migration so that
				// to know when to stop (re)cursing.
				dest = &migration.Migration{
					Revision: migration.Revision{
						Namespace: namespace,
						Version:   0,
					},
				}
				direction = DIRECTION_BACKWARD
				versionNum = 0
			} else {
				versionNum, err = strconv.Atoi(version)
				if err != nil {
					return nil, direction, err
				}
				var found bool
				if dest, found = p.findMigration(namespace, versionNum); !found {
					return nil, direction, ErrInvalidVersion
				}
			}
			if versionNum < currentVersion {
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
	if start == dest {
		// we've reached the end. However, it could be that this is the last node to be applied
		// so let's check against that:
		if !p.isApplied(start) {
			plan.add(start)
		}
		return nil
	}
	var finished = false
	for !finished {
		if direction == DIRECTION_FORWARD {
			for _, d := range start.Dependencies {
				if plan.visited[d] {
					continue
				}
				// so this basically means - if we've stumbled upon a dependency let's make
				// sure that all migrations that lead up to this dependency are applied as well.
				startOfDependency, _ := p.findMigration(d.Namespace, 1)
				dependency, _ := p.findMigration(d.Namespace, d.Version)
				if err := p.pathBetween(plan, startOfDependency, dependency, direction); err != nil {
					return err
				}
			}
			if !p.isApplied(start) {
				plan.add(start)
			}
			// can we find next migration ?
			var found bool
			start, found = p.findNextMigration(start.Revision)
			if !found {
				// if not then it means we've reached the end of this namespace.
				return nil
			}
			return p.pathBetween(plan, start, dest, direction)
		} else {
			// let's find all migrations that depend on the one we're trying to unapply.
			dependencies := p.findDependenciesOf(start.Revision)
			for _, dependency := range dependencies {
				if plan.visited[dependency.Revision] {
					continue
				}
				// now we need to migrate these dependencies to one version less.
				dependencyCurrent, found := p.findMigration(dependency.Revision.Namespace, dependency.Revision.Version)
				if !found {
					// if it wasn't found then all the better - we've got less work to do.
					// this simply means that the dependency was not (yet) applied.
					continue
				}
				// otherwise, let's go on:
				if dependency.Revision.Version > 1 {
					dependencyMigration, _ := p.findMigration(dependency.Revision.Namespace, dependency.Revision.Version-1)
					if err := p.pathBetween(plan, dependencyCurrent, dependencyMigration, direction); err != nil {
						return err
					}
				} else {
					if p.isApplied(dependency) {
						plan.add(dependency)
					}
				}
			}
			if p.isApplied(start) {
				plan.add(start)
			}
			if start.Revision.Namespace == dest.Revision.Namespace && start.Revision.Version == dest.Revision.Version+1 {
				finished = true
			} else {
				// if not then we need to calculate next index. which is the one from previous version from
				// within the same namespace.
				start, _ = p.findMigration(start.Revision.Namespace, start.Revision.Version-1)
			}
		}
	}

	return nil
}
