package operation

import (
	"github.com/toudi/gorpheus/v2/internal/ddl/column"
	"github.com/toudi/gorpheus/v2/internal/ddl/constraint"
	"github.com/toudi/gorpheus/v2/internal/ddl/index"
)

type CreateTable struct {
	TableName   string                   `yaml:"name"`
	Columns     []*column.Column         // populated based on `fields`
	Indexes     []*index.Index           `yaml:"indexes"`
	Constraints []*constraint.Constraint `yaml:"constraints"` // for foreign keys that span multiple columns
	FieldsMap   fieldsMapType            `yaml:"fields"`
}

func (ct *CreateTable) SideEffects() ([]Operation, error) {
	var sideEffects []Operation

	// the indexes can be defined explicitely
	for _, index := range ct.Indexes {
		sideEffects = append(sideEffects, Operation{
			CreateIndex: &CreateIndex{
				Table: ct.TableName,
				Index: index,
			},
		})
	}

	// or non-explicitely, during the column creation.
	for _, column := range ct.Columns {
		columnSideEffects, err := columnSideEffects(column, ct.TableName)
		if err != nil {
			return nil, err
		}
		for _, cse := range columnSideEffects {
			if cse.AddConstraint != nil {
				ct.Constraints = append(ct.Constraints, cse.AddConstraint.Constraint)
				continue
			}
			sideEffects = append(sideEffects, cse)
		}
	}

	// if there are any constraints, we have to validate them.
	// this is because the unique constraint can specify multiple columns and
	// then it's name and definition are generated dynamically.
	for _, constraint := range ct.Constraints {
		if err := constraint.Validate(); err != nil {
			return nil, err
		}
	}

	return sideEffects, nil
}
