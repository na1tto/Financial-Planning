package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-financial-planning/backend/internal/repository"
)

func (h *Handler) createTransaction(w http.ResponseWriter, r *http.Request) {
	user, ok := getUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	input, ok := h.parseTransactionInput(w, r)
	if !ok {
		return
	}

	transaction, err := h.repo.CreateTransaction(r.Context(), user.ID, input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create transaction")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"transaction": toTransactionResponse(transaction),
	})
}

func (h *Handler) listTransactions(w http.ResponseWriter, r *http.Request) {
	user, ok := getUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	transactions, err := h.repo.ListTransactions(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list transactions")
		return
	}

	response := make([]transactionResponse, 0, len(transactions))
	for _, item := range transactions {
		response = append(response, toTransactionResponse(item))
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"transactions": response,
	})
}

func (h *Handler) updateTransaction(w http.ResponseWriter, r *http.Request) {
	user, ok := getUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	transactionID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || transactionID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	input, ok := h.parseTransactionInput(w, r)
	if !ok {
		return
	}

	transaction, err := h.repo.UpdateTransaction(r.Context(), user.ID, transactionID, input)
	if err != nil {
		if repository.IsNotFound(err) {
			writeError(w, http.StatusNotFound, "transaction not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to update transaction")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"transaction": toTransactionResponse(transaction),
	})
}

func (h *Handler) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	user, ok := getUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	transactionID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || transactionID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	if err := h.repo.DeleteTransaction(r.Context(), user.ID, transactionID); err != nil {
		if repository.IsNotFound(err) {
			writeError(w, http.StatusNotFound, "transaction not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to delete transaction")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) forecastMonthly(w http.ResponseWriter, r *http.Request) {
	user, ok := getUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	months := 12
	if raw := strings.TrimSpace(r.URL.Query().Get("months")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < 1 || value > 36 {
			writeError(w, http.StatusBadRequest, "months must be between 1 and 36")
			return
		}
		months = value
	}

	start := time.Now().UTC()
	start = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, months, 0)

	aggregates, err := h.repo.GetMonthlyForecast(r.Context(), user.ID, start, end)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to build forecast")
		return
	}

	aggregateMap := map[string]repository.ForecastAggregate{}
	for _, item := range aggregates {
		aggregateMap[item.Month] = item
	}

	type monthlyPoint struct {
		Month   string  `json:"month"`
		Income  float64 `json:"income"`
		Expense float64 `json:"expense"`
		Net     float64 `json:"net"`
	}

	result := make([]monthlyPoint, 0, months)
	for i := 0; i < months; i++ {
		current := start.AddDate(0, i, 0)
		monthKey := current.Format("2006-01")
		aggregate := aggregateMap[monthKey]
		income := amountFromCents(aggregate.IncomeCents)
		expense := amountFromCents(aggregate.ExpenseCents)

		result = append(result, monthlyPoint{
			Month:   monthKey,
			Income:  income,
			Expense: expense,
			Net:     income - expense,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"months": months,
		"data":   result,
	})
}

func (h *Handler) parseTransactionInput(w http.ResponseWriter, r *http.Request) (repository.TransactionInput, bool) {
	var request struct {
		Kind        string  `json:"kind"`
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		DueDate     string  `json:"dueDate"`
	}

	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return repository.TransactionInput{}, false
	}

	request.Kind = strings.TrimSpace(strings.ToLower(request.Kind))
	request.Description = strings.TrimSpace(request.Description)

	if request.Kind != "income" && request.Kind != "expense" {
		writeError(w, http.StatusBadRequest, "kind must be income or expense")
		return repository.TransactionInput{}, false
	}

	if len(request.Description) < 3 {
		writeError(w, http.StatusBadRequest, "description must have at least 3 characters")
		return repository.TransactionInput{}, false
	}

	amountCents, err := centsFromAmount(request.Amount)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return repository.TransactionInput{}, false
	}

	dueDate, err := parseDate(request.DueDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return repository.TransactionInput{}, false
	}

	return repository.TransactionInput{
		Kind:        request.Kind,
		Description: request.Description,
		AmountCents: amountCents,
		DueDate:     dueDate,
	}, true
}
