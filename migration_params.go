package gorpheus

const (
	DirectionUp   = iota
	DirectionDown = iota
)

type MigrationParams struct {
	// when specified, the migrations utility will roll back all of the migrations
	Zero bool
	// when specified along with namespace and revision, the migrations utility will
	// skip the actual migration method and only add revision entry to the control table
	Fake bool
	// connection URL for the database to connect to. If empty, gorpheus will try to use environment variable.
	ConnectionURL string
	// the environment variable that holds database URL. This defaults to DATABASE_URL
	EnvKeyName string
	// if specified, only migrations from the matching namespace will be applied
	// (with dependencies)
	Namespace string
	// revision within namespace. A special value of -1 can be used to roll back the specified
	// namespace by exactly 1 revision
	VersionNo int
	// if set, gorpheus will clean the migrations table from the database
	Vacuum bool
}
