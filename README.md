gorpheus.

This code is a proof of concept migation manager.

You can compile it like so:

`go build cmd/gorph/test.go`

it demonstrates the two possibilities to manage database migations:

1. Trough files
2. Trough in-memory structs

Here's how this works in principle:

First, you have to initialize the migation collection:

`c := gorpheus.Collection_init()`

then you can register the migrations. If you want to recursively go trough a
directory, do the following:

`storage.ScanDirectory("some-dir", c)`

you can also add some in-memory migations too:

```
var newMigration = &NewMigration{
	Migration: migration.Migration{
		Version:   "0003_inmem",
		Namespace: "users",
		Depends:   []string{"users/0002_something"},
	},
}

func (n *NewMigration) UpScript() (string, uint8, error) {
	return `create_table("users_inmemory")`, migration.TypeFizz, nil
}

func (n *NewMigration) DownScript() (string, uint8, error) {
	return `DROP TABLE users_inmemory`, migration.TypeSQL, nil
}

// don't forget to register it:

c.Register(newMigration)
```

you can register the migrations in any order you wish - if they have the dependencies,
they will be sorted. If they don't have dependencies, sorting defaults to
name sorting - i.e. default/0002_anything will come after default/0001_anything.

each migration has:

- version (required)
- namespace (this is an optional parameter)
  When omitted, it defaults to "default". The reason for this parameter is so
  that one could have multiple directories with migrations. This is somewhat
  similar to what django offers in it's apps - each app has it's own directory
  and the app directory has it's own migations.
- dependencies (this is an optional parameter)
  this is an array of strings which describe the dependencies for migration.
  each dependency must be constructed in a form of `namespace/version`

the migration also defines two functions:

`UpScript() (string, uint8, error)`
this function returns an UP script, it's type and (potentially) an error. Currently
there are two types supported: Fizz and SQL

`DownScript() (string, uint8, error)`
like the above, just for Down operation.

then you can migrate your collection:

```
	tx := db.MustBegin()
	log.Debug("Migrating up")
	err := c.MigrateUp(tx)
	if err != nil {
		log.WithError(err).Error("Cannot migrate up")
		tx.Rollback()
	} else {
		tx.Commit()
	}
```

migrating down is equally simple:

```
    tx := db.MustBegin()
    err := c.MigrateDownTo("default/0001")
    ...
```
