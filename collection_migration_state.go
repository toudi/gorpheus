package gorpheus

import (
	"errors"
	"time"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const createFizz = `
create_table("gorpheus_revisions") {
	t.Column("revision", "string", {"size": 32})
}
`

const findSQL = `SELECT COUNT(1) FROM gorpheus_revisions WHERE revision = ?`
const insertSQL = `INSERT INTO gorpheus_revisions (revision, created_at, updated_at) VALUES (?, ?, ?)`
const deleteSQL = `DELETE FROM gorpheus_revisions WHERE revision = ?`

func (c *Collection) createTable(tx *sqlx.Tx) error {
	var bubbler *fizz.Bubbler

	if tx.DriverName() == "sqlite3" {
		bubbler = fizz.NewBubbler(translators.NewSQLite("sqlite://example.sqlite3"))
	} else {
		return errors.New("Unsupported driver")
	}

	createSQL, err := bubbler.Bubble(createFizz)
	if err != nil {
		log.WithError(err).Error("Cannot translate fizz migration")
		return err
	}
	_, err = tx.Exec(createSQL)
	return err
}

func (c *Collection) Exists(tx *sqlx.Tx, revision string) (bool, error) {
	var result uint8
	row := tx.QueryRow(findSQL, revision)
	err := row.Scan(&result)
	if err != nil {
		return false, c.createTable(tx)
	}
	return result == 1, err
}

func (c *Collection) InsertRevision(tx *sqlx.Tx, revision string) error {
	_, err := tx.Exec(insertSQL, revision, time.Now().String(), time.Now().String())
	return err
}

func (c *Collection) RemoveRevision(tx *sqlx.Tx, revision string) error {
	_, err := tx.Exec(deleteSQL, revision)
	return err
}
