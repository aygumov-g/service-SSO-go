package app

func (a *App) Run() {
	a.logger.Info("http server started", "addr", a.httpServer.Addr())

	if err := a.httpServer.Start(); err != nil {
		a.logger.Error("http server failed", "error", err)
	}
}
