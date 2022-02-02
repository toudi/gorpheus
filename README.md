# gorpheus

This code is a proof of concept migration manager. Please refer to examples directory for practical usage examples.

## When it might and when it might not be helpful to me?

I wrote this code because I was looking for a library (in go) that would implement migrations as django does them. This means:

- the migrations are database-backend agnostic (i.e. you write migrations once and can execute them on any database you wish)
- the migrations can be split across multiple directories
- the migrations can depend on each other
- the migrations can be applied within a single namespace (i.e. `./manage.py migrate my_namespace` in Django)
- the migrations have human-readable names (i.e. 0001, 0002, and so on)
- the migration can be a python function (so-called data migrations)
- the forwards and backwards actions are specified within the same file / module

There are several libraries that implement migrations in go, however the one that I saw required that the migrations would be in a single directory, or would have versions that would
depend on timestamps or that forwards / backwards actions would be split into two files, or they did not support running go code as a migration at all.

while using Django, I quickly discovered that relying on timestamps is not that helpful when it comes to numbering. Imagine for a second that 2 or more developers start working on their
own branches of the code and each of the developers needs to modify the database. If they do that with incrementing numbers as a base then all that needs to be done during the merge process
is renaming of the files. it is still possible to do that with timestamps of course, but it's not as user-friendly.

On the other hand, if you know exactly what would be the target database for your code then there's hardly any point in using my code - the other tools are definetely more suitable for that.

## Basic concepts

Ok, if you made it here this means that you think the project might be beneficial to you :-) Here are the key concepts of the project.

### Collection

gorpheus organizes migrations in a collection. Each item of the collection must implement a `Migration` interface. The reason for this
is simple - I wanted to support many possible sources of migration (file / embedded files / S3 storage, you name it.)

### Migration

Each element of the `Collection` array is a `Migration` interface. Each migration consists of the following things:

#### Type (`uint8`)

| Value | Description |
| --- | --- |
| `TypeSQL` | is straight forward - you need to provide the SQL that will be executed during forward / backwards migration. |
| `TypeFizz` | will be translated to the target dialect with the help of wonderful gobuffalo/fizz library. |
| `TypeGo` | is not an interpreted migration script, but rather a function that will perform the migration. |

#### Version (`string`)

each migration is represented by a version, which is a human-friendly string that contains the namespace. For instance:

```
       version
/----------------------\
|                      |
users/0001_create_table.sql
^^^^^ ^^^^ ^^^^^^^^^^^^
  |	    |       +---------- name
  |     +------------------ version number
  +------------------------ namespace
```

#### Depends (`[]string`)

dependencies array. If you don't define any, gorpheus will assign the previous one from the same namespace (if it exists). The reason for that behavior
is that gorpheus traverses the migrations collection and builds a graph of dependencies. Of course it would be tedious to manually enter dependencies
in each of the migrations so that's why this process is optional.

#### Reader (`io.ReadSeekCloser`)

this is how gorpheus will be able to read and parse sections of the migration script. This is also how it is possible to support other sources of storage
that I haven't concidered (i.e. S3, FTP, you name it)

### Migration file and sections

Here's an example migration file:

```SQL
-- gorph DEPENDS 
users/0001_create_users 
-- end --
-- gorph UP
CREATE TABLE foo (
    invoice_id INTEGER,
    user_id INTEGER,
    FOREIGN KEY(invoice_id) REFERENCES invoices(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
);
-- end --
-- gorph DOWN
DROP TABLE foo;
-- end --
```

as you can see, it consists of 3 sections

- dependencies section. This one is optional as I already mentioned - if you don't have any depedencies from other namespaces then gorpheus will automatically set the previous migration as a dependency
- up / forwards section. this defines what will be done when forwards migration will be applied
- down / backwards section. this defines what will be done when backwards migration will be applied

### Go / data migrations

What if you want to migrate your database but there's some logic that you want to apply which is hard (or even impossible) to write as an SQL query? This is where Go migrations come in place:

```go
type MyDataMigration struct {
	migration.GoMigration
}

func (dm *MyDataMigration) Up(tx *sqlx.Tx) error {
	_, err := tx.Exec(tx.Rebind("UPDATE users_inmemory SET email = ?;"), "foo@bar.baz")
	return err
}

func (dm *MyDataMigration) Down(tx *sqlx.Tx) error {
	_, err := tx.Exec("UPDATE users_inmemory SET email = null;")
	fmt.Printf("Err=%v", err)
	return err
}

var myDataMigration = &MyDataMigration{
	GoMigration: migration.GoMigration{
		Migration: migration.Migration{
			Version:   "users/0004_inmem",
			Depends:   []string{"users/0003_inmem"},
		},
	},
}

collection.Register(myDataMigration)
```

As you can see, the Go migration receives a database transaction (`tx`) and needs to return the error that will indicate whether a transaction will be commited or rolled back. You could do all sorts of things in here - for-loops, selects, etc, etc. This example is very trivial and could be as well defined as an regular sql migration, however I wanted to show the basics.

For obvious reasons, these migrations cannot live within the filesystem / embedFS but have to be compiled into the binary.

## Tutorial
Ok, now that we have the basics out of the way let's get to some practical tutorial:

First of all, you need to initialize a collection:

```go
c := gorpheus.Collection_init()
```

then you can register the migrations. If you want to recursively go trough a
directory, do the following:

```go
storage.RegisterFSRecurse(c, "some-dir")
```

This will walk recursively trough the directory "some-dir" and will register the migrations. The first level of directory
will get the namespace assigned as `default` and as the function recurses, it will change the namespace based on the subdirectory names.

For example:

```
migrations/0001_something.sql         # will be recognized as default/0001_something.sql
migrations/0002_something.sql         # will be recognized as default/0002_something.sql
migrations/foobar/
migrations/foobar/0001_something.sql  # will be recognized as foobar/0001_something.sql
migrations/foobar/0002_something.sql  # will be recognized as foobar/0001_something.sql
```

If you want to have more control over the namespaces then don't use the recurse function but specify everything by hand:

```go
storage.RegisterFS(c, "users", "users-ns")
```

If you have some files that you want to embed inside the final binary (go 1.16+ is required for this) you can do the following thing:

```go
//go:embed migrations/*
var embeddedMigrations embed.FS

storage.RegisterEmbedFS(c, embeddedMigrations, "namespace")
```

you can also add some in-memory migrations too (I'm not really sure if this is beneficial, however I wrote this code way prior to go 1.16 when there was no embed functionality):

```go
type NewMigration struct {
	migration.Migration
}

var newMigration = &NewMigration{
	Migration: migration.Migration{
		Version:   "users/0003_inmem",
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

## The migration parameters

This is where you inform gorpheus how you wish to proceed - whether you want to migrate only a single namespace, or all of them, or to which revision.

```go
// in all of the examples you can either rely on a environment variable for database connection or specify the connectionURL in the parameters:
var myParams = gorpheus.MigrationParams{
	ConnectionURL: "sqlite://example.sqlite",
}
var myParams = gorpheus.MigratinoParams{
	EnvKeyName: "DATABASE_URL",
}
// I'm only skipping that for brevity

// this is the default action - gorpheus will apply all outstanding forward migrations
var myParams = gorpheus.MigrationParams{}

// this will migrate a single namespace to the latest version (including dependencies)
var myParams = gorpheus.MigrationParams{
	Namespace: "users",
}

// this will migrate a single namespace to a version specified by number
var myParams = gorpheus.MigrationParams{
	Namespace: "users",
	VersionNo: 123,
}

Please note that if the database will be already migrated to, say, version number 125 within this namespace then gorpheus will roll back the namespace instead

// this will instruct gorpheus to roll back all migrations from the namespace "users"
var myParams = gorpheus.MigrationParams{
	Namespace: "users",
	Zero     : true,
}

// this will instruct gorpheus to roll back all migrations from all namespaces and destroy it's own state table
var myParams = gorpheus.MigrationParams{
	Vaccuum: true,
}

// once you've decided on what you want to do, call Migrate function:

err := c.Migrate(&params)

```