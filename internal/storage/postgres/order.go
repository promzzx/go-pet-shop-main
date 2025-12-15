package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
)

// order.go реализует методы работы с заказами в PostgreSQL.
func (s *Storage) CreateOrder(ctx context.Context, o models.Order) (int, error) {
	const fn = "storage.postgres.order.CreateOrder"
	const q = `INSERT INTO orders (user_id, total_price) VALUES ($1, $2)RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, q, o.CustomerID, o.Totalprice).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) AddOrderItem(ctx context.Context, orderItem models.OrderItem) error {
	const fn = "storage.postgress.order.AddOrderItem"

	const q = `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`

	_, err := s.db.Exec(ctx, q, orderItem)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	const fn = "storage.postgres.order.GetOrderByID"

	const q = `SELECT id, user_id, total_price, created_at FROM orders WHERE id = $1`

	var o models.Order
	err := s.db.QueryRow(ctx, q, id).Scan(&o.ID, &o.CustomerID, &o.Totalprice, &o.CreatedAt)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", fn, err)
	}
	return o, nil
}

func (s *Storage) GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error) {
	const fn = "storage.postgres.order.GetOrdersByUserEmail"

	const q = `SELECT o.id, o.user_id, o.total_price, o.created_at FROM orders o JOIN users u ON o.user_id = u.id WHERE u.email = $1`

	rows, err := s.db.Query(ctx, q, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.Totalprice, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return orders, nil
}

func (s *Storage) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	const fn = "storage.postgres.Order.GetOrderItemsByOrderID"

	const q = "SELECT * FROM orders_items WHERE orders_id=$1"

	rows, err := s.db.Query(ctx, q, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var oi []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		oi = append(oi, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return oi, nil
}
