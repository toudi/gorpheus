package commands

import "fmt"

type Show struct {
	migrationsDiscovery
	Namespace string `optional:"" arg:""`
}

func (s *Show) Run(ctx *Context) error {
	if err := s.migrationsDiscovery.handleMigrationSources(ctx); err != nil {
		return err
	}

	historyManager := ctx.Gorpheus.HistoryManager()

	for namespace, migrations := range ctx.Gorpheus.Migrations() {
		// we want to limit migrations to only a single namespace and this ain't the one.
		if s.Namespace != "" && namespace != s.Namespace {
			continue
		}

		fmt.Println(namespace)

		for _, migration := range migrations {
			var appliedMark = " "
			if historyManager.IsApplied(namespace, migration) {
				appliedMark = "x"
			}

			fmt.Printf("[%s] %s\n", appliedMark, migration)
		}
	}

	return nil
}
