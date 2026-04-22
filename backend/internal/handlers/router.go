package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-financial-planning/backend/internal/app"
	"github.com/go-financial-planning/backend/internal/repository"
)

const sessionCookieName = "finance_session"

type Handler struct {
	cfg  app.Config
	repo *repository.Repository
}

func New(cfg app.Config, repo *repository.Repository) *Handler {
	return &Handler{cfg: cfg, repo: repo}
}

func (h *Handler) Router() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(15 * time.Second))
	router.Use(middleware.CleanPath)
	router.Use(h.securityHeaders)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{h.cfg.FrontendOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(authRouter chi.Router) {
			authRouter.Post("/register", h.register)
			authRouter.Post("/login", h.login)
			authRouter.Post("/logout", h.logout)
			authRouter.With(h.requireAuth).Get("/me", h.me)
		})

		r.Group(func(privateRouter chi.Router) {
			privateRouter.Use(h.requireAuth)
			privateRouter.Get("/transactions", h.listTransactions)
			privateRouter.Post("/transactions", h.createTransaction)
			privateRouter.Put("/transactions/{id}", h.updateTransaction)
			privateRouter.Delete("/transactions/{id}", h.deleteTransaction)
			privateRouter.Get("/forecast/monthly", h.forecastMonthly)
		})
	})

	return router
}

func (h *Handler) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none';")

		next.ServeHTTP(w, r)
	})
}
