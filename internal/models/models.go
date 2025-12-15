package models

import "time"

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int // количество на складе
}

type Customer struct {
	ID    int
	Name  string
	Email string
}

type Order struct {
	ID         int
	CustomerID int
	Totalprice float64
	CreatedAt  time.Time
}

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	Quantity  int
}

type Transaction struct {
	ID              int
	OrderID         int
	Amount          float64
	Status          string
	TransactionDate time.Time
}

type OrderDetail struct {
	OrderID    int
	Email      string
	TotalPrice int64
	CreatedAt  time.Time

	ProductID   int
	ProductName string
	Quantity    int

	TxStatus string
	Amount   int64
}

type PopularProduct struct {
	ProductID int
	Name      string
	TotalSold int64
}
