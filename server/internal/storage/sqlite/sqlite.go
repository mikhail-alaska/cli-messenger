package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"sort"

	"github.com/mattn/go-sqlite3"
	"github.com/mikhail-alaska/cli-messenger/server/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) CreateUser(userName string, openKey int, password string) (int64, error) {
	const op = "storage.sqlite.CreateUser"

	stmt, err := s.db.Prepare("INSERT INTO users(username, openkey, password) VALUES(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	res, err := stmt.Exec(userName, openKey, password)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s, %w", op, storage.ErrUserNameExists)
		}
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s, failed to get last imported id %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetAllUsers() ([]string, error) {
	const op = "storage.sqlite.GetOpenKey"

	var res []string
	stmt, err := s.db.Prepare("SELECT username FROM users")
	if err != nil {
		return res, fmt.Errorf("%s, %w", op, err)
	}

	rows, err := stmt.Query()
	var temp string
	if err != nil {
		return res, fmt.Errorf("%s, while executing statement %w", op, err)
	}
	for rows.Next() {
		rows.Scan(&temp)
		res = append(res, temp)
	}

	return res, nil
}

func (s *Storage) GetOpenKey(username string) (int, error) {
	const op = "storage.sqlite.GetOpenKey"

	stmt, err := s.db.Prepare("SELECT openkey FROM users WHERE username=?")
	if err != nil {
		return -1, fmt.Errorf("%s, %w", op, err)
	}

	var res int
	err = stmt.QueryRow(username).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, storage.ErrUserNameNotFound
	}
	if err != nil {
		return -1, fmt.Errorf("%s, while executing statement %w", op, err)
	}

	return res, nil
}

func (s *Storage) GetPassword(username string) (string, error) {
	const op = "storage.sqlite.GetPassword"

	var res string
	stmt, err := s.db.Prepare("SELECT openkey FROM users WHERE username=?")
	if err != nil {
		return res, fmt.Errorf("%s, %w", op, err)
	}

	err = stmt.QueryRow(username).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		return res, storage.ErrUserNameNotFound
	}
	if err != nil {
		return res, fmt.Errorf("%s, while executing statement %w", op, err)
	}

	return res, nil
}

func (s *Storage) CreateMessage(userName1, userName2, msg1, msg2 string) (int64, error) {
	const op = "storage.sqlite.CreateMessage"

	stmt, err := s.db.Prepare("INSERT INTO messages(username1, username2, messagefor1, messagefor2) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	res, err := stmt.Exec(userName1, userName2, msg1, msg2)
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s, failed to get last imported id %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetChatsByUsername(username string) ([]string, error) {
	const op = "storage.sqlite.GetAliasByUrl"
	var usernames []string

	stmt, err := s.db.Prepare("SELECT username2 FROM messages WHERE username1=?")
	if err != nil {
		return usernames, fmt.Errorf("%s, %w", op, err)
	}

	var res string
	rows, err := stmt.Query(username)
	if err != nil {
		return usernames, fmt.Errorf("%s, while executing statement %w", op, err)
	}
	for rows.Next() {
		rows.Scan(&res)
		if !slices.Contains(usernames, res) {
			usernames = append(usernames, res)
		}
	}

	stmt, err = s.db.Prepare("SELECT username1 FROM messages WHERE username2=?")
	if err != nil {
		return usernames, fmt.Errorf("%s, %w", op, err)
	}

	rows, err = stmt.Query(username)
	if err != nil {
		return usernames, fmt.Errorf("%s, while executing statement %w", op, err)
	}
	for rows.Next() {
		rows.Scan(&res)
		if !slices.Contains(usernames, res) {
			usernames = append(usernames, res)
		}
	}
	return usernames, nil
}

func (s *Storage) GetMessages(username1, username2 string) ([]storage.StorageMessages, error) {
	const op = "storage.sqlite.GetAliasByUrl"
	var messages []storage.StorageMessages
	var res storage.StorageMessages

	stmt, err := s.db.Prepare("SELECT id, messagefor1 FROM messages WHERE username1=? AND username2=?")
	if err != nil {
		return messages, fmt.Errorf("%s, %w", op, err)
	}

	rows1, err := stmt.Query(username1, username2)
	if err != nil {
		return messages, fmt.Errorf("%s, while executing statement %w", op, err)
	}
	for rows1.Next() {
		rows1.Scan(&res.Id, &res.Message)
		messages = append(messages, res)
	}

	rows2, err := stmt.Query(username2, username1)
	if err != nil {
		return messages, fmt.Errorf("%s, while executing statement %w", op, err)
	}
	for rows2.Next() {
		rows2.Scan(&res.Id, &res.Message)
		messages = append(messages, res)
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Id < messages[j].Id
	})

	return messages, nil
}
