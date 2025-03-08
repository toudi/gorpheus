package executor

import (
	"errors"

	"github.com/toudi/gorpheus/v2/internal/ddl/column"
	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
	"github.com/toudi/gorpheus/v2/internal/utils"
)

// the db state object stores the internal db state as seen by
// migrations. this allows the migration writer to only specify
// changes that are relevant (e.g. change the columns width)
// and the db state will yield the deltas.
// it should also be noted that the vast majority of the changes can
// be reversed directly, without the need of this table. However,
// alter column operations do need to track of the column's internal
// state.

type columnState struct {
	currentVersion int
	availableSince int // which revision of the table first references this column
	revisions      []map[string]interface{}
}
type tableState struct {
	version int // each preparation operation on the table bumps this number
	// in other words, this is a "theoretical" version of the table
	// incremented after each altering operation, but it does not
	// represent the "live" / actual version of the table
	migratedVersion int // each operation performed on the table bumps this number
	// this means that this represents the actual version at which the table is in.
	columns map[string]*columnState
}

type dbState struct {
	tables   map[string]*tableState
	executor *Executor
}

var _executor *Executor

func (cs *columnState) getCurrentColumnRevision(columnName string) (*column.Column, error) {
	var rawColumnData map[string]interface{}

	for revision, revisionData := range cs.revisions {
		if revision > cs.currentVersion {
			break
		}
		if rawColumnData == nil {
			rawColumnData = make(map[string]interface{})
		}
		for fieldName, fieldValue := range revisionData {
			rawColumnData[fieldName] = fieldValue
		}
	}

	return _executor.parseField(columnName, rawColumnData)
}

func (d *dbState) AddColumnVersion(table string, columnName string, column map[string]interface{}) {
	if _, ok := d.tables[table]; !ok {
		d.tables[table] = &tableState{}
	}
	if d.tables[table].columns == nil {
		d.tables[table].columns = make(map[string]*columnState)
	}
	if _, exists := d.tables[table].columns[columnName]; !exists {
		d.tables[table].columns[columnName] = &columnState{
			availableSince: d.tables[table].version,
		}
	}

	d.tables[table].columns[columnName].revisions = append(
		d.tables[table].columns[columnName].revisions, column,
	)
	d.tables[table].columns[columnName].currentVersion += 1
	currentVersion := d.tables[table].columns[columnName].currentVersion
	// we have to copy the type so that the alter column carries it over as well..
	if currentVersion >= 2 {
		_, exists := d.tables[table].columns[columnName].revisions[currentVersion-1]["type"]
		// but only if it was not overriden by the user.
		// for instance, it's allowed to change the column type from, say, smallint to bigint
		if !exists {
			d.tables[table].columns[columnName].revisions[currentVersion-1]["type"] = d.tables[table].columns[columnName].revisions[currentVersion-2]["type"]
		}
	}
}

func (d *dbState) GetColumnVersion(table string, columnName string) int {
	return d.tables[table].columns[columnName].currentVersion
}

// this function returns the changed attributes
func (d *dbState) Delta(alter *operation.AlterField, versionFrom int, versionTo int) (err error) {
	var table = alter.TableName
	var columnName = alter.FieldName

	var columnFrom = d.tables[table].columns[columnName].revisions[versionFrom]
	var columnTo = d.tables[table].columns[columnName].revisions[versionTo]

	if versionTo > versionFrom {
		// if this is a forwards migration then it's actually trivial - we simply have to deserialize column state
		// into altered attributes
		if err = utils.Recast(columnTo, &alter.Delta); err != nil {
			return err
		}
	} else {
		// we have to iterate over `from` attributes, read their corresponding properties from `to`
		// and apply that..
		var rawDelta map[string]interface{} = make(map[string]interface{})
		for propertyName := range columnFrom {
			rawDelta[propertyName] = columnTo[propertyName]
		}
		// columnFrom = rawDelta
		if err = utils.Recast(rawDelta, &alter.Delta); err != nil {
			return err
		}
	}

	// let's see if the column data type has changed:
	if utils.IsDifferent(columnFrom, columnTo) {
		// this means that the definition has changed. let's save it:
		toColumn, err := d.executor.parseField(columnName, columnTo)
		if err != nil {
			return err
		}
		alter.Column = toColumn
	}

	return err
}

func (st *dbState) TableColumns(tableName string) ([]*column.Column, error) {
	tableState, exist := st.tables[tableName]
	if !exist {
		return nil, errors.New("unknown table")
	}

	var columns []*column.Column

	for columnName, columnState := range tableState.columns {
		if columnState.availableSince > tableState.migratedVersion {
			continue
		}
		column, err := columnState.getCurrentColumnRevision(columnName)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}

func (st *dbState) IncrementTableVersion(tableName string) error {
	tableState, exists := st.tables[tableName]
	if !exists {
		return errors.New("unknown table")
	}
	tableState.version += 1
	return nil
}

func (st *dbState) AdjustMigratedTableVersion(adjust *operation.AdjustTableVersion) error {
	tableState, exists := st.tables[adjust.TableName]
	if !exists {
		return errors.New("unknown table")
	}
	tableState.migratedVersion += adjust.Delta
	return nil
}

func (st *dbState) SetFieldVersion(tableName, columnName string, newCurrentRevision int) error {
	tableState, exists := st.tables[tableName]
	if !exists {
		return errors.New("unknown table")
	}
	columnState, exists := tableState.columns[columnName]
	if !exists {
		return errors.New("unknown column")
	}
	columnState.currentVersion = newCurrentRevision
	return nil
}

func (st *dbState) AdjustVersions(o operation.Operation) {
	var tableName string

	if o.AlterField != nil {
		tableName = o.AlterField.TableName
		st.SetFieldVersion(
			tableName,
			o.AlterField.FieldName,
			o.AlterField.FieldVersion,
		)
	}
}

func NewDBState(e *Executor) *dbState {
	_executor = e
	return &dbState{
		tables:   make(map[string]*tableState),
		executor: e,
	}
}
