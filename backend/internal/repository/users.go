package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

func (r *Repository) CreateUser(ctx context.Context, name, email, passwordHash string) (User, error) {
	query := `
		INSERT INTO users (name, email, password_hash, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		RETURNING id, name, email, password_hash, created_at, updated_at
	`

	row := r.db.QueryRowContext(ctx, query, strings.TrimSpace(name), strings.ToLower(strings.TrimSpace(email)), passwordHash)
	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, strings.ToLower(strings.TrimSpace(email))).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID int64) (User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
