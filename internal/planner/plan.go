package planner

import "github.com/toudi/gorpheus/v2/internal/migration"

type Plan struct {
	migrations []*migration.Migration
	visited    map[migration.Revision]bool
}

func (p *Plan) add(m *migration.Migration) {
	// add `m` to the current plan
	if !p.visited[m.Revision] {
		p.migrations = append(p.migrations, m)
		p.visited[m.Revision] = true
	}
}
