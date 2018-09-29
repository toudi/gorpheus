package gorpheus_test

import (
	"reflect"
	"testing"

	"github.com/toudi/gorpheus"
)

func migration_factory(version string, previous string, mtype uint8) gorpheus.Migration {
	return gorpheus.Migration{
		Version:         version,
		PreviousVersion: previous,
		Type:            mtype,
		UpFn:            nil,
		DownFn:          nil,
	}
}

func Test_RegisterMigration(t *testing.T) {
	gorpheus.Reset()
	gorpheus.Register(
		migration_factory("0001_initial", "", gorpheus.TypeSQL),
	)

	gorpheus.Register(
		migration_factory("0002_some_changes", "", gorpheus.TypeFizz),
	)

	if len(gorpheus.Migrations) != 2 {
		t.Error("Length of migrations should be equal to 2")
	}
}

func Test_SortMigrations(t *testing.T) {
	gorpheus.Reset()
	gorpheus.Register(
		migration_factory("0001_initial", "", gorpheus.TypeSQL),
	)
	gorpheus.Register(
		migration_factory("0002_some_changes", "", gorpheus.TypeFizz),
	)
	gorpheus.Register(
		migration_factory("0003_foo", "", gorpheus.TypeSQL),
	)
	gorpheus.Register(
		migration_factory("0003_bar", "0002_some_changes", gorpheus.TypeFizz),
	)
	gorpheus.Register(
		migration_factory("0004_foo", "0003_foo", gorpheus.TypeSQL),
	)
	gorpheus.Register(
		migration_factory("0004_bar", "", gorpheus.TypeFizz),
	)
	gorpheus.Register(
		migration_factory("0005_bar", "", gorpheus.TypeFizz),
	)
	gorpheus.Register(
		migration_factory("0006_bar", "", gorpheus.TypeFizz),
	)
	gorpheus.SortMigrations()
	var expected_names_sorted = []string{
		"0001_initial", "0002_some_changes", "0003_bar", "0003_foo", "0004_foo",
		"0004_bar", "0005_bar", "0006_bar",
	}
	if !reflect.DeepEqual(gorpheus.MigrationNames(), expected_names_sorted) {
		t.Errorf("Sorting failed: expected: %v, value: %v\n", expected_names_sorted, gorpheus.MigrationNames())
	}

}
