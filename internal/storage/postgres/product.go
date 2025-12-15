package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
)

// product.go реализует методы работы с продуктами в PostgreSQL.
func (s *Storage) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	const fn = "storage.postgres.product.GetAllProducts"

	rows, err := s.db.Query(ctx, `SELECT * FROM products`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *Storage) CreateProduct(ctx context.Context, p models.Product) error {
	const fn = "storage.postgres.product.CreateProduct"

	_, err := s.db.Exec(ctx,
		`INSERT INTO products (name, price, stock) VALUES ($1, $2, $3)`,
		p.Name, p.Price, p.Stock)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) DeleteProduct(ctx context.Context, id string) error {
	const fn = "storage.postgres.product.DeleteProduct"

	_, err := s.db.Exec(ctx,
		`DELETE FROM products WHERE id = $1`,
		id)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) UpdateProduct(ctx context.Context, p models.Product) error {
	const fn = "storage.postgres.product.UpdateProduct"

	_, err := s.db.Exec(ctx,
		`UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4`,
		p.Name, p.Price, p.Stock, p.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) GetProductByID(ctx context.Context, id string) (models.Product, error) {
	const fn = "storage.postgres.product.GetProductByID"

	const q = `SELECT id, name, price, stock FROM products WHERE id = $1`

	var p models.Product
	err := s.db.QueryRow(ctx, q, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		return models.Product{}, fmt.Errorf("%s: %w", fn, err)
	}
	return p, nil
}
