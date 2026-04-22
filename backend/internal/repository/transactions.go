package repository

import (
	"context"
	"database/sql"
	"time"
)

type TransactionInput struct {
	Kind        string
	Description string
	AmountCents int64
	DueDate     time.Time
}

func (r *Repository) CreateTransaction(ctx context.Context, userID int64, input TransactionInput) (Transaction, error) {
	query := `
		INSERT INTO transactions (user_id, kind, description, amount_cents, due_date, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		RETURNING id, user_id, kind, description, amount_cents, due_date, created_at, updated_at
	`

	row := r.db.QueryRowContext(ctx, query, userID, input.Kind, input.Description, input.AmountCents, input.DueDate.Format("2006-01-02"))
	var transaction Transaction
	if err := row.Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Kind,
		&transaction.Description,
		&transaction.AmountCents,
		&transaction.DueDate,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	); err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

func (r *Repository) ListTransactions(ctx context.Context, userID int64) ([]Transaction, error) {
	query := `
		SELECT id, user_id, kind, description, amount_cents, due_date, created_at, updated_at
		FROM transactions
		WHERE user_id = ?
		ORDER BY due_date ASC, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]Transaction, 0)
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Kind,
			&transaction.Description,
			&transaction.AmountCents,
			&transaction.DueDate,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		); err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, rows.Err()
}

func (r *Repository) UpdateTransaction(ctx context.Context, userID, transactionID int64, input TransactionInput) (Transaction, error) {
	query := `
		UPDATE transactions
		SET kind = ?, description = ?, amount_cents = ?, due_date = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
		RETURNING id, user_id, kind, description, amount_cents, due_date, created_at, updated_at
	`

	row := r.db.QueryRowContext(ctx, query, input.Kind, input.Description, input.AmountCents, input.DueDate.Format("2006-01-02"), transactionID, userID)
	var transaction Transaction
	if err := row.Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Kind,
		&transaction.Description,
		&transaction.AmountCents,
		&transaction.DueDate,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	); err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

func (r *Repository) DeleteTransaction(ctx context.Context, userID, transactionID int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM transactions WHERE id = ? AND user_id = ?`, transactionID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) GetMonthlyForecast(ctx context.Context, userID int64, startDate, endDate time.Time) ([]ForecastAggregate, error) {
	query := `
		SELECT strftime('%Y-%m', due_date) AS month,
		       SUM(CASE WHEN kind = 'income' THEN amount_cents ELSE 0 END) AS income_cents,
		       SUM(CASE WHEN kind = 'expense' THEN amount_cents ELSE 0 END) AS expense_cents
		FROM transactions
		WHERE user_id = ?
		  AND due_date >= date(?)
		  AND due_date < date(?)
		GROUP BY month
		ORDER BY month ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aggregates := make([]ForecastAggregate, 0)
	for rows.Next() {
		var item ForecastAggregate
		if err := rows.Scan(&item.Month, &item.IncomeCents, &item.ExpenseCents); err != nil {
			return nil, err
		}
		aggregates = append(aggregates, item)
	}

	return aggregates, rows.Err()
}
