package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
)

func (s *Storage) GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error) {
	const fn = "storage.postgres.history.GetUserOrderHistory"

	var orders []models.OrderDetail
	// такой селект делаю просто чтобы тупо не запутаться в названиях таблиц и полей а то можно с ума сойти как вьюшку делать
	const q = `SELECT
		o.id 				AS order_id, 
		u.email         	AS email, 
		o.total_price 		AS total_price, 
		o.created_at		AS created_at,
		oi.product_id 		AS product_id,
		p.name 				AS product_name,
		oi.quantity 		AS quantity, 
		t.status 			AS transaction_status
		t.amount 			AS transaction_amount
	FROM orders o 
	JOIN users u ON o.user_id = u.id
	JOIN orders_items oi ON oi.order_id = o.id
	JOIN products p ON oi.product_id = p.id
	JOIN transactions t ON t.order_id = o.id
	WHERE u.email = $1
	ORDER BY o.created_at DESC`

	rows, err := s.db.Query(ctx, q, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	defer rows.Close()

	for rows.Next() {
		var o models.OrderDetail
		if err := rows.Scan(&o.OrderID, &o.Email, &o.TotalPrice, &o.CreatedAt, &o.ProductID, &o.ProductName, &o.Quantity, &o.TxStatus, &o.Amount); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return orders, nil
}

func (s *Storage) GetPopularProducts(ctx context.Context) ([]models.PopularProduct, error) {
	const fn = "storage.postgres.history.GetPopularProducts"

	var products []models.PopularProduct

	const q = `SELECT oi.product_id, p.name, 
	SUM(oi.quantity)::bigint AS total_sold
	FROM orders_items oi
	JOIN products p ON oi.product_id = p.id
	GROUP BY oi.product_id, p.name
	ORDER BY total_sold DESC
	LIMIT 10
	`

	rows, err := s.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.PopularProduct
		if err := rows.Scan(&p.ProductID, &p.Name, &p.TotalSold); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return products, nil

}
