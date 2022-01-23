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

func main() {
	var err error

	var params gorpheus.MigrationParams

	flag.StringVar(&params.EnvKeyName, "env", "DATABASE_URL", "Environment variable that will contain database url")
	flag.BoolVar(&params.Fake, "fake", false, "This option is only meaningful if combined with namespace and target revision. Then, if set, it pretends that a migration was applied.")

	flag.Parse()

	if flag.NArg() > 1 {
		fmt.Println(flag.NArg())
		// that's a namespace and potentially a revision as well.
		fmt.Println(flag.Args())
		params.Namespace = flag.Arg(0)
		if flag.NArg() == 2 {
			if flag.Arg(1) == "zero" {
				params.Zero = true
			} else {
				if params.VersionNo, err = strconv.Atoi(flag.Arg(1)); err != nil {
					fmt.Printf("unable to parse revision: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}

	log.SetLevel(log.DebugLevel)
	log.Debug("gorpheus started")
	collection := gorpheus.Collection_init()
	storage.ScanDirectory("migrations", collection)
	collection.Register(newMigration)
	collection.Register(newMigration2)
	log.Debugf("Collection")
	printCollection(collection)
	log.Debug("Performing migrations")
	err = collection.Migrate(&params)
	//err := collection.MigrateUp(db)
	if err != nil {
		log.WithError(err).Error("Cannot perform migrations")
	}
}
