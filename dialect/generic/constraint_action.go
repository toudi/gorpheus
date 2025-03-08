package generic

import (
	"errors"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/constraint"
)

var ErrUnknownPolicy = errors.New("unknown constraint policy")

func ConstraintAction(policy *string, buffer *strings.Builder) error {
	policyConstraint, exists := constraint.PolicyFromAlias[*policy]
	if !exists {
		return ErrUnknownPolicy
	}
	switch policyConstraint {
	case constraint.CASCADE:
		buffer.WriteString("CASCADE")
	case constraint.NO_ACTION:
		buffer.WriteString("NO ACTION")
	case constraint.SET_DEFAULT:
		buffer.WriteString("SET DEFAULT")
	case constraint.SET_NULL:
		buffer.WriteString("SET NULL")
	case constraint.RESTRICT:
		buffer.WriteString("RESTRICT")
	}
	return nil
}
