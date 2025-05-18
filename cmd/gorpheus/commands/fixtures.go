package commands

import (
	"database/sql"
	"strings"

	"github.com/toudi/gorpheus/v2/fixtures"
)

type LoadData struct {
	databaseUrl
	FixtureFile string `arg:"" type:"existingfile" help:"fixture file to load into the database (either JSON or YAML)" required:"true"`
}
type DumpData struct {
	databaseUrl
	Table  string `name:"table" short:"t" help:"table name to dump data from"`
	Output string `name:"output" short:"o" type:"path" help:"path to output file. either JSON or YAML"`
}

func (dd *DumpData) Run(ctx *Context) error {
	if err := dd.handleDbUrl(ctx); err != nil {
		return err
	}

	var ff = fixtures.FixtureFile{
		Table: dd.Table,
	}

	return ctx.Gorpheus.Transaction(func(tx *sql.Tx) error {
		var err error

		rows, err := tx.Query("SELECT * FROM " + ff.Table)
		if err != nil {
			return err
		}
		defer rows.Close()

		if err = ff.ReadDataFromQuery(rows); err != nil {
			return err
		}

		return ff.Save(dd.Output)
	})
}

func (ld *LoadData) Run(ctx *Context) error {
	if err := ld.handleDbUrl(ctx); err != nil {
		return err
	}

	ff, err := fixtures.LoadFromFile(ld.FixtureFile)
	if err != nil {
		return err
	}

	return ctx.Gorpheus.Transaction(func(tx *sql.Tx) error {
		var err error
		// let's iterate over the rows one by one and insert them into the database
		for _, row := range ff.Data {
			var columnNames = make([]string, len(row))
			var columnValues = make([]any, len(row))
			var columnValuePlaceholders = make([]string, len(row))
			var colIdx int = 0

			for colName, colValue := range row {
				columnNames[colIdx] = colName
				columnValues[colIdx] = colValue
				columnValuePlaceholders[colIdx] = "?"
				colIdx++
			}

			if _, err = tx.Exec("INSERT INTO "+ff.Table+" ("+strings.Join(columnNames, ", ")+") VALUES ("+strings.Join(columnValuePlaceholders, ",")+")", columnValues...); err != nil {
				return err
			}

		}

		return nil
	})
}
