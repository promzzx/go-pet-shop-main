package product

import (
	"context"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Products interface {
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	CreateProduct(ctx context.Context, product models.Product) error
	DeleteProduct(ctx context.Context, id string) error
	UpdateProduct(ctx context.Context, product models.Product) error
	GetProductByID(ctx context.Context, id string) (models.Product, error)
}

func GetAllProducts(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.GetAllProducts"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		products, err := products.GetAllProducts(r.Context())

		if err != nil {
			log.Error("failed to get products", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Retrieved products successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, products)
	}
}

func CreateProduct(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.CreateProduct"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Creating new product", slog.String("url", r.URL.String()))

		var product models.Product
		if err := render.DecodeJSON(r.Body, &product); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := products.CreateProduct(r.Context(), product); err != nil {
			log.Error("failed to create product", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Product created successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "Product created successfully"})
	}
}

func DeleteProduct(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.DeleteProduct"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Deleting product", slog.String("url", r.URL.String()))

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Error("empty id")
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if err := products.DeleteProduct(r.Context(), id); err != nil {
			log.Error("failed to delete product", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Deleted product successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "Product deleted successfully"})
	}
}

func UpdateProduct(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.UpdateProduct"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Updating product", slog.String("url", r.URL.String()))

		var product models.Product
		if err := render.DecodeJSON(r.Body, &product); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := products.UpdateProduct(r.Context(), product); err != nil {
			log.Error("failed to update product", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Product updated successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "Product updated successfully"})
	}
}

func GetProductByID(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.GetProductByID"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Error("empty id")
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		res, err := products.GetProductByID(r.Context(), id)
		if err != nil {
			log.Error("failed to get product by ID", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Retrieved products successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, res)
	}
}
