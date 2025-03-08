package planner

import (
	"testing"

	"github.com/toudi/gorpheus/v2/interfaces"
	"github.com/toudi/gorpheus/v2/internal/migration"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestGraphConstructor(t *testing.T) {
	t.Run("test sorting of source migrations", func(t *testing.T) {
		t.Parallel()

		planner := New(
			WithMigrations([]*migration.Migration{
				{
					Revision: migration.Revision{
						Namespace: "users",
						Name:      "1",
						Version:   1,
					},
				},
				{
					Revision: migration.Revision{
						Namespace: "users",
						Name:      "2",
						Version:   2,
					},
					Dependencies: []migration.Revision{
						{
							Namespace: "users",
							Name:      "1",
							Version:   1,
						},
					},
				},
				{
					Revision: migration.Revision{
						Namespace: "users",
						Name:      "3",
						Version:   3,
					},
					Dependencies: []migration.Revision{
						{
							Namespace: "users",
							Name:      "2",
							Version:   2,
						},
					},
				},
				{
					Revision: migration.Revision{
						Namespace: "books",
						Name:      "1",
						Version:   1,
					},
				},
				{
					Revision: migration.Revision{
						Namespace: "books",
						Name:      "2",
						Version:   2,
					},
				},
				{
					Revision: migration.Revision{
						Namespace: "books",
						Name:      "3",
						Version:   3,
					},
					Dependencies: []migration.Revision{
						{
							Namespace: "users",
							Name:      "3",
							Version:   3,
						},
					},
				},
			}),
		)

		require.Equal(
			t,
			[]string{
				"books/1",
				"books/2",
				"users/1",
				"users/2",
				"users/3",
				"books/3",
			},
			lo.Map(planner.migrations, func(m *migration.Migration, _ int) string {
				return migration.LookupID(m.Revision.Namespace, m.Revision.Version)
			}),
		)
	})

	t.Run("test calculating current version within namespaces", func(t *testing.T) {
		t.Parallel()

		planner := New(
			WithApplied([]interfaces.MigrationsHistoryRow{
				{
					Namespace: "users",
					Migration: "0001_initial",
				},
				{
					Namespace: "users",
					Migration: "0002_foo",
				},
				{
					Namespace: "users",
					Migration: "0004_bar",
				},
			}),
		)

		require.Equal(t, map[string]int{"users": 4}, planner.currentVersion)
	})

	t.Run("test set applied from string slice", func(t *testing.T) {
		t.Parallel()

		planner := New(
			WithAppliedFromStringSlice([]string{"users/0001_initial", "users/0002_foo", "users/0004_bar"}),
		)

		require.Equal(t, map[string]int{"users": 4}, planner.currentVersion)
	})
}
