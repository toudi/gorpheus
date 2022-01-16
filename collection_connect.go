package gorpheus

import (
	"fmt"
	"os"

	"github.com/gobuffalo/fizz/translators"
	"github.com/jmoiron/sqlx"
	"github.com/xo/dburl"
)

func (c *Collection) connectToDb(params *MigrationParams) (*sqlx.DB, error) {
	databaseUrl := os.Getenv(params.EnvKeyName)
	if databaseUrl == "" {
		return nil, fmt.Errorf("`%s` environment variable is empty", params.EnvKeyName)
	}
	// let's parse the database URL
	url, err := dburl.Parse(databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}
	fmt.Printf("detected driver: %s", url.Driver)
	switch url.Driver {
	case "sqlite3":
		c.SetTranslator(translators.NewSQLite(url.OriginalScheme))
	case "postgres":
		c.SetTranslator(translators.NewPostgres())
	case "mysql":
		c.SetTranslator(translators.NewMySQL(url.OriginalScheme, url.Path))
	default:
	}
	return sqlx.Open(url.Driver, url.DSN)
}
