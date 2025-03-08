package index_test

import (
	"testing"

	"github.com/toudi/gorpheus/v2/internal/ddl/index"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	t.Run("test Parse", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			t.Run("plain index", func(t *testing.T) {
				t.Parallel()

				idx, err := index.Parse(true, "foo")
				require.NoError(t, err)
				require.Equal(
					t,
					&index.Index{
						Name:    "foo_idx",
						Columns: []string{"foo"},
					},
					idx,
				)
			})
			t.Run("unique index", func(t *testing.T) {
				t.Parallel()

				idx, err := index.Parse("unique", "foo")
				require.NoError(t, err)
				require.Equal(
					t,
					&index.Index{
						Name:    "foo_uniq_idx",
						Columns: []string{"foo"},
						Unique:  true,
					},
					idx,
				)
			})

			t.Run("input is a map with index definition", func(t *testing.T) {
				t.Parallel()

				idx, err := index.Parse(map[string]interface{}{
					"name":   "foo_idx",
					"fields": []string{"foo", "bar"},
				}, "")
				require.NoError(t, err)
				require.Equal(
					t,
					&index.Index{
						Name:    "foo_idx",
						Columns: []string{"foo", "bar"},
					},
					idx,
				)
			})
		})

		t.Run("errors", func(t *testing.T) {
			t.Run("boolean passed as string", func(t *testing.T) {
				t.Parallel()

				// actually this probably should return an error.
				// TODO
				idx, err := index.Parse("true", "foo")
				require.Nil(t, idx)
				require.Nil(t, err)
			})
		})
	})
}
