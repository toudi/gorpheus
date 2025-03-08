package operation_test

import (
	"testing"

	"github.com/toudi/gorpheus/v2/internal/ddl/index"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"

	"github.com/stretchr/testify/require"
)

func TestCreateTable(t *testing.T) {
	t.Run("test side effects", func(t *testing.T) {
		t.Run("empty create table should not yield any side effects", func(t *testing.T) {
			t.Parallel()

			input := &operation.CreateTable{}
			sideEffects, err := input.SideEffects()
			require.NoError(t, err)
			require.Nil(t, sideEffects)
		})

		t.Run("if there are indexes defined, create table should yield them", func(t *testing.T) {
			t.Parallel()

			input := &operation.CreateTable{
				TableName: "test_table",
				Indexes: []*index.Index{
					{
						Name: "test_index",
					},
				},
			}

			sideEffects, err := input.SideEffects()
			require.NoError(t, err)
			require.Equal(t, []operation.Operation{
				{
					CreateIndex: &operation.CreateIndex{
						Table: "test_table",
						Index: &index.Index{
							Name: "test_index",
						},
					},
				},
			}, sideEffects,
			)
		})
	})
}
