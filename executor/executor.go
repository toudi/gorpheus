package executor

import (
	"io"
	"path/filepath"
	"slices"
	"strings"

	migrationPkg "github.com/toudi/gorpheus/v2/internal/migration"

	"github.com/toudi/gorpheus/v2/dialect/generic"
	"github.com/toudi/gorpheus/v2/interfaces"
	"github.com/toudi/gorpheus/v2/internal/migration"
)

type Executor struct {
	dialect    interfaces.Dialect
	history    *generic.MigrationsHistoryManager
	migrations []*migration.Migration
	dbState    *dbState
	logger     interfaces.Logger
}

func New(initializers ...func(e *Executor)) *Executor {
	executor := &Executor{logger: &interfaces.DevNullLogger{}}
	executor.dbState = NewDBState(executor)

	for _, initializer := range initializers {
		initializer(executor)
	}

	return executor
}

func WithLogger(logger interfaces.Logger) func(e *Executor) {
	return func(e *Executor) {
		e.logger = logger
	}
}

// this function return just the migration names, for display purposes
func (e *Executor) Migrations() map[string][]string {
	var output = make(map[string][]string)

	for _, migration := range e.migrations {
		namespace := migration.Revision.Namespace
		output[namespace] = append(output[namespace], migration.Revision.Name)
	}

	// now we need to sort them alphabetically
	for namespace, _ := range output {
		slices.Sort(output[namespace])
	}

	return output
}

func (e *Executor) HistoryManager() *generic.MigrationsHistoryManager {
	return e.history
}

func (e *Executor) AddSourceToNamespace(namespace string, name string, source io.Reader) error {
	migration, err := e.Parse(source)
	if err != nil {
		return err
	}

	migration.Revision.Namespace = namespace
	migration.Revision.Name = strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	migration.Revision.Version, err = migrationPkg.NumericVersion(migration.Revision.Name)
	if err != nil {
		return err
	}
	if migration.RawDependencies != nil {
		for _, rawDependency := range migration.RawDependencies {
			revision, err := migrationPkg.RevisionFromString(rawDependency)
			if err != nil {
				return err
			}
			migration.Dependencies = append(migration.Dependencies, revision)
		}
	}
	if err = e.PrepareFinalOperations(migration); err != nil {
		return err
	}

	e.migrations = append(e.migrations, migration)

	return nil
}
