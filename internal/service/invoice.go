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

func (s Service) UpdateInvoice(ctx context.Context, arg model.UpdateInvoiceRequest, params string) (string, error) {

	tx, err := s.Db.Begin(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to begin transaction : ", err)
		return "", err
	}
	defer tx.Rollback(ctx)

	status, err := pkg.ToPaymentStatus(arg.PaymentType)
	if err != nil {
		return "", err
	}

	invoice := model.Invoice{
		InvoiceNumber: params,
		Date:          time.Now(),
		CustomerName:  arg.CustomerName,
		Salesperson:   arg.SalesPersonName,
		Notes:         arg.Notes,
		PaymentType:   status,
	}

	qTx := s.Querier.WithTx(tx)
	insertInvoice, err := qTx.UpdateInvoice(ctx, invoice)
	if err != nil {
		slog.WarnContext(ctx, "failed to update invoice : ", err)
	}
	if insertInvoice.RowsAffected() == 0 {
		return "", errors.New("given id not found")
	}
	_, err = qTx.DeleteProduct(ctx, params)
	if err != nil {
		slog.WarnContext(ctx, "failed to update invoice : ", err)
		return "", errors.New("failed to update invoice")
	}
	for _, product := range arg.Product {
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
			InvoiceNumber:  invoice.InvoiceNumber,
		}
		_, err = qTx.InsertProduct(ctx, productModel)
		if err != nil {
			slog.ErrorContext(ctx, "failed to insert product : ", err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to commit transaction : ", err)
	}
	return invoice.InvoiceNumber, nil
}

func (s Service) GetAllInvoice(ctx context.Context, startTime time.Time, endTime time.Time, page, size int) (model.GetInvoiceResponse, model.PaginationIndex, error) {
	count, err := s.Querier.CountAllInvoiceWithGivenDate(ctx, startTime, endTime)
	if err != nil {
		slog.ErrorContext(ctx, "failed to count invoices : ", err)
	}
	date, err := s.Querier.GetAllInvoiceWithGivenDate(ctx, startTime, endTime, page, size)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get invoices : ", err)
	}
	pagination := pkg.CalculatePagination(page, size, int(count))

	response := model.GetInvoiceResponse{
		Invoice: date,
	}
	return response, pagination, nil
}
