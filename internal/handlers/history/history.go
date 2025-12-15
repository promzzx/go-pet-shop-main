package history

import (
	"context"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type History interface {
	GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error)
	GetPopularProducts(ctx context.Context) ([]models.PopularProduct, error)
}

func HandleUserhistory(log *slog.Logger, history History) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.history.HandleUserhistory"
		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		email := r.URL.Query().Get("email")
		if email == "" {
			log.Error("email parameter is missing")
			http.Error(w, "email parameter is required", http.StatusBadRequest)
			return
		}

		orders, err := history.GetUserOrderHistory(r.Context(), email)
		if err != nil {
			log.Error("failed to get user order history", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("Retrieved user order history successfully", slog.String("url", r.URL.String()))
		render.JSON(w, r, orders)
	}
}

func HandlePopularProducts(log *slog.Logger, history History) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.history.HandlePopularProducts"
		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		products, err := history.GetPopularProducts(r.Context())
		if err != nil {
			log.Error("failed to get popular products", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("Retrieved popular products successfully", slog.String("url", r.URL.String()))
		render.JSON(w, r, products)
	}
}
