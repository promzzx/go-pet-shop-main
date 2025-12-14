package customer

import (
	"context"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Customers interface {
	CreateUser(ctx context.Context, p models.Customer) error
	GetAllUsers(ctx context.Context) ([]models.Customer, error)
	GetUserByEmail(ctx context.Context, email string) (models.Customer, error)
}

func GetAllUsers(log *slog.Logger, customers Customers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.GetAllUsers"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		users, err := customers.GetAllUsers(r.Context())

		if err != nil {
			log.Error("failed to get users", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Retrieved products successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, users)
	}
}

func CreateUser(log *slog.Logger, customers Customers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.customers.CreateUser"
		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("Creating new user", slog.String("url", r.URL.String()))

		var user models.Customer
		if err := render.DecodeJSON(r.Body, &customers); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := customers.CreateUser(r.Context(), user); err != nil {
			log.Error("failed to create user", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("User created successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "User created successfully"})
	}
}

func GetUserByEmail(log *slog.Logger, customers Customers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.customers.GetUserByEmail"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		email := chi.URLParam(r, "email")
		if email == "" {
			log.Error("empty id")
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		res, err := customers.GetUserByEmail(r.Context(), email)
		if err != nil {
			log.Error("failed to get user by email", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Retrieved user successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, res)
	}
}
