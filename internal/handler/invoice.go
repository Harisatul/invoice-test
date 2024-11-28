package handler

import (
	"encoding/json"
	"invoice-test/internal/model"
	"invoice-test/pkg"
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
	}
	pkg.WriteSuccessResponse(w, http.StatusOK, "success create invoice", invoice)
}
