package service

import (
	"context"
	"errors"
	"invoice-test/internal/model"
	"invoice-test/pkg"
	"log/slog"
	"time"
)

func (s Service) CreateInvoice(ctx context.Context, request model.CreateInvoiceRequest) (string, error) {

	status, err := pkg.ToPaymentStatus(request.PaymentType)
	if err != nil {
		return "", err
	}
	invoice := model.Invoice{
		InvoiceNumber: pkg.InvoiceNumberGenerator(),
		Date:          time.Now(),
		CustomerName:  request.CustomerName,
		Salesperson:   request.SalesPersonName,
		Notes:         request.Notes,
		PaymentType:   status,
	}

	tx, err := s.Db.Begin(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to begin transaction : ", err)
		return "", err
	}
	defer tx.Rollback(ctx)

	qTx := s.Querier.WithTx(tx)
	insertInvoice, err := qTx.InsertInvoice(ctx, invoice)
	if err != nil {
		slog.ErrorContext(ctx, "failed to insert invoice : ", err)
	}
	for _, product := range request.Product {
		id, err := pkg.GenerateId()
		if err != nil {
			return "", err
		}
		productModel := model.Product{
			ID:             id,
			ItemName:       product.Name,
			Quantity:       product.Quantity,
			TotalCOGS:      product.TotalCostOfGoodSold,
			TotalPriceSold: product.TotalPriceSold,
			InvoiceNumber:  insertInvoice,
		}
		_, err = qTx.InsertProduct(ctx, productModel)
		if err != nil {
			slog.ErrorContext(ctx, "failed to insert product : ", err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to commit transaction : ", err)
	}
	return insertInvoice, nil
}

func (s Service) DeleteInvoice(ctx context.Context, id string) error {
	row, err := s.Querier.DeleteInvoice(ctx, id)
	if row.RowsAffected() == 0 {
		slog.WarnContext(ctx, "failed to delete Invoice : ", err)
		return errors.New("given id not found")
	}
	return nil
}
