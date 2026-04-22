package repository

import (
	"context"
	"time"
)

func (r *Repository) CreateSession(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO sessions (user_id, token_hash, expires_at)
		VALUES (?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt.UTC())
	return err
}

func (r *Repository) GetUserBySessionHash(ctx context.Context, tokenHash string) (User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password_hash, u.created_at, u.updated_at
		FROM sessions s
		INNER JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = ?
		  AND s.expires_at > CURRENT_TIMESTAMP
		LIMIT 1
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, tokenHash).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) DeleteSession(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM sessions WHERE token_hash = ?`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (r *Repository) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at <= CURRENT_TIMESTAMP`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
