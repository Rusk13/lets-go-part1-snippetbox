package models

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*User, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := "INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3,  now() at time zone 'utc')"
	_, err = m.DB.Exec(context.Background(), query, name, email, string(hashedPassword))
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" && strings.Contains(pgxError.Message, "users_email_key") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	query := "SELECT id, hashed_password FROM users WHERE email = $1"
	err := m.DB.QueryRow(context.Background(), query, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT true FROM users WHERE id = $1)`

	err := m.DB.QueryRow(context.Background(), query, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) Get(id int) (*User, error) {
	var u User
	query := `SELECT * FROM users where id = $1`
	err := m.DB.QueryRow(context.Background(), query, id).Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.Created)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &User{}, ErrNoRecord
		} else {
			return &User{}, err
		}
	}
	u.HashedPassword = []byte{}
	return &u, nil

}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	var u User
	query := `SELECT * FROM users WHERE id=$1`
	err := m.DB.QueryRow(context.Background(), query, id).Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.Created)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}

	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	query = `UPDATE users SET hashed_password = $1 WHERE id=$2 `
	m.DB.Exec(context.Background(), query, newHashedPassword, id)
	return nil
}
