package models

import (
	"database/sql"
	"errors"
	"time"
)

const (
	// ExtensionWindow is how close to expiry a snippet must be before its
	// expiry can be extended.
	ExtensionWindow = 72 * time.Hour
	// ExtensionPeriod is how much extra life an extension grants.
	ExtensionPeriod = 7 * 24 * time.Hour
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (Snippet, error)
	Latest() ([]Snippet, error)
	ExtendExpiry(id int) (time.Time, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Snippet{}, ErrNoRecord
		default:
			return Snippet{}, err
		}
	}

	return s, nil
}

// ExtendExpiry pushes the expiry of a live snippet a further 7 days into the
// future, but only if the snippet is already within 3 days of expiring. If the
// snippet doesn't exist, has already expired, or is still too far from expiry,
// ErrNotExtendable is returned and nothing is changed.
func (m *SnippetModel) ExtendExpiry(id int) (time.Time, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return time.Time{}, err
	}
	defer tx.Rollback()

	updateStmt := `UPDATE snippets SET expires = DATE_ADD(expires, INTERVAL 7 DAY)
	WHERE id = ? AND expires > UTC_TIMESTAMP()
	AND expires <= DATE_ADD(UTC_TIMESTAMP(), INTERVAL 3 DAY)`

	result, err := tx.Exec(updateStmt, id)
	if err != nil {
		return time.Time{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return time.Time{}, err
	}
	if rows == 0 {
		return time.Time{}, ErrNotExtendable
	}

	selectStmt := `SELECT expires FROM snippets WHERE id = ?`

	var expires time.Time

	err = tx.QueryRow(selectStmt, id).Scan(&expires)
	if err != nil {
		return time.Time{}, err
	}

	err = tx.Commit()
	if err != nil {
		return time.Time{}, err
	}

	return expires, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
