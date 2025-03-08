package executor

import (
	"github.com/toudi/gorpheus/v2/internal/ddl/column"
	"github.com/toudi/gorpheus/v2/internal/utils"
)

// the purpose of the ParseColumn function is to interpret the "base" column
// properties as well as it's distinct ones. For instance, if we have the
// Decimal column property, the "base" properties would be stuff like
// name, null/not null, primary key and so on whereas the decimal type
// would have precision and scale.

func (e *Executor) parseFieldsMap(fieldsMap map[string]map[string]interface{}, handler func(column *column.Column, rawData map[string]interface{})) error {
	for fieldName, fieldDefinition := range fieldsMap {
		column, err := e.parseField(fieldName, fieldDefinition)
		if err != nil {
			return err
		}

		// call the handler function so that the caller may do something with
		// the parsed column.
		handler(column, fieldDefinition)
	}
	return nil
}

func (e *Executor) parseField(fieldName string, fieldDefinition map[string]interface{}) (*column.Column, error) {
	var column column.Column
	var err error

	// let's start by parsing the "generic" column information.
	if err = utils.Recast(fieldDefinition, &column); err != nil {
		return nil, err
	}

	column.Name = fieldName

	// now we can run the validator
	if err = column.Validate(); err != nil {
		return nil, err
	}

	return &column, nil
}
