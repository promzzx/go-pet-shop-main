package product

import "go-pet-shop/internal/models"

// ProductsMock — тестовая структура, которая имитирует (мокает) репозиторий продуктов.
// Она копирует его методы, но вместо настоящей логики хранит функции,
// которые можно задавать прямо в тестах.
//
// Зачем это нужно?
// — Чтобы изолировать тесты хендлеров от настоящей базы данных.
// — Чтобы проверять только работу HTTP-слоя, а не всей системы.
// — Чтобы гибко задавать поведение (успех, ошибка, проверка аргументов) прямо внутри теста.
type ProductsMock struct {
	// Функции-заглушки, которые тесты могут переопределять.
	// Каждая из них полностью заменяет соответствующий метод интерфейса.
	GetAllProductsFunc func() ([]models.Product, error)
	CreateProductFunc  func(product models.Product) error
	DeleteProductFunc  func(id string) error
	UpdateProductFunc  func(product models.Product) error
}

// Каждый из методов ниже просто вызывает соответствующую функцию,
// которую тест заполняет заранее.
//
// Почему так?
// — Мы сохраняем “поведение интерфейса” (методы остаются такими же).
// — Но саму логику тест может менять как угодно.

// Мок-версия метода GetAllProducts.
// Вместо выполнения запроса в БД вызовет заранее заданную функцию.
func (m *ProductsMock) GetAllProducts() ([]models.Product, error) {
	return m.GetAllProductsFunc()
}

// Мок метода CreateProduct — полностью контролируется тестом.
func (m *ProductsMock) CreateProduct(p models.Product) error {
	return m.CreateProductFunc(p)
}

// Мок метода DeleteProduct — тест может заставить его:
// ✓ вернуть успех
// ✓ вернуть ошибку
// ✓ проверить правильность аргументов (id)
func (m *ProductsMock) DeleteProduct(id string) error {
	return m.DeleteProductFunc(id)
}

// Мок метода UpdateProduct — аналогично, поведение задаётся тестом.
func (m *ProductsMock) UpdateProduct(p models.Product) error {
	return m.UpdateProductFunc(p)
}
