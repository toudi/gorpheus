package planner

import (
	"github.com/toudi/gorpheus/v2/interfaces"
	"github.com/toudi/gorpheus/v2/internal/migration"
)

type Planner struct {
	migrations     []*migration.Migration // source migrations array
	currentVersion map[string]int         // what version (if any) is currently applied per namespace
	lookupIndex    map[string]int         // position of a *revision* inside `migrations` slice
}

// provide a default constructor and optionally an array of initializer functions
func New(initializers ...func(p *Planner)) *Planner {
	p := &Planner{
		currentVersion: make(map[string]int),
		lookupIndex:    make(map[string]int),
	}

	for _, init := range initializers {
		init(p)
	}

	return p
}

func WithMigrations(migrations []*migration.Migration) func(p *Planner) {
	return func(p *Planner) {
		p.migrations = sortMigrations(migrations)
		// prepare the lookup index.
		for idx, _migration := range p.migrations {
			p.lookupIndex[migration.LookupID(_migration.Revision.Namespace, _migration.Revision.Version)] = idx
		}
	}
}

func WithApplied(applied []interfaces.MigrationsHistoryRow) func(p *Planner) {
	return func(p *Planner) {
		for _, revision := range applied {
			// let's parse the number from the migration id and set it as the current version for that namespace
			version, _ := migration.NumericVersion(revision.Migration)
			// `applied` array is guaranteed to be sorted by the timestamp in ascending order
			// so we can simply override the current version for that namespace.
			p.currentVersion[revision.Namespace] = version
		}
	}
}

// equivalent of SetApplied, but which takes a slice of strings as input
func WithAppliedFromStringSlice(applied []string) func(p *Planner) {
	return func(p *Planner) {
		for _, revision := range applied {
			namespace, _, version, _ := migration.ParseId(revision)
			p.currentVersion[namespace] = version
		}
	}
}
