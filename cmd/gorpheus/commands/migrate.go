package commands

type Migrate struct {
	migrationsDiscovery
	Target []string `name:"dest" arg:"" optional:""`
}

func (m *Migrate) Run(ctx *Context) error {
	if err := m.migrationsDiscovery.handleMigrationSources(ctx); err != nil {
		return err
	}

	var destNamespace string = ""
	var destRevision string = ""

	if len(m.Target) > 0 {
		destNamespace = m.Target[0]
		if len(m.Target) > 1 {
			destRevision = m.Target[1]
		}
	}

	return ctx.Gorpheus.Migrate(destNamespace, destRevision)
}
