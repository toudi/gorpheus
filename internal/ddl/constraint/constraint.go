package constraint

import (
	"errors"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/utils"
)

var (
	ErrUnableToGenerateConstraintName = errors.New("unable to generate constraint name")
	ErrUnableToParseUnique            = errors.New("unable to parse unique")
)

type Constraint struct {
	Name       string
	ForeignKey *ForeignKey `yaml:"foreign-key"`
	SrcUnique  any         `yaml:"unique"`
	Unique     *Unique     `yaml:"-"`
}

func parseUnique(input any) (*Unique, error) {
	var uniq Unique
	var err error

	if tmpString, ok := input.(string); ok {
		uniq.Columns = []string{tmpString}
	} else if uniqueList, ok := input.([]string); ok {
		uniq.Columns = uniqueList
	} else if tmpMap, ok := input.(map[string]any); ok {
		// let's see if this is already a unique object:
		err = utils.Recast(tmpMap, &uniq)
	} else if uniqueList, ok := input.([]any); ok {
		for _, column := range uniqueList {
			if columnStr, ok := column.(string); !ok {
				err = ErrUnableToParseUnique
				break
			} else {
				uniq.Columns = append(uniq.Columns, columnStr)
			}
		}
	} else {
		err = ErrUnableToParseUnique
	}
	return &uniq, err
}

func (c *Constraint) Validate() error {
	if c.SrcUnique != nil {
		unique, err := parseUnique(c.SrcUnique)
		if err != nil {
			return err
		}
		c.Unique = unique
	}

	if c.ForeignKey != nil {
		if err := c.ForeignKey.Validate(); err != nil {
			return err
		}
	}

	if c.Name == "" {
		if c.Unique != nil {
			// we can automatically generate a constraint's name based on columns.
			c.Name = strings.Join(c.Unique.Columns, "_") + "_uniq"
		} else if c.ForeignKey != nil {
			c.Name = c.ForeignKey.Table + "_" + strings.Join(c.ForeignKey.Columns, "_") + "_fkey"
		} else {
			return ErrUnableToGenerateConstraintName
		}
	}
	return nil
}
