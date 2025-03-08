package operation

import "github.com/toudi/gorpheus/v2/internal/ddl/constraint"

type AddConstraint struct {
	Table      string                 `yaml:"table"`
	Null       *bool                  `yaml:"null"`
	Constraint *constraint.Constraint `yaml:",inline"`
	SideEffect bool
}
