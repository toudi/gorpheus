package operation_test

import (
	"strings"
	"testing"

	"github.com/toudi/gorpheus/v2/internal/ddl/index"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"

	"github.com/stretchr/testify/require"
)

func TestCreateIndex(t *testing.T) {
	type testCase struct {
		expected string
		input    operation.CreateIndex
	}

	var buffer strings.Builder

	for _, test := range []testCase{
		{
			expected: "CREATE INDEX IF NOT EXISTS col_idx ON test_table (col);",
			input: operation.CreateIndex{
				Table: "test_table",
				Index: &index.Index{
					Columns: []string{"col"},
				},
			},
		},
		{
			expected: "CREATE UNIQUE INDEX IF NOT EXISTS col_uniq_idx ON test_table (col);",
			input: operation.CreateIndex{
				Table: "test_table",
				Index: &index.Index{
					Columns: []string{"col"},
					Unique:  true,
				},
			},
		},
		{
			expected: "CREATE INDEX IF NOT EXISTS index_name ON test_table (col);",
			input: operation.CreateIndex{
				Table: "test_table",
				Index: &index.Index{
					Name:    "index_name",
					Columns: []string{"col"},
				},
			},
		},
		{
			expected: "CREATE INDEX IF NOT EXISTS col_col2_idx ON test_table (col, col2);",
			input: operation.CreateIndex{
				Table: "test_table",
				Index: &index.Index{
					Columns: []string{"col", "col2"},
				},
			},
		},
	} {
		t.Run(test.expected, func(t *testing.T) {
			buffer.Reset()

			err := test.input.DDL(&buffer)
			require.Equal(t, test.expected, buffer.String())
			require.NoError(t, err)
		})
	}
}
