package product

import (
	"bytes"
	"context"
	"errors"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateProduct_Success(t *testing.T) {
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, p models.Product) error {
			if p.ID != 1 {
				t.Errorf("expected ID=1, got %d", p.ID)
			}
			if p.Name != "Cat Food" {
				t.Errorf("expected Name=Cat Food, got %q", p.Name)
			}
			return nil
		},
	}

	body := []byte(`{"id":1,"name":"Cat Food"}`)
	req := httptest.NewRequest(http.MethodPut, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := UpdateProduct(slog.Default(), mock)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}
func TestUpdateProduct_Fail(t *testing.T) {
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, p models.Product) error {
			return errors.New("DB error")
		},
	}

	body := []byte(`{"id":1,"name":"Cat Food"}`)
	req := httptest.NewRequest(http.MethodPut, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := UpdateProduct(slog.Default(), mock)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body=%s", w.Code, w.Body.String())
	}
}
