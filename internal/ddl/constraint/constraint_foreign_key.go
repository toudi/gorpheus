package constraint

import (
	"errors"
	"slices"
)

const (
	NO_ACTION = iota
	CASCADE
	SET_NULL
	SET_DEFAULT
	RESTRICT
)

var ErrNoColumnsInForeignKey = errors.New("no columns in foreign key")

type ForeignKey struct {
	SrcColumns []string `yaml:"from"`
	Table      string   `yaml:"table"`
	Column     string   `yaml:"column"`
	Columns    []string `yaml:"columns"`
	OnUpdate   *string  `yaml:"on-update"`
	OnDelete   *string  `yaml:"on-delete"`
}

func (fk *ForeignKey) Validate() error {
	if !slices.Contains(fk.Columns, fk.Column) {
		fk.Columns = append(fk.Columns, fk.Column)
	}

	if fk.Columns == nil {
		return ErrNoColumnsInForeignKey
	}

	return nil
}

var policyToAlias = map[int][]string{
	NO_ACTION:   {"no action", "no-action"},
	CASCADE:     {"cascade"},
	SET_NULL:    {"set null", "set-null", "null"},
	SET_DEFAULT: {"set default", "set-default", "default"},
	RESTRICT:    {"restrict"},
}

var PolicyFromAlias map[string]int

func init() {
	PolicyFromAlias = make(map[string]int)

	for policy, aliases := range policyToAlias {
		for _, alias := range aliases {
			PolicyFromAlias[alias] = policy
		}
	}
}
