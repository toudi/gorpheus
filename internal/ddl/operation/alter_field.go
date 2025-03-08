package operation

import (
	"github.com/toudi/gorpheus/v2/internal/ddl/column"
)

type AlteredAttributes struct {
	Null    *bool        `yaml:"null"`
	Default *interface{} `yaml:"default"`
}

type AlterField struct {
	TableName string `yaml:"table"`
	FieldName string `yaml:"name"`
	// this contains changed attributes of a column
	FieldDelta   rawFieldType `yaml:"field"`
	FieldVersion int
	Column       *column.Column // this property will only be populated when something that affects
	// column type will change. For instance - if a column changes from Varchar(32) to varchar(512)
	// or a column changes from smallint to bigint and so on.
	Delta AlteredAttributes
}
