package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/go-financial-planning/backend/internal/repository"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != nil && !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single JSON object")
	}

	return nil
}

func centsFromAmount(amount float64) (int64, error) {
	if amount <= 0 {
		return 0, errors.New("amount must be greater than zero")
	}

	return int64(math.Round(amount * 100)), nil
}

func amountFromCents(cents int64) float64 {
	return float64(cents) / 100
}

func parseDate(raw string) (time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	date, err := time.Parse("2006-01-02", trimmed)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, use YYYY-MM-DD")
	}

	return date, nil
}

func hashSessionToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

type userResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func toUserResponse(user repository.User) userResponse {
	return userResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

type transactionResponse struct {
	ID          int64   `json:"id"`
	Kind        string  `json:"kind"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	AmountCents int64   `json:"amountCents"`
	DueDate     string  `json:"dueDate"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

func toTransactionResponse(transaction repository.Transaction) transactionResponse {
	return transactionResponse{
		ID:          transaction.ID,
		Kind:        transaction.Kind,
		Description: transaction.Description,
		Amount:      amountFromCents(transaction.AmountCents),
		AmountCents: transaction.AmountCents,
		DueDate:     transaction.DueDate.Format("2006-01-02"),
		CreatedAt:   transaction.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   transaction.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
