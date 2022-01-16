package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus/v1"
	"github.com/toudi/gorpheus/v1/migration"
	"github.com/toudi/gorpheus/v1/storage"
)

const cmdMigrate = "migrate"
const cmdRollback = "rollback"

func printCollection(c *gorpheus.Collection) {
	for _, m := range c.Versions {
		log.Debugf("-> %s", m.Revision())
	}
}

type NewMigration struct {
	migration.Migration
}

type NewMigration2 struct {
	migration.GoMigration
}

var newMigration = &NewMigration{
	Migration: migration.Migration{
		Version:   "0003_inmem",
		Namespace: "users",
		Depends:   []string{"users/0002_something"},
	},
}

func (n *NewMigration) UpScript() (string, uint8, error) {
	return `
create_table("users_inmemory") {
  t.Column("email", "string", {"null": true})
  t.DisableTimestamps()
}`, migration.TypeFizz, nil
}

func (n *NewMigration) DownScript() (string, uint8, error) {
	return `drop_table("users_inmemory")`, migration.TypeFizz, nil
	//return `DROP TABLE "users_inmemory";`, migration.TypeSQL, nil
}

var newMigration2 = &NewMigration2{
	GoMigration: migration.GoMigration{
		Migration: migration.Migration{
			Version:   "0004_inmem",
			Namespace: "users",
			Depends:   []string{"users/0003_inmem"},
		},
	},
}

func (n2 *NewMigration2) Up(tx *sqlx.Tx) error {
	_, err := tx.Exec(tx.Rebind("UPDATE users_inmemory SET email = ?;"), "foo@bar.baz")
	return err
}

func (n2 *NewMigration2) Down(tx *sqlx.Tx) error {
	_, err := tx.Exec("UPDATE users_inmemory SET email = null;")
	fmt.Printf("Err=%v", err)
	return err
}

func parseNamespaceAndRevision(flagset *flag.FlagSet, dest *gorpheus.MigrationParams) error {
	var err error

	if flagset.NArg() >= 1 {
		// if there are any positional arguments, the first one is expected to be the namespace
		dest.Namespace = flagset.Arg(0)
		if flagset.NArg() == 2 {
			// if there is a second positional argument, it will be the desired revision within namespace
			if dest.Revision, err = strconv.Atoi(flagset.Arg(1)); err != nil {
				return fmt.Errorf("unable to parse revision: %v", err)
			}
		}
	}

	return nil
}

func main() {
	var err error

	var params gorpheus.MigrationParams

	flag.StringVar(&params.EnvKeyName, "env", "DATABASE_URL", "Environment variable that will contain database url")
	flag.BoolVar(&params.Fake, "fake", false, "This option is only meaningful if combined with namespace and target revision. Then, if set, it pretends that a migration was applied.")
	flag.BoolVar(&params.Zero, "zero", false, "Unnapply all migrations, including gorpheus_migrations table")

	flag.Parse()

	if flag.NArg() > 1 {
		fmt.Println(flag.NArg())
		// that's a namespace and potentially a revision as well.
		fmt.Println(flag.Args())
		params.Namespace = flag.Arg(0)
		if flag.NArg() == 2 {
			if params.Revision, err = strconv.Atoi(flag.Arg(1)); err != nil {
				fmt.Printf("unable to parse revision: %v\n", err)
				os.Exit(1)
			}
		}
	}

	log.SetLevel(log.DebugLevel)
	log.Debug("gorpheus started")
	collection := gorpheus.Collection_init()
	storage.ScanDirectory("migrations", collection)
	collection.Register(newMigration)
	collection.Register(newMigration2)
	log.Debugf("Original collection")
	printCollection(collection)
	collection.Sort()
	log.Debugf("Collection after sorting:")
	printCollection(collection)
	log.Debug("Performing migrations")
	err = collection.Migrate(&params)
	//err := collection.MigrateUp(db)
	if err != nil {
		log.WithError(err).Error("Cannot perform migrations")
	}
}
