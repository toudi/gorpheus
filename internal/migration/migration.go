package migration

import (
	"errors"
	"strconv"
	"strings"

	"github.com/toudi/gorpheus/v2/internal/ddl/operation"
)

var (
	ErrUnknownColumnType = errors.New("unknown column type")
	ErrInvalidId         = errors.New("invalid migation ID")
)

type Migration struct {
	Revision Revision
	// dependencies are parsed dependencies - meaning - the ones that can be
	// passed into the planner object
	Dependencies []Revision
	// raw dependencies are the ones that are being read from the source yaml file
	RawDependencies []string `yaml:"depends"`
	// if all of the migrations will be performed on the same table,
	// you can specify it here.
	Table      string                `yaml:"table"`
	Operations []operation.Operation `yaml:"operations"`
}

func NumericVersion(name string) (version int, err error) {
	// first part is the namespace, second one is the name
	nameParts := strings.SplitN(name, "_", 2)
	if len(nameParts) < 1 {
		return -1, ErrInvalidId
	}
	return strconv.Atoi(nameParts[0])
}

func (m *Migration) IsDependencyOf(other *Migration) bool {
	for _, d := range other.Dependencies {
		if d.Namespace == m.Revision.Namespace {
			// if a dependency has revision >= than `m` then by definition
			// it must mean that m is a dependency of `other`.
			// for instance, imagine that
			// m is namespace/0001_initial
			// and you're trying to establish if m could be a dependency of
			// zzzz/0002_something where zzzz/0002_something has a dependency
			// of namespace/0003_foo. Then because namespace/0003_foo has to
			// be applied first then by definition namespace/0001_initial must
			// also be applied first.
			return d.Version >= m.Revision.Version
		}
	}

	return false
}
