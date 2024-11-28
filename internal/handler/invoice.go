package handler

import (
	"encoding/json"
	"invoice-test/internal/model"
	"invoice-test/pkg"
	"log/slog"
	"net/http"
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
	pkg.WriteSuccessResponse(w, http.StatusOK, "success create invoice", invoice)
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
	pkg.WriteSuccessResponse(w, http.StatusOK, "success delete invoice", nil)
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
	pkg.WriteSuccessResponse(w, http.StatusOK, "success update invoice", invoice)
	return
}
