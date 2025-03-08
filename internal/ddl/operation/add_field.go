package operation

import (
	"github.com/toudi/gorpheus/v2/internal/ddl/column"
)

type AddField struct {
	TableName string        `yaml:"table"`
	FieldName string        `yaml:"name"` // only applicable when `Field` is populated
	Field     *rawFieldType `yaml:"field"`
	FieldsMap fieldsMapType `yaml:"fields"`
	Column    *column.Column
	Columns   []*column.Column `yaml:"columns"`
}

func (af *AddField) SingleFieldAdds() []Operation {
	var operations []Operation

	for _, column := range af.Columns {
		operations = append(operations, Operation{
			AddField: &AddField{
				TableName: af.TableName,
				Column:    column,
			},
		})
	}

	return operations
}

func (af *AddField) SideEffects() ([]Operation, error) {
	return columnSideEffects(af.Column, af.TableName)
}
