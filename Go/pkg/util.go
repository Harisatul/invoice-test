package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	model2 "invoice-test/internal/model"
	"math/rand"
	"net/http"
	"time"
)

func GenerateId() (uuid.UUID, error) {
	v7, err := uuid.NewV7()
	return v7, err
}

func InvoiceNumberGenerator() string {
	rand.Seed(time.Now().UnixNano())

	randomString := fmt.Sprintf("%04d", rand.Intn(10000))

	randomNumber := fmt.Sprintf("%06d", rand.Intn(1000000))

	invoiceNumber := fmt.Sprintf("INV-%s-%s", randomString, randomNumber)

	return invoiceNumber
}

func ToPaymentStatus(input string) (model2.PaymentStatus, error) {
	switch input {
	case string(model2.PaymentStatusCASH):
		return model2.PaymentStatusCASH, nil
	case string(model2.PaymentStatusCREDIT):
		return model2.PaymentStatusCREDIT, nil
	default:
		return "", errors.New("invalid payment status")
	}
}

func CalculatePagination(page, pageSize, totalCount int) model2.PaginationIndex {
	// Calculate pagination metadata
	totalPages := (totalCount + pageSize - 1) / pageSize
	return model2.PaginationIndex{
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		HasPrevious: page > 1,
		HasNext:     page < totalPages,
	}
}

func WriteSuccessResponse(w http.ResponseWriter, status int, message string, data interface{}, index interface{}) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model2.APIResponse{
		Status:          status,
		Message:         message,
		Data:            data,
		PaginationIndex: index,
	})
}

func WriteErrorResponse(w http.ResponseWriter, status int, message string, err interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model2.APIResponse{
		Status:  status,
		Message: message,
		Errors:  err,
	})
}
