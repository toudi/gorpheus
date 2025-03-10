package planner_test

import (
	"testing"

	"github.com/toudi/gorpheus/v2/internal/migration"
	"github.com/toudi/gorpheus/v2/internal/planner"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestPathCalculation(t *testing.T) {
	// regardless of the performed test, we will always use the same
	// array of migrations.
	var usersMigrations = []*migration.Migration{
		{
			Revision: migration.Revision{
				Namespace: "users",
				Name:      "0001_initial",
				Version:   1,
			},
		},
		{
			Revision: migration.Revision{
				Namespace: "users",
				Name:      "0002_add_email",
				Version:   2,
			},
		},
		{
			Revision: migration.Revision{
				Namespace: "users",
				Name:      "0003_add_password",
				Version:   3,
			},
		},
	}

	var booksMigrations = []*migration.Migration{
		{
			Revision: migration.Revision{
				Namespace: "books",
				Name:      "0001_initial",
				Version:   1,
			},
		},
		{
			Revision: migration.Revision{
				Namespace: "books",
				Name:      "0002_add_author",
				Version:   2,
			},
			Dependencies: []migration.Revision{
				{
					Namespace: "users", Name: "0003_add_password", Version: 3,
				},
			},
		},
	}

	var allMigrations []*migration.Migration
	allMigrations = append(allMigrations, booksMigrations...)
	allMigrations = append(allMigrations, usersMigrations...)

	// let's also define a helper method to generate a list of migration id's for
	// easier verification.
	var migrationIds = func(input []*migration.Migration) []string {
		return lo.Map(input, func(m *migration.Migration, _ int) string {
			return m.Revision.String()
		})
	}

	t.Run("calculate path forwards", func(t *testing.T) {
		t.Run("empty set", func(t *testing.T) {
			t.Parallel()

			planner := planner.New()
			plan, _, err := planner.PreparePlan("", "")
			require.NoError(t, err)
			require.Nil(t, plan)
		})

		t.Run("no applied migrations", func(t *testing.T) {
			t.Run("apply everything", func(t *testing.T) {
				t.Parallel()

				var expected = []string{
					"books/0001_initial",
					"users/0001_initial",
					"users/0002_add_email",
					"users/0003_add_password",
					// in order to apply books/0002 we have to first apply users/0002
					// however because users/0002 is not applied yet we fist have to
					// apply users/0001
					"books/0002_add_author",
				}

				planner := planner.New(planner.WithMigrations(allMigrations))

				plan, _, err := planner.PreparePlan("", "")
				require.NoError(t, err)
				require.Equal(t, expected, migrationIds(plan))
			})

			t.Run("only a single namespace", func(t *testing.T) {
				t.Parallel()

				var expected = []string{
					"users/0001_initial",
					"users/0002_add_email",
					"users/0003_add_password",
				}

				planner := planner.New(planner.WithMigrations(allMigrations))

				plan, _, err := planner.PreparePlan("users", "")
				require.NoError(t, err)
				require.Equal(t, expected, migrationIds(plan))
			})

			t.Run("a single namespace, with dependencies", func(t *testing.T) {
				t.Parallel()

				var expected = []string{
					"books/0001_initial",
					"users/0001_initial",
					"users/0002_add_email",
					"users/0003_add_password",
					// in order to apply books/0002 we have to first apply users/0002
					// however because users/0002 is not applied yet we fist have to
					// apply users/0001
					"books/0002_add_author",
				}

				planner := planner.New(planner.WithMigrations(allMigrations))

				plan, _, err := planner.PreparePlan("books", "")
				require.NoError(t, err)
				require.Equal(t, expected, migrationIds(plan))
			})
		})

		t.Run("applied migrations", func(t *testing.T) {
			t.Run("no dependencies", func(t *testing.T) {
				t.Parallel()

				var expected = []string{
					"users/0003_add_password",
				}

				planner := planner.New(
					planner.WithMigrations(allMigrations),
					planner.WithAppliedFromStringSlice(
						[]string{
							"users/0001_initial",
							"users/0002_add_email",
						},
					),
				)

				plan, _, err := planner.PreparePlan("users", "")
				require.NoError(t, err)
				require.Equal(t, expected, migrationIds(plan))
			})

			t.Run("everything done", func(t *testing.T) {
				t.Parallel()

				var expected = []string{}

				planner := planner.New(
					planner.WithMigrations(allMigrations),
					planner.WithAppliedFromStringSlice(
						[]string{
							"users/0001_initial",
							"users/0002_add_email",
							"users/0003_add_password",
						},
					),
				)

				plan, _, err := planner.PreparePlan("users", "")
				require.NoError(t, err)
				require.Equal(t, expected, migrationIds(plan))
			})

			t.Run("everything done (specific version)", func(t *testing.T) {
				t.Parallel()

				var expected = []string{}

				planner := planner.New(
					planner.WithMigrations(allMigrations),
					planner.WithAppliedFromStringSlice(
						[]string{
							"users/0001_initial",
							"users/0002_add_email",
						},
					),
				)

				plan, _, err := planner.PreparePlan("users", "2")
				require.NoError(t, err)
				require.Equal(t, expected, migrationIds(plan))
			})

			t.Run("dependencies", func(t *testing.T) {
				t.Parallel()

				var expected = []string{
					"books/0001_initial",
					"users/0002_add_email",
					"users/0003_add_password",
					"books/0002_add_author",
				}

				planner := planner.New(
					planner.WithMigrations(allMigrations),
					planner.WithAppliedFromStringSlice(
						[]string{
							"users/0001_initial",
						},
					),
				)

				plan, _, err := planner.PreparePlan("books", "")
				require.NoError(t, err)
				require.Equal(t, expected, migrationIds(plan))
			})

			t.Run("dependencies should not be touched", func(t *testing.T) {
				t.Parallel()

				var expected = []string{
					"books/0001_initial",
					"books/0002_add_author",
				}

				planner := planner.New(
					planner.WithMigrations(allMigrations),
					planner.WithAppliedFromStringSlice(
						[]string{
							"users/0001_initial",
							"users/0002_add_email",
							"users/0003_add_password",
						},
					),
				)

				plan, _, err := planner.PreparePlan("books", "")
				require.NoError(t, err)
				require.Equal(t, expected, migrationIds(plan))
			})
			t.Run("detect invalid version (positive)", func(t *testing.T) {
				t.Parallel()

				var expected = []string{}

				_planner := planner.New(
					planner.WithMigrations(allMigrations),
					planner.WithAppliedFromStringSlice([]string{"users/0001_initial"}),
				)

				plan, direction, err := _planner.PreparePlan("users", "10")
				require.Error(t, planner.ErrInvalidVersion, err)
				require.Equal(t, planner.DIRECTION_FORWARD, direction)
				require.Equal(t, expected, migrationIds(plan))

			})

		})
	})

	t.Run("calculate path backwards", func(t *testing.T) {
		t.Run("parse 0", func(t *testing.T) {
			t.Parallel()

			var expected = []string{
				"users/0001_initial",
			}

			_planner := planner.New(
				planner.WithMigrations(allMigrations),
				planner.WithAppliedFromStringSlice([]string{"users/0001_initial"}),
			)

			plan, direction, err := _planner.PreparePlan("users", "0")
			require.NoError(t, err)
			require.Equal(t, planner.DIRECTION_BACKWARD, direction)
			require.Equal(t, expected, migrationIds(plan))

		})
		t.Run("parse zero", func(t *testing.T) {
			t.Parallel()

			var expected = []string{
				"users/0001_initial",
			}

			_planner := planner.New(
				planner.WithMigrations(allMigrations),
				planner.WithAppliedFromStringSlice([]string{"users/0001_initial"}),
			)

			plan, direction, err := _planner.PreparePlan("users", "zero")
			require.NoError(t, err)
			require.Equal(t, planner.DIRECTION_BACKWARD, direction)
			require.Equal(t, expected, migrationIds(plan))

		})
		t.Run("detect invalid version (negative)", func(t *testing.T) {
			t.Parallel()

			var expected = []string{}

			_planner := planner.New(
				planner.WithMigrations(allMigrations),
				planner.WithAppliedFromStringSlice([]string{"users/0001_initial"}),
			)

			plan, direction, err := _planner.PreparePlan("users", "-1")
			require.Error(t, planner.ErrInvalidVersion, err)
			require.Equal(t, planner.DIRECTION_FORWARD, direction)
			require.Equal(t, expected, migrationIds(plan))

		})

		t.Run("no dependencies", func(t *testing.T) {
			t.Parallel()

			var expected = []string{
				"users/0002_add_email",
				"users/0001_initial",
			}

			_planner := planner.New(
				planner.WithMigrations(allMigrations),
				planner.WithAppliedFromStringSlice([]string{"users/0001_initial", "users/0002_add_email"}),
			)

			plan, direction, err := _planner.PreparePlan("users", "0")
			require.NoError(t, err)
			require.Equal(t, planner.DIRECTION_BACKWARD, direction)
			require.Equal(t, expected, migrationIds(plan))
		})
		t.Run("no dependencies (books)", func(t *testing.T) {
			t.Parallel()

			var expected = []string{
				"books/0002_add_author",
			}

			_planner := planner.New(
				planner.WithMigrations(allMigrations),
				planner.WithAppliedFromStringSlice(
					[]string{
						"users/0001_initial",
						"users/0002_add_email",
						"books/0001_initial",
						"books/0002_add_author",
					},
				),
			)

			plan, direction, err := _planner.PreparePlan("books", "1")
			require.NoError(t, err)
			require.Equal(t, planner.DIRECTION_BACKWARD, direction)
			require.Equal(t, expected, migrationIds(plan))
		})
		t.Run("dependencies (books)", func(t *testing.T) {
			t.Parallel()

			// so basically books/0001_initial has a dependency of users/0003
			// however this is not a conflict - we are not touching users namespace
			// in this test and so it is safe to unapply just books namespace migrations.
			var expected = []string{
				"books/0002_add_author",
				"books/0001_initial",
			}

			_planner := planner.New(
				planner.WithMigrations(allMigrations),
				planner.WithAppliedFromStringSlice(
					[]string{
						"users/0001_initial",
						"users/0002_add_email",
						"users/0003_add_password",
						"books/0001_initial",
						"books/0002_add_author",
					},
				),
			)

			plan, direction, err := _planner.PreparePlan("books", "0")
			require.NoError(t, err)
			require.Equal(t, planner.DIRECTION_BACKWARD, direction)
			require.Equal(t, expected, migrationIds(plan))
		})

		t.Run("unapply migration that is a dependency", func(t *testing.T) {
			t.Parallel()

			var expected = []string{
				"books/0002_add_author",
				"users/0003_add_password",
			}

			_planner := planner.New(
				planner.WithMigrations(allMigrations),
				planner.WithAppliedFromStringSlice(
					[]string{
						"users/0001_initial",
						"users/0002_add_email",
						"users/0003_add_password",
						"books/0001_initial",
						"books/0002_add_author",
					},
				),
			)

			plan, direction, err := _planner.PreparePlan("users", "2")
			require.NoError(t, err)
			require.Equal(t, planner.DIRECTION_BACKWARD, direction)
			require.Equal(t, expected, migrationIds(plan))
		})

		t.Run("unapply migration that is a dependency but dependency itself is not applied", func(t *testing.T) {
			t.Parallel()

			var expected = []string{
				"users/0003_add_password",
			}

			_planner := planner.New(
				planner.WithMigrations(allMigrations),
				planner.WithAppliedFromStringSlice(
					[]string{
						"users/0001_initial",
						"users/0002_add_email",
						"users/0003_add_password",
					},
				),
			)

			plan, direction, err := _planner.PreparePlan("users", "2")
			require.NoError(t, err)
			require.Equal(t, planner.DIRECTION_BACKWARD, direction)
			require.Equal(t, expected, migrationIds(plan))
		})

	})
}
