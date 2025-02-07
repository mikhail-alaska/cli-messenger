package sqlite
import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS users(
            id INTEGER PRIMARY KEY,
            username TEXT NOT NULL UNIQUE,
            openkey INTEGER NOT NULL);
       `)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}


	stmt3, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS messages(
            id INTEGER PRIMARY KEY,
            username1 TEXT NOT NULL,
            username2 TEXT NOT NULL,
            messagefor1 TEXT
            messagefor2 TEXT);
       `)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	_, err = stmt3.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return &Storage{db: db}, nil
}
