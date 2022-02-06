package gorpheus

import (
	"fmt"
	"os"

	"github.com/gobuffalo/fizz/translators"
	"github.com/jmoiron/sqlx"
	"github.com/xo/dburl"
)

func (c *Collection) connectToDb(params *MigrationParams) (*sqlx.DB, error) {
	var connectionURL string = params.Connection.ConnectionURL
	var dbConn *sqlx.DB
	var url *dburl.URL
	var err error

	if connectionURL == "" && params.Connection.Conn == nil {
		connectionURL = os.Getenv(params.Connection.EnvKeyName)
		if connectionURL == "" {
			return nil, fmt.Errorf("`%s` environment variable is empty", params.Connection.EnvKeyName)
		}
	}
	// let's parse the database URL
	url, err = dburl.Parse(connectionURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	if params.Connection.Conn != nil {
		dbConn = sqlx.NewDb(params.Connection.Conn, url.Driver)
	} else {
		dbConn, err = sqlx.Open(url.Driver, url.DSN)
	}

	if url.Driver == "" {
		return nil, fmt.Errorf("unable to detect database driver")
	}

	switch url.Driver {
	case "sqlite3":
		c.SetTranslator(translators.NewSQLite(url.DSN))
	case "postgres":
		c.SetTranslator(translators.NewPostgres())
	case "mysql":
		c.SetTranslator(translators.NewMySQL(url.DSN, url.Path))
	default:
	}
	return dbConn, err
}
