package app

import (
	"context"

	"github.com/aygumov-g/service-SSO-go/internal/config"
	"github.com/aygumov-g/service-SSO-go/internal/http/server"
	"github.com/aygumov-g/service-SSO-go/internal/logger"
	"github.com/aygumov-g/service-SSO-go/internal/storage/postgres"
)

type App struct {
	httpServer *server.Server
	db         *postgres.DB
	logger     logger.Logger
}

func New(ctx context.Context) (*App, error) {
	log := logger.New()
	cfg := config.Load()

	db, err := postgres.New(ctx, cfg.MainDB.DSN())
	if err != nil {
		return nil, err
	}

	httpServer := buildHTTP(cfg, db.Get(), log)

	return &App{
		httpServer: httpServer,
		db:         db,
		logger:     log,
	}, nil
}
