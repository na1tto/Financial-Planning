package handlers

import (
	"database/sql"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/go-financial-planning/backend/internal/auth"
	"github.com/go-financial-planning/backend/internal/repository"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	if request.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if _, err := mail.ParseAddress(request.Email); err != nil {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if len(strings.TrimSpace(request.Password)) < 8 {
		writeError(w, http.StatusBadRequest, "password must have at least 8 characters")
		return
	}

	if _, err := h.repo.GetUserByEmail(r.Context(), request.Email); err == nil {
		writeError(w, http.StatusConflict, "email already in use")
		return
	} else if err != nil && !repository.IsNotFound(err) {
		writeError(w, http.StatusInternalServerError, "failed to verify user")
		return
	}

	passwordHash, err := auth.HashPassword(request.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to secure password")
		return
	}

	user, err := h.repo.CreateUser(r.Context(), request.Name, request.Email, passwordHash)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	if err := h.startSession(w, r, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to start session")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"user": toUserResponse(user),
	})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.repo.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		if repository.IsNotFound(err) {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to authenticate")
		return
	}

	valid, err := auth.VerifyPassword(request.Password, user.PasswordHash)
	if err != nil || !valid {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := h.startSession(w, r, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to start session")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user": toUserResponse(user),
	})
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		_ = h.repo.DeleteSession(r.Context(), hashSessionToken(sessionCookie.Value))
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cfg.IsProduction,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	user, ok := getUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user": toUserResponse(user),
	})
}

func (h *Handler) startSession(w http.ResponseWriter, r *http.Request, userID int64) error {
	token, tokenHash, err := auth.GenerateSessionToken()
	if err != nil {
		return err
	}

	if err := h.repo.CreateSession(r.Context(), userID, tokenHash, time.Now().UTC().Add(h.cfg.SessionTTL)); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cfg.IsProduction,
		Expires:  time.Now().Add(h.cfg.SessionTTL),
	})

	return nil
}

func (h *Handler) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		user, err := h.repo.GetUserBySessionHash(r.Context(), hashSessionToken(sessionCookie.Value))
		if err != nil {
			if repository.IsNotFound(err) || err == sql.ErrNoRows {
				writeError(w, http.StatusUnauthorized, "session expired")
				return
			}

			writeError(w, http.StatusInternalServerError, "failed to validate session")
			return
		}

		next.ServeHTTP(w, withUser(r, user))
	})
}
