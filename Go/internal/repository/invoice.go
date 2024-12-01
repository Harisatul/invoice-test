package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	model2 "invoice-test/internal/model"
	"time"
)

const insertInvoiceQuery = `
	INSERT INTO invoice (invoice_number, date, customer_name, salesperson, notes, payment_type)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING invoice_number;
`

func (q *Queries) InsertInvoice(ctx context.Context, arg model2.Invoice) (string, error) {
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

func (q *Queries) InsertProduct(ctx context.Context, arg model2.Product) (pgconn.CommandTag, error) {
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

func (q *Queries) UpdateInvoice(ctx context.Context, arg model2.Invoice) (pgconn.CommandTag, error) {
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

const countAllInvoiceWithGivenDateQuery = `
SELECT COUNT(*) FROM invoice WHERE date BETWEEN $1 AND $2;
`

func (q *Queries) CountAllInvoiceWithGivenDate(ctx context.Context, startDate time.Time, endDate time.Time) (int64, error) {
	var count int64
	return count, q.db.QueryRow(ctx, countAllInvoiceWithGivenDateQuery, startDate, endDate).Scan(&count)
}

const countAllInvoiceWithGivenDateAndCashOnlyQuery = `
SELECT COUNT(*) FROM invoice WHERE payment_type= 'CASH' AND date BETWEEN $1 AND $2;
`

func (q *Queries) CountAllInvoiceWithGivenDateAndCashOnly(ctx context.Context, startDate time.Time, endDate time.Time) (int64, error) {
	var count int64
	return count, q.db.QueryRow(ctx, countAllInvoiceWithGivenDateAndCashOnlyQuery, startDate, endDate).Scan(&count)
}

const getAllInvoiceWithPagination = `
SELECT invoice_number, date, customer_name, salesperson, notes, payment_type
FROM invoice WHERE date BETWEEN $1 AND $2
LIMIT $3 OFFSET $4;	
`

func (q *Queries) GetAllInvoiceWithGivenDate(ctx context.Context, startDate time.Time, endDate time.Time, page, size int) ([]model2.Invoice, error) {

	offset := (page - 1) * size

	rows, err := q.db.Query(ctx, getAllInvoiceWithPagination, startDate, endDate, size, offset)
	if err != nil {
		return nil, fmt.Errorf("error when querying: %w", err)
	}
	defer rows.Close()

	var invoices []model2.Invoice
	for rows.Next() {
		var invoice model2.Invoice
		err := rows.Scan(
			&invoice.InvoiceNumber,
			&invoice.Date,
			&invoice.CustomerName,
			&invoice.Salesperson,
			&invoice.Notes,
			&invoice.PaymentType,
		)
		if err != nil {
			return nil, fmt.Errorf("error when scanning: %w", err)
		}
		invoices = append(invoices, invoice)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error when scanning: %w", err)
	}
	return invoices, nil
}

const getSumofAllInvoice = `
SELECT invoice.invoice_number, p.total_cogs, p.total_price_sold
FROM invoice
JOIN public.product p on invoice.invoice_number = p.invoice_number
WHERE date BETWEEN $1 AND $2;
`

func (q *Queries) GetSumofAllInvoice(ctx context.Context, startDate time.Time, endDate time.Time) ([]model2.Product, error) {
	rows, err := q.db.Query(ctx, getSumofAllInvoice, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var invoices []model2.Product
	for rows.Next() {
		var inv model2.Product
		if err := rows.Scan(&inv.InvoiceNumber, &inv.TotalCOGS, &inv.TotalPriceSold); err != nil {
			return nil, fmt.Errorf("failed to scan invoice row: %w", err)
		}
		invoices = append(invoices, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating invoice rows: %w", err)
	}
	return invoices, nil
}
