package app

import (
	"context"

	"github.com/aygumov-g/service-SSO-go/internal/config"
	"github.com/aygumov-g/service-SSO-go/internal/infrastructure/postgres"
	"github.com/aygumov-g/service-SSO-go/internal/transport/http/server"
	"github.com/aygumov-g/service-SSO-go/pkg/logger"
)

type App struct {
	logger     logger.Logger
	httpServer *server.Server
	db         *postgres.DB
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.Load()
	log := logger.New()

	db, err := postgres.New(ctx, cfg.DB.DSN())
	if err != nil {
		return nil, err
	}

	httpServer := buildHTTP(cfg, db.Get(), log)

	return &App{
		logger:     log,
		httpServer: httpServer,
		db:         db,
	}, nil
}

func (a *App) Run() {
	a.logger.Info("http server started", "addr", a.httpServer.Addr())

	if err := a.httpServer.Start(); err != nil {
		a.logger.Error("http server failed", "error", err)
	}
}

func (a *App) Shutdown(ctx context.Context) {
	a.logger.Info("shutdown started")

	_ = a.httpServer.Shutdown(ctx)
	a.db.Close()

	a.logger.Info("shutdown completed")
}
