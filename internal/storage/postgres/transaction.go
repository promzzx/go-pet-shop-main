package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
)

func (s *Storage) PlaceOrder(ctx context.Context, email string, items []models.OrderItem) (int, error) {
	const fn = "storage.postgres.transaction.PlaceOrder"

	if email == "" || len(items) == 0 {
		return 0, fmt.Errorf("%s: invalid input data", fn)
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: begin tx: %w", fn, err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// change stock
	const updateStockQuery = `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`
	for _, item := range items {
		res, err := tx.Exec(ctx, updateStockQuery, item.Quantity, item.ProductID)
		if err != nil {
			return 0, fmt.Errorf("%s: update stock: %w", fn, err)
		}
		if res.RowsAffected() == 0 {
			return 0, fmt.Errorf("%s: insufficient stock for product ID %d", fn, item.ProductID)
		}
	}

	// calculate total price
	const getPriceQuery = `SELECT price FROM products WHERE id = $1`

	var totalPrice int64
	for _, item := range items {
		var price int64
		if err := tx.QueryRow(ctx, getPriceQuery, item.ProductID).Scan(&price); err != nil {
			return 0, fmt.Errorf("%s: get price for product %d: %w", fn, item.ProductID, err)
		}
		totalPrice += price * int64(item.Quantity)
	}

	// create order
	const insertOrderQuery = `INSERT INTO orders (user_id, total_price) VALUES ((SELECT id FROM users WHERE email = $1), $2) RETURNING id`

	var orderID int
	if err = tx.QueryRow(ctx, insertOrderQuery, email, totalPrice).Scan(&orderID); err != nil {
		return 0, fmt.Errorf("%s: insert order: %w", fn, err)
	}

	// create order items
	const insertOrderItemQuery = `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`
	for _, item := range items {
		if _, err := tx.Exec(ctx, insertOrderItemQuery, orderID, item.ProductID, item.Quantity); err != nil {
			return 0, fmt.Errorf("%s: insert order item: %w", fn, err)
		}
	}

	// create transaction
	const insertTransactionQuery = `INSERT INTO transactions (order_id, amount, status) VALUES ($1, $2, $3)`

	const status = "pending"

	_, err = tx.Exec(ctx, insertTransactionQuery, orderID, totalPrice, status)
	if err != nil {
		return 0, fmt.Errorf("%s: insert transaction: %w", fn, err)
	}

	// commit full transaction
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: commit: %w", fn, err)
	}

	return orderID, nil
}
