package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func newDb(ctx context.Context, cfg *viper.Viper) *pgxpool.Pool {
	username := cfg.GetString("db.user")
	password := cfg.GetString("db.password")
	host := cfg.GetString("db.host")
	port := cfg.GetString("db.port")
	database := cfg.GetString("db.name")
	maxConn := cfg.GetString("db.pool.max")
	minConn := cfg.GetString("db.pool.min")
	timezone := cfg.GetString("server.timezone")
	maxIdleTime := cfg.GetString("db.pool.max_idle_time")
	maxConnLifetime := cfg.GetDuration("db.pool.max_conn_lifetime")

	connTemplate := "postgres://%s:%s@%s:%d/%s?pool_max_conns=%d&pool_min_conns=%d&timezone=%s&pool_max_conn_lifetime=%s&pool_max_connection_idle=%s"
	connString := fmt.Sprintf(connTemplate, username, password, host, port, database, maxConn, minConn, timezone, maxConnLifetime, maxIdleTime)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		panic(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		panic(err)
	}
	return pool
}
