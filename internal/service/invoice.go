package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"invoice-test/internal/model"
	"invoice-test/pkg"
	"log/slog"
	"mime/multipart"
	"strconv"
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
	cahsOnly, err := s.Querier.CountAllInvoiceWithGivenDateAndCashOnly(ctx, startTime, endTime)
	if err != nil {
		slog.ErrorContext(ctx, "failed to count cash only invoices : ", err)

	}
	invoice, err := s.Querier.GetSumofAllInvoice(ctx, startTime, endTime)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get invoices : ", err)
	}

	var sumOftotalCogs int64
	var sumOftotalPrice int64
	for _, inv := range invoice {
		sumOftotalCogs += inv.TotalCOGS
		sumOftotalPrice += inv.TotalPriceSold
	}

	pagination := pkg.CalculatePagination(page, size, int(count))
	aggregateResponse := model.InvoiceAggregateResponse{
		TotalProfit:            sumOftotalPrice - sumOftotalCogs,
		TotalOfCashTransaction: cahsOnly,
	}

	response := model.GetInvoiceResponse{
		Invoice:                  date,
		InvoiceAggregateResponse: aggregateResponse,
	}
	return response, pagination, nil
}

func (s Service) ImportXLSX(ctx context.Context, file multipart.File) error {
	xlsxFile, err := excelize.OpenReader(file)
	if err != nil {
		slog.ErrorContext(ctx, "failed to open xlsx file : ", err)
		return errors.New("failed to open xlsx file")
	}
	defer xlsxFile.Close()

	sheetInvoiceName := xlsxFile.GetSheetName(0)
	invoiceSheet, err := xlsxFile.Rows(sheetInvoiceName)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read sheet invoice : ", err)
	}

	sheetProductName := xlsxFile.GetSheetName(1)
	productSheet, err := xlsxFile.Rows(sheetProductName)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read sheet product : ", err)
	}

	rowCount := 0
	for invoiceSheet.Next() {

		fmt.Println(rowCount)

		// Skip the first row
		if rowCount == 0 {
			rowCount++
			continue
		}

		columns, err := invoiceSheet.Columns()
		if err != nil {
			slog.ErrorContext(ctx, "failed to read columns: ", err)
		}
		var startDate time.Time
		value := columns[1]
		startDate, err = time.Parse("01-02-06", value)
		if err != nil {
			rowCount++
			continue
		}
		status, err := pkg.ToPaymentStatus(columns[4])
		if err != nil {
			rowCount++
			continue
		}
		invoice := model.Invoice{
			InvoiceNumber: columns[0],
			Date:          startDate,
			CustomerName:  columns[2],
			Salesperson:   columns[3],
			Notes:         columns[5],
			PaymentType:   status,
		}
		_, err = s.Querier.InsertInvoice(ctx, invoice)
		if err != nil {
			rowCount++
			continue
		}
		rowCount++
	}

	productRowCount := 0
	for productSheet.Next() {

		fmt.Println(productRowCount)
		if productRowCount == 0 {
			productRowCount++
			continue
		}

		columns, err := productSheet.Columns()
		if err != nil {
			slog.ErrorContext(ctx, "failed to read columns: ", err)
		}
		slog.InfoContext(ctx, "read columns: ", slog.String("column", columns[0]))

		id, err := pkg.GenerateId()
		if err != nil {
			rowCount++
			continue
		}

		quantity, err := strconv.Atoi(columns[2])
		if err != nil {
			rowCount++
			continue
		}

		totalCogs, err := strconv.Atoi(columns[3])
		if err != nil {
			rowCount++
			continue
		}
		totalPrice, err := strconv.Atoi(columns[4])
		if err != nil {
			rowCount++
			continue
		}

		product := model.Product{
			ID:             id,
			ItemName:       columns[1],
			Quantity:       quantity,
			TotalCOGS:      int64(totalCogs),
			TotalPriceSold: int64(totalPrice),
			InvoiceNumber:  columns[0],
		}
		_, err = s.Querier.InsertProduct(ctx, product)
		if err != nil {
			productRowCount++
			slog.ErrorContext(ctx, "failed to insert product : ", err)
		}
		productRowCount++
	}
	return nil
}
