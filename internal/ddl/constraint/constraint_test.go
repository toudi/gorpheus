package constraint_test

import (
	"testing"

	"github.com/toudi/gorpheus/v2/internal/ddl/constraint"

	"github.com/stretchr/testify/require"
)

func TestConstraint(t *testing.T) {
	t.Run("test Unique", func(t *testing.T) {
		t.Run("unique (single column)", func(t *testing.T) {
			t.Parallel()
			c := &constraint.Constraint{
				SrcUnique: "foo",
			}
			err := c.Validate()
			require.NoError(t, err)
			require.Equal(t, &constraint.Unique{
				Columns: []string{"foo"},
			}, c.Unique)
		})
		t.Run("unique (multiple columns)", func(t *testing.T) {
			t.Parallel()
			c := &constraint.Constraint{
				SrcUnique: []string{"foo", "bar"},
			}
			err := c.Validate()
			require.NoError(t, err)
			require.Equal(t, &constraint.Unique{
				Columns: []string{"foo", "bar"},
			}, c.Unique)
		})
		t.Run("unique (map)", func(t *testing.T) {
			t.Parallel()
			c := &constraint.Constraint{
				SrcUnique: map[string]any{
					"columns": []string{"foo"},
				},
			}
			err := c.Validate()
			require.NoError(t, err)
			require.Equal(t, &constraint.Unique{
				Columns: []string{"foo"},
			}, c.Unique)
		})
		t.Run("unique (invalid)", func(t *testing.T) {
			t.Parallel()
			c := &constraint.Constraint{
				SrcUnique: 42,
			}
			err := c.Validate()
			require.ErrorIs(t, err, constraint.ErrUnableToParseUnique)
		})
	})

	t.Run("test foreign key", func(t *testing.T) {
		t.Run("automatically add column to slice of columns", func(t *testing.T) {
			t.Parallel()

			c := &constraint.Constraint{
				ForeignKey: &constraint.ForeignKey{
					Column: "foo",
				},
			}
			err := c.Validate()
			require.NoError(t, err)
			require.Equal(t, []string{"foo"}, c.ForeignKey.Columns)
		})
	})
}
