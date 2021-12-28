package main

import (
	"github.com/gobuffalo/fizz/translators"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus/v1"
	"github.com/toudi/gorpheus/v1/migration"
	"github.com/toudi/gorpheus/v1/storage"
)

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
}`, migration.TypeFizz, nil
}

func (n *NewMigration) DownScript() (string, uint8, error) {
	return `DROP TABLE users_inmemory`, migration.TypeSQL, nil
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
	_, err := tx.Exec("UPDATE users_inmemory SET email = ?;", "foo@bar.baz")
	return err
}

func (n2 *NewMigration2) Down(tx *sqlx.Tx) error {
	_, err := tx.Exec("UPDATE users_inmemory SET email = null;")
	return err
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("gorpheus started")
	collection := gorpheus.Collection_init()
	collection.SetTranslator(translators.NewSQLite("sqlite:///example.sqlite"))
	storage.ScanDirectory("migrations", collection)
	collection.Register(newMigration)
	collection.Register(newMigration2)
	log.Debugf("Original collection")
	printCollection(collection)
	collection.Sort()
	log.Debugf("Collection after sorting:")
	printCollection(collection)
	// connect to db
	db, _ := sqlx.Open("sqlite3", "./example.sqlite3")
	// tx := db.MustBegin()
	log.Debug("Migrating up")
	err := collection.MigrateUp(db)
	if err != nil {
		log.WithError(err).Error("Cannot migrate up")
		// tx.Rollback()
	}
	// } else {
	// 	tx.Commit()
	// }
	// log.Debugf("Migrating down")
	// tx = db.MustBegin()
	// err = collection.MigrateDownTo(db, "users/0003")
	// if err != nil {
	// 	tx.Rollback()
	// } else {
	// 	tx.Commit()
	// }
}
