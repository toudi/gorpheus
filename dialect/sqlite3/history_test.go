package sqlite3

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHistoryImplementation(t *testing.T) {
	t.Run("test HistoryTableExists", func(t *testing.T) {
		t.Run("should return false when the table does not exist", func(t *testing.T) {
			s3 := &sqlite3Dialect{
				dsn: "sqlite3://:memory:",
			}
			manager := s3.HistoryManager()
			var exists bool
			err := s3.Transaction(func(transaction *sql.Tx) error {
				var err error
				exists, err = manager.HistoryTableExists(transaction)
				return err
			})

			require.NoError(t, err)
			require.False(t, exists)
		})

		t.Run("should return true when the table does exist", func(t *testing.T) {
			s3 := &sqlite3Dialect{
				dsn: "sqlite3://:memory:",
			}
			manager := s3.HistoryManager()
			var exists bool
			err := s3.Transaction(func(transaction *sql.Tx) error {
				// we don't really care for the schema itself, we just need for the table to exist
				_, _ = transaction.Exec("CREATE TABLE migrations_history (id INTEGER PRIMARY KEY AUTOINCREMENT);")
				var err error
				exists, err = manager.HistoryTableExists(transaction)
				return err
			})

			require.NoError(t, err)
			require.True(t, exists)

		})

	})
}
