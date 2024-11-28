package cmd

import (
	"context"
	"fmt"
	"golang.org/x/exp/slog"
	"invoice-test/config"
	"invoice-test/internal/handler"
	"invoice-test/internal/repository"
	"invoice-test/internal/service"
	"net/http"
	"time"
)

func RegisterRoute(mux *http.ServeMux, handler handler.Handler) {
	mux.HandleFunc("GET /api/health-check", handler.HealthCheckHandler)
	mux.HandleFunc("POST /api/invoice", handler.CreateInvoice)
	mux.HandleFunc("DELETE /api/invoice", handler.DeleteInvoice)
	mux.HandleFunc("PUT /api/invoice", handler.UpdateInvoice)
}

func runHTTPServer(ctx context.Context) {

	cfg := config.NewViper("app")

	db := config.NewDb(ctx, cfg)
	defer db.Close()

	queries := repository.New(db)
	svc := service.Service{Db: db, Querier: queries}
	handler := handler.Handler{Service: svc}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         cfg.GetString("server.port"),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.GetInt("server.read_timeout")) * time.Second,
		WriteTimeout: time.Duration(cfg.GetInt("server.write_timeout")) * time.Second,
	}

	RegisterRoute(mux, handler)

	go func() {
		slog.Info(fmt.Sprintf("starting server at port %s", cfg.GetString("server.port")))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("unable to start server", err)
		}
	}()
	<-ctx.Done()

	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	slog.Info("shutting down server")
	if err := srv.Shutdown(timeout); err != nil {
		slog.Error("unable to shutdown server")
	}
	slog.Info("server down gracefully")
}
