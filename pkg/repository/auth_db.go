package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"Forum"
)

type AuthSQL struct {
	db *sql.DB
}

func NewAuthSQL(db *sql.DB) *AuthSQL {
	return &AuthSQL{db: db}
}

func (r *AuthSQL) CreateUser(user Forum.User) error {
	records := fmt.Sprintf("INSERT INTO users (Email, Username, Password) values ($1, $2, $3)")
	query, err := r.db.Prepare(records)
	if err != nil {
		return err
	}
	_, err = query.Exec(user.Email, user.Username, user.Password)
	if err != nil {
		return errors.New("user is already exists")
	}
	return nil
}

func (r *AuthSQL) GetUser(username, email string) (Forum.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var user Forum.User
	request := fmt.Sprintf("SELECT * FROM users WHERE Username = $1 or Email=$2")
	rows := r.db.QueryRowContext(ctx, request, username, email)
	if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		return user, err
	}
	return user, nil
}

func (r *AuthSQL) CheckInvalid(username, email string) (Forum.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var user Forum.User
	if err := r.db.QueryRowContext(ctx, "SELECT id, email, username FROM users WHERE Email = $1 AND Username=$2 ", email, username).Scan(&user.Id, &user.Email, &user.Username); err != nil {
		return user, err
	}
	return user, nil
}
