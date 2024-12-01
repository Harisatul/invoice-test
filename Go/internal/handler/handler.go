package handler

import (
	"invoice-test/internal/service"
	"invoice-test/pkg"
	"net/http"
)

type Handler struct {
	Service service.Service
}

func (h Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	pkg.WriteSuccessResponse(w, http.StatusOK, "success", "hello from server 👋", nil)
}
