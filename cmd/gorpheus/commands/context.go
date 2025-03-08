package commands

import (
	"github.com/toudi/gorpheus/v2/executor"
	"github.com/toudi/gorpheus/v2/interfaces"
)

type Context struct {
	Gorpheus *executor.Executor
	Logger   interfaces.Logger
}
