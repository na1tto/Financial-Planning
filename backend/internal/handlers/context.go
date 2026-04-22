package handlers

import (
	"context"
	"net/http"

	"github.com/go-financial-planning/backend/internal/repository"
)

type contextKey string

const userContextKey contextKey = "current-user"

func withUser(r *http.Request, user repository.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func getUser(r *http.Request) (repository.User, bool) {
	user, ok := r.Context().Value(userContextKey).(repository.User)
	return user, ok
}
