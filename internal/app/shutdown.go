package app

import "context"

func (a *App) Shutdown(ctx context.Context) {
	a.logger.Info("shutdown started")

	_ = a.httpServer.Shutdown(ctx)
	a.db.Close()

	a.logger.Info("shutdown completed")
}
