package operation

import (
	"errors"
	"strings"
)

var (
	ErrCannotParseFunctionName = errors.New("cannot parse the function name")
)

type CreateFunction struct {
	Definition string `yaml:"definition"`
}

func (o *CreateFunction) GetFunctionName() (string, error) {
	parenthesisIndex := strings.Index(o.Definition, "(")
	if parenthesisIndex == -1 {
		return "", ErrCannotParseFunctionName
	}

	return strings.TrimLeft(o.Definition[:parenthesisIndex], " "), nil
}
