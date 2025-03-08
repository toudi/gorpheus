package executor

import (
	"errors"

	"github.com/toudi/gorpheus/v2/dialect"
	"github.com/toudi/gorpheus/v2/dialect/generic"

	"github.com/xo/dburl"
)

var ErrUnknownDialect = errors.New("unknown dialect")

func (e *Executor) SetDialect(dialectName string, dsn string) error {
	var err error

	dialectConstructor, exists := dialect.GetConstructor(dialectName)
	if !exists {
		return errors.Join(ErrUnknownDialect, errors.New(dialectName))
	}

	e.dialect = dialectConstructor(e.dbState, dsn)
	e.history, err = generic.NewMigrationsHistoryManager(e.dialect, e.logger)

	return err
}

func (e *Executor) SetDbURL(url string) error {
	connectionInfo, err := dburl.Parse(url)
	if err != nil {
		return err
	}

	return e.SetDialect(connectionInfo.Driver, url)
}
