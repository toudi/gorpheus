package main

import (
	"context"
	"log/slog"

	"gorpheus-cli/commands"

	"github.com/toudi/gorpheus/v2/executor"
	"github.com/toudi/gorpheus/v2/interfaces"

	_ "gorpheus-cli/dialects"

	"github.com/alecthomas/kong"
	"github.com/lmittmann/tint"
)

var CLI struct {
	Verbosity int               `short:"v" type:"counter" help:"Set verbosity level"`
	Migrate   commands.Migrate  `cmd:"" help:"migrate database"`
	Show      commands.Show     `cmd:"" help:"show applied migrations"`
	Loaddata  commands.LoadData `cmd:"" help:"load data from fixture file into database"`
	Dumpdata  commands.DumpData `cmd:"" help:"dump data from database to a fixture file"`
}

func main() {
	ctx := kong.Parse(&CLI)

	logger := commands.NewLogger(context.Background(), CLI.Verbosity)

	gorpheus := executor.New(executor.WithLogger(logger))

	err := ctx.Run(&commands.Context{
		Gorpheus: gorpheus,
		Logger:   logger,
	})
	if err != nil {
		logger.Log(interfaces.LogMessage{
			Level:   slog.LevelDebug,
			Message: "unable to run subcommand",
			Payload: []any{
				tint.Err(err),
			},
		})
	}
}
