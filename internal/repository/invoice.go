package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
	"invoice-test/internal/model"
	"time"
)

const insertInvoiceQuery = `
	INSERT INTO invoice (invoice_number, date, customer_name, salesperson, notes, payment_type)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING invoice_number;
`

func (q *Queries) InsertInvoice(ctx context.Context, arg model.Invoice) (string, error) {
	var invoiceNumber string
	err := q.db.QueryRow(ctx, insertInvoiceQuery,
		arg.InvoiceNumber,
		arg.Date,
		arg.CustomerName,
		arg.Salesperson,
		arg.Notes,
		arg.PaymentType,
	).Scan(&invoiceNumber)
	return invoiceNumber, err
}

type PaymentStatus string

const (
	PaymentStatusCASH   PaymentStatus = "CASH"
	PaymentStatusCREDIT PaymentStatus = "CREDIT"
)

const insertProductQuery = `
	INSERT INTO product (id, item_name, quantity, total_cogs, total_price_sold, invoice_number)
VALUES ($1, $2, $3, $4, $5, $6)
`

func (q *Queries) InsertProduct(ctx context.Context, arg model.Product) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, insertProductQuery,
		arg.ID,
		arg.ItemName,
		arg.Quantity,
		arg.TotalCOGS,
		arg.TotalPriceSold,
		arg.InvoiceNumber,
	)
}

const deleteInvoiceQuery = `
	DELETE FROM invoice WHERE invoice_number = $1
`

func (q *Queries) DeleteInvoice(ctx context.Context, id string) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, deleteInvoiceQuery, id)
}

const updateInvoiceQuery = `
UPDATE invoice
SET date = $2, customer_name = $3, salesperson = $4, notes = $5, updated_at = $6, payment_type = $7 
WHERE invoice_number = $1
`

func (q *Queries) UpdateInvoice(ctx context.Context, arg model.Invoice) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, updateInvoiceQuery,
		arg.InvoiceNumber,
		arg.Date,
		arg.CustomerName,
		arg.Salesperson,
		arg.Notes,
		time.Now(),
		arg.PaymentType,
	)
}

const deleteProductQuery = `
DELETE FROM product
WHERE invoice_number = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id string) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, deleteProductQuery, id)
}