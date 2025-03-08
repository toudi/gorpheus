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
