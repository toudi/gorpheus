package sqlite3

import (
	"fmt"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
	"github.com/toudi/gorpheus/v2/internal/utils"
)

type TempTableOps struct {
	createTempTable      operation.Operation
	insertIntoSelectFrom operation.Operation
	dropOriginal         operation.Operation
	renameTempTable      operation.Operation
	sideEffects          []operation.Operation
}

func (s3 *sqlite3Dialect) createTempTable(
	srcTableName string,
	includeColumn func(columnName string) bool,
) (*TempTableOps, error) {
	// step 1. Let's create a short random string that we'll append to the new temporary table name.
	suffix := utils.RandStringBytes(10)
	tempTableName := srcTableName + "_" + suffix
	// step 2. Create a new table without the column that we're modifying.
	// in order to do that we have to use the db state and interate through the columns
	// and simply skip the one that we're altering
	createTempTable := &operation.CreateTable{
		TableName: tempTableName,
	}

	var fieldNames []string

	// grab list of columns for this table
	// the db state will take care of selecting the "current" revision.
	tableColumns, err := s3.dbState.TableColumns(srcTableName)
	if err != nil {
		return nil, err
	}

	for _, column := range tableColumns {
		fieldNames = append(fieldNames, column.Name)
		if !includeColumn(column.Name) {
			continue
		}
		createTempTable.Columns = append(createTempTable.Columns, column)
	}

	sideEffects, err := createTempTable.SideEffects()
	if err != nil {
		return nil, err
	}

	var joinedFieldNames = strings.Join(fieldNames, ", ")

	return &TempTableOps{
		createTempTable: operation.Operation{
			CreateTable: createTempTable,
		},
		insertIntoSelectFrom: operation.Operation{
			RunStatement: &operation.RunStatement{
				DDL: fmt.Sprintf(
					"INSERT INTO %s (%s) SELECT %s FROM %s",
					tempTableName,
					joinedFieldNames,
					joinedFieldNames,
					srcTableName,
				),
			},
		},
		dropOriginal: operation.Operation{
			RunStatement: &operation.RunStatement{DDL: "DROP TABLE " + srcTableName},
		},
		renameTempTable: operation.Operation{
			RunStatement: &operation.RunStatement{DDL: "ALTER TABLE " + createTempTable.TableName + " RENAME TO " + srcTableName},
		},
		sideEffects: sideEffects,
	}, nil
}
