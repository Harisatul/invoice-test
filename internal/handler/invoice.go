package handler

import (
	"encoding/json"
	"invoice-test/internal/model"
	"invoice-test/pkg"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func (h Handler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var createInvoiceRequest model.CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&createInvoiceRequest); err != nil {
		pkg.WriteErrorResponse(w, http.StatusBadRequest, "invalid json body", err)
		return
	}
	invoice, err := h.Service.CreateInvoice(r.Context(), createInvoiceRequest)
	if err != nil {
		pkg.WriteErrorResponse(w, http.StatusInternalServerError, "create invoice", err)
		return
	}
	pkg.WriteSuccessResponse(w, http.StatusOK, "success create invoice", invoice, nil)
	return
}

func (h Handler) DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := h.Service.DeleteInvoice(r.Context(), id)
	if err != nil {
		slog.Warn(err.Error())
		pkg.WriteErrorResponse(w, http.StatusBadRequest, "failed to delete invoice", err.Error())
		return
	}
	pkg.WriteSuccessResponse(w, http.StatusOK, "success delete invoice", nil, nil)
	return
}

func (h Handler) UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var updateInvoiceRequest model.UpdateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&updateInvoiceRequest); err != nil {
		pkg.WriteErrorResponse(w, http.StatusBadRequest, "invalid json body", err)
		return
	}
	invoice, err := h.Service.UpdateInvoice(r.Context(), updateInvoiceRequest, id)
	if err != nil {
		if err.Error() == "given id not found" || err.Error() == "invalid payment status" {
			pkg.WriteErrorResponse(w, http.StatusBadRequest, "failed to update invoice", err.Error())
			return
		}
		pkg.WriteErrorResponse(w, http.StatusInternalServerError, "failed to update invoice", err.Error())
		return
	}
	pkg.WriteSuccessResponse(w, http.StatusOK, "success update invoice", invoice, nil)
	return
}

func (h Handler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	var startDate, endDate time.Time
	var err error
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			pkg.WriteErrorResponse(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD.", err.Error())
			return
		}
	}
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			pkg.WriteErrorResponse(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD.", err.Error())
			return
		}
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil && pageStr != "" {
		pkg.WriteErrorResponse(w, http.StatusBadRequest, "Invalid page number.", err.Error())
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil && sizeStr != "" {
		pkg.WriteErrorResponse(w, http.StatusBadRequest, "Invalid size number.", err.Error())
		return
	}
	invoice, index, err := h.Service.GetAllInvoice(r.Context(), startDate, endDate, page, size)
	pkg.WriteSuccessResponse(w, http.StatusOK, "success update invoice", invoice, index)
	return
}
