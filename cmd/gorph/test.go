package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus"
	"github.com/toudi/gorpheus/storage"
)

func printCollection(c *gorpheus.Collection) {
	for _, m := range c.Versions {
		log.Debugf("[ ] %s\n", m.Revision())
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("gorpheus started")
	collection := gorpheus.Collection_init()
	storage.ScanDirectory("migrations", collection)
	log.Debugf("Original collection")
	printCollection(collection)
	collection.Sort()
	log.Debugf("Collection after sorting:")
	printCollection(collection)
	// connect to db
	db, _ := sqlx.Open("sqlite3", "./example.sqlite3")
	tx := db.MustBegin()
	err := collection.MigrateUpTo(tx, "0001")
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	// fmt.Println("Migration names")
	// fmt.Println(gorpheus.MigrationNames(collection))
}
