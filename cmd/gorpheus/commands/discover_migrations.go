package commands

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/toudi/gorpheus/v2/interfaces"
	"github.com/xo/dburl"
)

type databaseUrl struct {
	DatabaseURL string `name:"db" default:"DATABASE_URL" help:"Either a database url or name of an environment variable that contains one."`
}

func (d *databaseUrl) handleDbUrl(ctx *Context) error {
	// first, let's parse the database url.
	var dbUrl *dburl.URL
	var err error

	ctx.Logger.Log(interfaces.LogMessage{
		Level:   slog.LevelDebug,
		Message: "read database url from",
		Payload: []any{slog.String("DatabaseURL", d.DatabaseURL)},
	})

	dbUrl, err = dburl.Parse(os.Getenv(d.DatabaseURL))

	if err != nil {
		ctx.Logger.Log(interfaces.LogMessage{
			Level:   slog.LevelDebug,
			Message: "reading database url from environment variable failed. trying to parse it as is."},
		)
		// maybe it's an actual url ?
		dbUrl, err = dburl.Parse(d.DatabaseURL)
		if err != nil {
			// there's not much that we can do.
			ctx.Logger.Log(interfaces.LogMessage{
				Level:   slog.LevelError,
				Message: "unable to parse database url",
				Payload: []any{tint.Err(err)},
			})
			return err
		}
	}

	gorpheus := ctx.Gorpheus

	ctx.Logger.Log(interfaces.LogMessage{
		Level:   slog.LevelDebug,
		Message: "parsed",
		Payload: []any{slog.String("database url", dbUrl.String())},
	})
	if err = gorpheus.SetDbURL(dbUrl.String()); err != nil {
		ctx.Logger.Log(interfaces.LogMessage{
			Level:   slog.LevelError,
			Message: "unable to set database url",
			Payload: []any{tint.Err(err)},
		})

		return err
	}

	return nil
}

type migrationsDiscovery struct {
	databaseUrl
	MigrationsPath []string `name:"path" short:"p" type:"existingdir" help:"Add specified path(s) to collection." required:""`
}

func (md *migrationsDiscovery) handleMigrationSources(ctx *Context) error {
	if err := md.handleDbUrl(ctx); err != nil {
		return errors.Join(fmt.Errorf("error handling database url"), err)
	}

	gorpheus := ctx.Gorpheus

	var err error

	for _, path := range md.MigrationsPath {
		ctx.Logger.Log(interfaces.LogMessage{
			Level:   slog.LevelDebug,
			Message: "adding migrations to collection",
			Payload: []any{slog.String("path", path)},
		})

		if err = gorpheus.AddMigrationsFromPath(path); err != nil {
			ctx.Logger.Log(interfaces.LogMessage{
				Level:   slog.LevelError,
				Message: "unable to add migrations from directory",
				Payload: []any{tint.Err(err)},
			})
			return err
		}
	}

	return nil
}
