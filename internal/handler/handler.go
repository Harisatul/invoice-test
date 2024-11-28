package handler

import (
	"encoding/json"
	"invoice-test/internal/service"
	"net/http"
)

type Handler struct {
	Service service.Service
}

func (h Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Hello From server 👋")
}
