package service

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"invoice-test/internal/repository"
	"log/slog"
)

type Service struct {
	Db      *pgxpool.Pool
	Querier *repository.Queries
}

func (s Service) HealthCheck(context.Context) {
	slog.Info("Server is running")
}
