package product

import (
	"context"
	"go-pet-shop/internal/models"
)

// ProductsMock — тестовая структура, которая имитирует (мокает) репозиторий продуктов.
// Она копирует его методы, но вместо настоящей логики хранит функции,
// которые можно задавать прямо в тестах.
//
// Зачем это нужно?
// — Чтобы изолировать тесты хендлеров от настоящей базы данных.
// — Чтобы проверять только работу HTTP-слоя, а не всей системы.
// — Чтобы гибко задавать поведение (успех, ошибка, проверка аргументов) прямо внутри теста.
type ProductsMock struct {
	GetAllProductsFunc func(ctx context.Context) ([]models.Product, error)
	CreateProductFunc  func(ctx context.Context, product models.Product) error
	DeleteProductFunc  func(ctx context.Context, id string) error
	UpdateProductFunc  func(ctx context.Context, product models.Product) error
	GetProductByIDFunc func(ctx context.Context, id string) (models.Product, error)
}

// Каждый из методов ниже просто вызывает соответствующую функцию,
// которую тест заполняет заранее.
//
// Почему так?
// — Мы сохраняем “поведение интерфейса” (методы остаются такими же).
// — Но саму логику тест может менять как угодно.

// Мок-версия метода GetAllProducts.
// Вместо выполнения запроса в БД вызовет заранее заданную функцию.
func (m *ProductsMock) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return m.GetAllProductsFunc(ctx)
}

// Мок метода CreateProduct — полностью контролируется тестом.
func (m *ProductsMock) CreateProduct(ctx context.Context, p models.Product) error {
	return m.CreateProductFunc(ctx, p)
}

// Мок метода DeleteProduct — тест может заставить его:
// ✓ вернуть успех
// ✓ вернуть ошибку
// ✓ проверить правильность аргументов (id)
func (m *ProductsMock) DeleteProduct(ctx context.Context, id string) error {
	return m.DeleteProductFunc(ctx, id)
}

// Мок метода UpdateProduct — аналогично, поведение задаётся тестом.
func (m *ProductsMock) UpdateProduct(ctx context.Context, p models.Product) error {
	return m.UpdateProductFunc(ctx, p)
}

func (m *ProductsMock) GetProductByID(ctx context.Context, id string) (models.Product, error) {
	return m.GetProductByIDFunc(ctx, id)
}
