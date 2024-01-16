package database

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteDatabase struct {
	db *sql.DB
}

func newSQLiteDatabase(db *sql.DB) *SQLiteDatabase {
    return &SQLiteDatabase{
        db: db,
    }
}

func (r *SQLiteDatabase) migrate() error {
    query := `
    CREATE TABLE IF NOT EXISTS books(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    );
    `

    _, err := r.db.Exec(query)
    return err
}

func ConnectSQLite(dbPath string) (*SQLiteDatabase, error) {
	connection, err := sql.Open("sqlite3", dbPath)
    if err != nil {
		return nil, err
    }

    database := newSQLiteDatabase(connection)

    if err := database.migrate(); err != nil {
		return nil, err
    }

	return database, nil
}
