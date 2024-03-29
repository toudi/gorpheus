package gorpheus

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const createFizz = `
create_table("gorpheus_revisions") {
	t.Column("revision", "string", {"size": 32})
	t.Column("applied", "timestamp", {})
	t.DisableTimestamps()
}
`
const dropRevisionsTable = `DROP TABLE gorpheus_revisions`

const findRevisionSQL = `SELECT COUNT(1) FROM gorpheus_revisions WHERE revision = ?`
const insertRevisionSQL = `INSERT INTO gorpheus_revisions (revision, applied) VALUES (?, ?)`
const deleteRevisionSQL = `DELETE FROM gorpheus_revisions WHERE revision = ?`

func (c *Collection) ensureMigrationsTableExists(db *sqlx.DB) error {
	var result uint
	var err error
	var tx *sqlx.Tx

	err = db.Get(&result, "SELECT COUNT(1) FROM gorpheus_revisions")
	if err != nil {
		c.Log(LoggerDebug, LogLevelDebug, "error detected: %v; trying to create revisions table\n", err)
		tx, err = db.Beginx()
		if err != nil {
			return fmt.Errorf("could not initialize a transaction: %v", err)
		}
		if err = c.createMigrationsTable(tx); err != nil {
			return fmt.Errorf("could not create migrations table: %v", err)
		}
		err = tx.Commit()
	}
	return err
}

func (c *Collection) createMigrationsTable(tx *sqlx.Tx) error {
	createSQL, err := c.TranslatedSQL(createFizz)
	c.Log(LoggerSQL, LogLevelDebug, createSQL)
	if err != nil {
		c.Log(LoggerDebug, LogLevelError, "cannot translate fizz migration: %v", err)
		return err
	}
	_, err = tx.Exec(createSQL)
	return err
}

func (c *Collection) dropMigrationsTable(db *sqlx.DB) error {
	_, err := db.Exec(dropRevisionsTable)
	return err
}

func (c *Collection) Exists(tx *sqlx.Tx, revision string) (bool, error) {
	var result uint8
	err := tx.Get(&result, tx.Rebind(findRevisionSQL), revision)
	return result == 1, err
}

func (c *Collection) InsertRevision(tx *sqlx.Tx, revision string) error {
	if _, exist := c.applied[revision]; !exist {
		_, err := tx.Exec(tx.Rebind(insertRevisionSQL), revision, time.Now())
		c.applied[revision] = true
		return err
	}
	return nil
}

func (c *Collection) RemoveRevision(tx *sqlx.Tx, revision string) error {
	if _, exist := c.applied[revision]; exist {
		_, err := tx.Exec(tx.Rebind(deleteRevisionSQL), revision)
		delete(c.applied, revision)
		return err
	}
	return nil
}
