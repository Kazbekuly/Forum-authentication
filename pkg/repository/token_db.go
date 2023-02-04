package repository

import (
	"context"
	"fmt"
	"time"

	"Forum"
)

func (r *AuthSQL) AddToken(token Forum.Token) (Forum.Token, error) {
	if err := r.DeleteTokenByUserID(token.UserId); err != nil {
		return token, err
	}
	records := fmt.Sprintf("INSERT INTO %s (UserId, auth_token, expires_at) values ($1, $2, $3)", tokenTable)
	query, err := r.db.Prepare(records)
	if err != nil {
		return token, err
	}
	_, err = query.Exec(token.UserId, token.AuthToken, token.ExpiresAT)
	if err != nil {
		return token, err
	}
	return token, nil
}

func (r *AuthSQL) GetToken(token string) (Forum.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30*time.Second))
	defer cancel()
	var userToken Forum.Token
	rows := r.db.QueryRowContext(ctx, "SELECT * FROM authorization_token WHERE auth_token=$1", token)
	if err := rows.Scan(&userToken.Id, &userToken.UserId, &userToken.AuthToken, &userToken.ExpiresAT); err != nil {
		return userToken, err
	}
	return userToken, nil
}

func (r *AuthSQL) GetUserByToken(token string) (Forum.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var userToken Forum.Token
	var user Forum.User
	rows := r.db.QueryRowContext(ctx, "SELECT * FROM authorization_token WHERE auth_token=$1", token)
	if err := rows.Scan(&userToken.Id, &userToken.UserId, &userToken.AuthToken, &userToken.ExpiresAT); err != nil {
		return user, err
	}
	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id=$1", userToken.UserId)
	if err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password); err != nil {
		return user, err
	}
	return user, nil
}

func (r *AuthSQL) DeleteToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15*time.Second))
	defer cancel()
	_, err := r.db.ExecContext(ctx, "DELETE FROM authorization_token WHERE auth_token=$1", token)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthSQL) DeleteTokenByUserID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15*time.Second))
	defer cancel()
	_, err := r.db.ExecContext(ctx, "DELETE FROM authorization_token WHERE UserId=$1", id)
	if err != nil {
		return err
	}
	return nil
}
