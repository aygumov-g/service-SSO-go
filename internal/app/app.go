package app

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_http "github.com/aygumov-g/service-SSO-go/internal/http"
	"github.com/aygumov-g/service-SSO-go/internal/logger"
)

func Run() {
	log := logger.NewLogger()
	log.Info("application starting")

	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))

	r := _http.NewRouter(tmpl, log)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("http server, listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server failed", "error", err)
		}
	}()

	<-ctx.Done()
	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "error", err)
	} else {
		log.Info("server stopped gracefully")
	}
}
