package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"invoice-test/internal/model"
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

func ToPaymentStatus(input string) (model.PaymentStatus, error) {
	switch input {
	case string(model.PaymentStatusCASH):
		return model.PaymentStatusCASH, nil
	case string(model.PaymentStatusCREDIT):
		return model.PaymentStatusCREDIT, nil
	default:
		return "", errors.New("invalid payment status")
	}
}

func WriteSuccessResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func WriteErrorResponse(w http.ResponseWriter, status int, message string, err interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.APIResponse{
		Status:  status,
		Message: message,
		Errors:  err,
	})
}
