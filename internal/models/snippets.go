package models

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	var id int
	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES ($1, $2, now() at time zone 'utc', DATE_ADD(now() at time zone 'utc', $3 * '1 day'::interval)) RETURNING id`
	err := m.DB.QueryRow(context.Background(), query, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > now() at time zone 'utc' and id = $1`
	s := &Snippet{}
	err := m.DB.QueryRow(context.Background(), query, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	query := `select id, title, content, created, expires from snippets
	where expires > now() at time zone 'utc' order by id desc limit 10`
	rows, err := m.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
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
