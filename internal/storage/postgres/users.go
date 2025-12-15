package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
)

func (s *Storage) CreateUser(ctx context.Context, p models.Customer) error {
	const fn = "storage.postgres.product.CreateUser"
	const q = `INSERT INTO users (name, email) VALUES ($1, $2)`
	_, err := s.db.Exec(ctx, q, p.Name, p.Email)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) GetAllUsers(ctx context.Context) ([]models.Customer, error) {
	const fn = "storage.postgres.product.GetAllUsers"
	rows, err := s.db.Query(ctx, `SELECT * FROM users`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()
	var users []models.Customer
	for rows.Next() {
		var u models.Customer
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.Customer, error) {
	const fn = "storage.postgres.product.GetUserByEmail"

	var u models.Customer
	err := s.db.QueryRow(ctx, `SELECT * FROM users WHERE email = $1`, email).
		Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return models.Customer{}, fmt.Errorf("%s: %w", fn, err)
	}
	return u, nil
}
