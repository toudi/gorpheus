package operation

import (
	"github.com/toudi/gorpheus/v2/internal/ddl/column"
	"github.com/toudi/gorpheus/v2/internal/ddl/constraint"
	"github.com/toudi/gorpheus/v2/internal/ddl/index"
)

func columnSideEffects(c *column.Column, tableName string) ([]Operation, error) {
	var sideEffects []Operation

	if c.Index != nil {
		idx, err := index.Parse(c.Index, c.Name)
		if err != nil {
			return nil, err
		}
		// let's add the index to sideEffects queue
		sideEffects = append(sideEffects, Operation{
			CreateIndex: &CreateIndex{
				Table: tableName,
				Index: idx,
			},
		})
	}

	if c.ForeignKey != nil {
		foreignKey := c.ForeignKey
		foreignKey.SrcColumns = []string{c.Name}
		sideEffects = append(sideEffects, Operation{
			AddConstraint: &AddConstraint{
				Table: tableName,
				Null:  c.Null,
				Constraint: &constraint.Constraint{
					Name:       c.Name + "_fk",
					ForeignKey: foreignKey,
				},
				SideEffect: true,
			},
		})
	}

	if c.Unique != nil {
		sideEffects = append(sideEffects, Operation{
			AddConstraint: &AddConstraint{
				Table: tableName,
				Null:  c.Null,
				Constraint: &constraint.Constraint{
					Name: c.Name + "_uniq",
					Unique: &constraint.Unique{
						Columns: []string{c.Name},
					},
				},
				SideEffect: true,
			},
		})
	}

	return sideEffects, nil
}
