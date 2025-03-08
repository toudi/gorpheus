package index

import (
	"errors"
	"fmt"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/utils"
)

var ErrInvalidDefinition = errors.New("invalid index definition")

type Index struct {
	Name       string   `yaml:"name"`
	Columns    []string `yaml:"fields"`
	Unique     bool     `yaml:"unique"`
	Definition string   `yaml:"definition"`
}

func (i Index) GetName() string {
	if i.Name != "" {
		return i.Name
	}

	var idxNameParts []string = i.Columns

	if i.Unique {
		idxNameParts = append(idxNameParts, "uniq")
	}

	idxNameParts = append(idxNameParts, "idx")

	return strings.Join(idxNameParts, "_")
}

func Parse(input interface{}, column string) (*Index, error) {
	if _, ok := input.(bool); ok {
		return &Index{
			Name:    fmt.Sprintf("%s_idx", column),
			Columns: []string{column},
		}, nil
	} else if tmpString, ok := input.(string); ok {
		if strings.EqualFold(tmpString, "unique") {
			return &Index{
				Name:    fmt.Sprintf("%s_uniq_idx", column),
				Columns: []string{column},
				Unique:  true,
			}, nil
		}
	} else if tmpMap, ok := input.(map[string]interface{}); ok {
		// let's try to parse the input map.
		idx := &Index{}
		err := utils.Recast(tmpMap, &idx)
		return idx, err
	}
	return nil, nil
}
