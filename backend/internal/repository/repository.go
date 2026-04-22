package repository

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func New(database *sql.DB) *Repository {
	return &Repository{db: database}
}

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Transaction struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	Kind        string    `json:"kind"`
	Description string    `json:"description"`
	AmountCents int64     `json:"amountCents"`
	DueDate     time.Time `json:"dueDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ForecastAggregate struct {
	Month        string
	IncomeCents  int64
	ExpenseCents int64
}
