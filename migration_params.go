package gorpheus

import "database/sql"

const (
	DirectionUp   = iota
	DirectionDown = iota
)

type DbConnection struct {
	// the environment variable that holds database URL. This defaults to DATABASE_URL
	EnvKeyName string
	// connection URL for the database to connect to. If empty, gorpheus will try to use environment variable.
	// if not empty, gorpheus will parse the URL to extract the driver name and connect to the database
	ConnectionURL string
	// if you already have established the connection by yourself then you can specify it here. gorpheus will
	// use it instead of creating a new one.
	Conn *sql.DB
}

type MigrationParams struct {
	// when specified, the migrations utility will roll back all of the migrations
	Zero bool
	// when specified along with namespace and revision, the migrations utility will
	// skip the actual migration method and only add revision entry to the control table
	Fake bool
	// if specified, only migrations from the matching namespace will be applied
	// (with dependencies)
	Namespace string
	// revision within namespace. A special value of -1 can be used to roll back the specified
	// namespace by exactly 1 revision
	VersionNo int
	// if set, gorpheus will clean the migrations table from the database
	Vacuum bool
	// if set, gorpheus will use this connection instead of relying on ConnectionURL / EnvKeyName
	Connection DbConnection
}
