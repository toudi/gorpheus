package sqlite3

import (
	"github.com/toudi/gorpheus/v2/internal/ddl/column"
	columnPkg "github.com/toudi/gorpheus/v2/internal/ddl/column"

	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestColumn(t *testing.T) {
	var buffer strings.Builder
	var column *column.Column
	var err error
	var s3 *sqlite3Dialect = &sqlite3Dialect{}
	t.Run("happy path", func(t *testing.T) {
		t.Run("varchar", func(t *testing.T) {
			buffer.Reset()
			column = &columnPkg.Column{
				Name:      "name",
				TypeAlias: "string",
				Type:      columnPkg.TypeString,
				MaxLength: 32,
			}
			err = s3.ColumnDDL(column, &buffer)
			require.NoError(t, err)
			require.Equal(t, "name VARCHAR (32) NOT NULL", buffer.String())
		})
		t.Run("varchar with a default value", func(t *testing.T) {
			buffer.Reset()
			column = &columnPkg.Column{
				Name:      "name",
				TypeAlias: "string",
				Default:   "'foo'",
				Type:      columnPkg.TypeString,
				MaxLength: 32,
			}
			err = s3.ColumnDDL(column, &buffer)
			require.NoError(t, err)
			require.Equal(t, "name VARCHAR (32) NOT NULL DEFAULT 'foo'", buffer.String())
		})
		t.Run("primary key", func(t *testing.T) {
			buffer.Reset()
			column = &columnPkg.Column{
				Name:       "name",
				TypeAlias:  "int",
				Type:       columnPkg.TypeInt,
				PrimaryKey: lo.ToPtr(true),
			}
			err = s3.ColumnDDL(column, &buffer)
			require.NoError(t, err)
			require.Equal(t, "name INTEGER PRIMARY KEY NOT NULL", buffer.String())
		})
		t.Run("primary key autoincrement", func(t *testing.T) {
			buffer.Reset()
			column = &columnPkg.Column{
				Name:          "name",
				TypeAlias:     "int",
				Type:          columnPkg.TypeInt,
				PrimaryKey:    lo.ToPtr(true),
				AutoIncrement: lo.ToPtr(true),
			}
			err = s3.ColumnDDL(column, &buffer)
			require.NoError(t, err)
			require.Equal(t, "name INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL", buffer.String())
		})

		t.Run("various column types", func(t *testing.T) {
			type testCase struct {
				column   *columnPkg.Column
				expected string
			}

			for _, test := range []testCase{
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "string",
						MaxLength: 32,
					},
					expected: "name VARCHAR (32) NOT NULL",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "text",
					},
					expected: "name TEXT NOT NULL",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "float",
					},
					expected: "name REAL NOT NULL",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "float64",
					},
					expected: "name DOUBLE PRECISION NOT NULL",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "bool",
					},
					expected: "name BOOLEAN NOT NULL",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "smallint",
					},
					expected: "name SMALLINT NOT NULL",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "int",
						Type:      columnPkg.TypeInt,
					},
					expected: "name INTEGER NOT NULL",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "bigint",
					},
					expected: "name BIGINT NOT NULL",
				},
				{
					column: &columnPkg.Column{
						Name:          "name",
						TypeAlias:     "decimal",
						MaxDigits:     10,
						DecimalPlaces: 2,
					},
					expected: "name NUMERIC (10, 2) NOT NULL",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "string",
						MaxLength: 32,
						Generated: &columnPkg.GeneratedColumn{
							Definition: "upper(some_other_field)",
						},
					},
					expected: "name VARCHAR (32) GENERATED ALWAYS AS (upper(some_other_field)) VIRTUAL",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "string",
						MaxLength: 32,
						Generated: &columnPkg.GeneratedColumn{
							Definition: "upper(some_other_field)",
							Stored:     true,
						},
					},
					expected: "name VARCHAR (32) GENERATED ALWAYS AS (upper(some_other_field)) STORED",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "date",
					},
					expected: "name DATE NOT NULL",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "time",
					},
					expected: "name TIME NOT NULL",
				},
				{
					column: &columnPkg.Column{
						Name:      "name",
						TypeAlias: "timestamp",
					},
					expected: "name DATETIME NOT NULL",
				},
			} {
				t.Run(test.expected, func(t *testing.T) {
					buffer.Reset()
					err = s3.ColumnDDL(test.column, &buffer)
					require.NoError(t, err)
					require.Equal(t, test.expected, buffer.String())
				})
			}
		})
	})

	t.Run("errors", func(t *testing.T) {
		t.Run("column must define a type alias", func(t *testing.T) {
			column = &columnPkg.Column{
				Name:      "name",
				TypeAlias: "string",
				Type:      columnPkg.TypeString,
			}
			err = s3.ColumnDDL(column, &buffer)
			require.Errorf(t, columnPkg.ErrInvalidColumnType, "")
		})
		t.Run("column of type string must have a length", func(t *testing.T) {
			column = &columnPkg.Column{
				Name:      "name",
				TypeAlias: "string",
				Type:      columnPkg.TypeString,
			}
			err = s3.ColumnDDL(column, &buffer)
			require.Equal(t, columnPkg.ErrInvalidMaxLength, err)
		})
		t.Run("unsupported type of int128", func(t *testing.T) {
			column = &columnPkg.Column{
				Name:      "name",
				TypeAlias: "int128",
				Type:      columnPkg.TypeInt128,
			}
			err = s3.ColumnDDL(column, &buffer)
			require.ErrorIs(t, err, ErrUnsupportedColumnType, "int128")
		})
		t.Run("unsupported autoincrement for type string", func(t *testing.T) {
			column = &columnPkg.Column{
				Name:          "name",
				TypeAlias:     "string",
				MaxLength:     1,
				Type:          columnPkg.TypeString,
				AutoIncrement: lo.ToPtr(true),
			}
			err = s3.ColumnDDL(column, &buffer)
			require.ErrorIs(t, err, ErrUnsupportedTypeForAutoIncrement)
		})
	})
}
