package model

import "github.com/google/uuid"

type Product struct {
	ID             uuid.UUID
	ItemName       string
	Quantity       int
	TotalCOGS      int64
	TotalPriceSold int64
	InvoiceNumber  string
}
