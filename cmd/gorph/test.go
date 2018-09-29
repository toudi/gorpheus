package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus"
	"github.com/toudi/gorpheus/migration"
	"github.com/toudi/gorpheus/storage"
)

func printCollection(c *gorpheus.Collection) {
	for _, m := range c.Versions {
		log.Debugf("-> %s", m.Revision())
	}
}

type NewMigration struct {
	migration.Migration
}

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

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("gorpheus started")
	collection := gorpheus.Collection_init()
	storage.ScanDirectory("migrations", collection)
	collection.Register(newMigration)
	log.Debugf("Original collection")
	printCollection(collection)
	collection.Sort()
	log.Debugf("Collection after sorting:")
	printCollection(collection)
	// connect to db
	db, _ := sqlx.Open("sqlite3", "./example.sqlite3")
	tx := db.MustBegin()
	log.Debug("Migrating up")
	err := collection.MigrateUp(tx)
	if err != nil {
		log.WithError(err).Error("Cannot migrate up")
		tx.Rollback()
	} else {
		tx.Commit()
	}
	log.Debugf("Migrating down")
	tx = db.MustBegin()
	err = collection.MigrateDownTo(tx, "users/0001")
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}
