package transaction

import (
	"context"
	"encoding/json"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
)

type Transactions interface {
	PlaceOrder(ctx context.Context, email string, items []models.OrderItem) (int, error)
}

// struct for incoming request
type placeOrderRequest struct {
	Email string `json:"email"`
	Items []struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	} `json:"items"`
}

// struct for successful rsponse
type placeOrderResponse struct {
	OrderID int `json:"order_id"`
}

// struct for err message
type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON is a helper to wite JSON resp
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func HandlePlaceOrder(log *slog.Logger, transactions Transactions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.transaction.HandlePlaceOrder"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req placeOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
			return
		}

		if req.Email == "" || len(req.Items) == 0 {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "email and items are required"})
			return
		}

		items := make([]models.OrderItem, 0, len(req.Items))
		for i, item := range req.Items {
			items = append(items, models.OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			})
			log.Debug("order item parsed", slog.Int("index", i), slog.Int("product_id", item.ProductID), slog.Int("quantity", item.Quantity))
		}

		orderID, err := transactions.PlaceOrder(r.Context(), req.Email, items)
		if err != nil {
			log.Error("failed transaction create", slog.Any("error", err))
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to place order"})
			return
		}

		log.Info("order placed successfully", slog.Int("order_id", orderID))
		writeJSON(w, http.StatusOK, placeOrderResponse{OrderID: orderID})
	}
}
