package product

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestDeleteProduct_Success(t *testing.T) {
	// ProductsMock — структура, которая подменяет настоящий репозиторий.
	// Мы задаём функцию DeleteProductFunc, которая будет вызвана хендлером.
	// Здесь мы проверяем, что хендлер передаёт правильный id ("42").
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id string) error {
			if id != "42" {
				t.Errorf("expected id 42, got %s", id)
			}
			return nil
		},
	}

	// Создаём новый chi роутер, так как нам нужно протестировать извлечение {id} из URL.
	// Без роутера мы не можем проверить path-параметры.
	r := chi.NewRouter()

	// Регистрируем обработчик DeleteProduct на эндпоинте /products/{id}.
	// Хендлер будет вызван при DELETE /products/{id}
	r.Delete("/products/{id}", DeleteProduct(slog.Default(), mock))

	// Создаем HTTP-запрос DELETE на /products/42.
	// Это симуляция реального запроса без запуска сервера.
	req := httptest.NewRequest(http.MethodDelete, "/products/42", nil)

	// httptest.NewRecorder записывает ответ хендлера.
	// В него попадёт статус-код, заголовки и тело ответа.
	w := httptest.NewRecorder()

	// Передаём запрос в роутер, который найдёт нужный хендлер и выполнит его.
	r.ServeHTTP(w, req)

	// Проверяем, что хендлер вернул статус 200 OK.
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeleteProduct_Fail(t *testing.T) {
	// Мок, который имитирует ошибку при удалении
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id string) error {
			// Проверяем корректность id
			if id != "42" {
				t.Errorf("expected id 42, got %s", id)
			}
			// Возвращаем ошибку — хендлер должен обработать её
			return fmt.Errorf("failed to delete product")
		},
	}

	// Роутер нужен для извлечения {id}
	r := chi.NewRouter()
	r.Delete("/products/{id}", DeleteProduct(slog.Default(), mock))

	// Запрос DELETE /products/42
	req := httptest.NewRequest(http.MethodDelete, "/products/42", nil)
	w := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(w, req)

	// Ожидаем 500 Internal Server Error
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
