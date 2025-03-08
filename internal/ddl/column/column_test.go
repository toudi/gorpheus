package column_test

import (
	"testing"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestColumn(t *testing.T) {
	t.Run("test Validate", func(t *testing.T) {
		t.Run("primary key cannot be null", func(t *testing.T) {
			t.Parallel()

			c := &column.Column{
				TypeAlias:  "int",
				PrimaryKey: lo.ToPtr(true),
				Null:       lo.ToPtr(true),
			}

			err := c.Validate()
			require.ErrorIs(t, err, column.ErrPrimaryKeyCannotBeNull)
		})

		t.Run("unknown type alias", func(t *testing.T) {
			t.Parallel()
			c := &column.Column{TypeAlias: "unknown"}
			err := c.Validate()
			require.ErrorIsf(t, err, column.ErrInvalidColumnType, "unknown")
		})
	})
}
