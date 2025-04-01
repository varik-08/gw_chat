package config

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
)

func initDB(dbCfg DB) *pgxpool.Pool {
	conn, err := newDB(dbCfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	pool = conn

	log.Println("Connected to database")

	return pool
}

func newDB(dbCfg DB) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//nolint
	connConfig, err := pgx.ParseConfig(
		fmt.Sprintf(
			"postgres://%s:%s@%s/%s?TimeZone=Europe/Moscow&search_path=%s",
			dbCfg.User,
			dbCfg.Password,
			net.JoinHostPort(dbCfg.Host, dbCfg.Port),
			dbCfg.Database,
			dbCfg.Schema,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create DSN for DB connection: %w", err)
	}
	dbc, err := pgxpool.New(ctx, connConfig.ConnString())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB : %w", err)
	}

	err = dbc.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return dbc, nil
}
