package order

import (
	"context"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Orders interface {
	CreateOrder(ctx context.Context, order models.Order) (int, error) // Возвращает ID
	AddOrderItem(ctx context.Context, orderItem models.OrderItem) error
	GetOrderByID(ctx context.Context, id int) (models.Order, error)
	GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error)
}

func HandleCreateOrder(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.order.CreateOrder"
		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Creating new order", slog.String("url", r.URL.String()))

		var order models.Order
		if err := render.DecodeJSON(r.Body, &order); err != nil {
			log.Error("failed to decode body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := orders.CreateOrder(r.Context(), order)
		if err != nil {
			log.Error("failed to create order", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("Order created successfully", slog.Int("order_id", id))
		render.JSON(w, r, map[string]int{"order_id": id})
	}
}

func HandleGetOrderByID(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.order.GetOrderByID"
		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		orderID, err := strconv.Atoi(chi.URLParam(r, "id"))

		if err != nil || orderID <= 0 {
			http.Error(w, "invalid order id", http.StatusBadRequest)
			return
		}

		order, err := orders.GetOrderByID(r.Context(), orderID)
		if err != nil {
			log.Error("failed to get order", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("Retrieved order successfully", slog.String("url", r.URL.String()))
		render.JSON(w, r, order)
	}
}

func HandleGetOrdersByUserEmail(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.order.GetOrderByUSEREmail"
		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		email := r.URL.Query().Get("email")
		if email == "" {
			log.Error("email is required")
			http.Error(w, "email is required", http.StatusBadRequest)
			return
		}
		order, err := orders.GetOrdersByUserEmail(r.Context(), email)
		if err != nil {
			log.Error("failed to get order", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("Retrieved order successfully", slog.String("url", r.URL.String()))
		render.JSON(w, r, order)
	}
}

func HandleGetOrderItemsByOrderID(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.order.GetOrderITEMSByOrderID"
		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		orderID, err := strconv.Atoi(chi.URLParam(r, "order_id"))
		if err != nil || orderID <= 0 {
			http.Error(w, "invalid order id", http.StatusBadRequest)
			return
		}
		order, err := orders.GetOrderItemsByOrderID(r.Context(), orderID)
		if err != nil {
			log.Error("failed to get order", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("Retrieved order successfully", slog.String("url", r.URL.String()))
		render.JSON(w, r, order)
	}
}

func HandleAddOrderItem(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.order.AddOrderItem"
		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("Adding order item", slog.String("url", r.URL.String()))

		var orderItem models.OrderItem

		idStr := chi.URLParam(r, "id")
		orderID, err := strconv.Atoi(idStr)
		if err != nil || orderID <= 0 {
			http.Error(w, "invalid order id", http.StatusBadRequest)
			return
		}

		if err := render.DecodeJSON(r.Body, &orderItem); err != nil {
			log.Error("failed to decode body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		orderItem.OrderID = orderID

		if err := orders.AddOrderItem(r.Context(), orderItem); err != nil {
			log.Error("failed to add order item", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("Order item added successfully", slog.String("url", r.URL.String()))
		render.JSON(w, r, map[string]string{"status": "Order item added successfully"})
	}
}
