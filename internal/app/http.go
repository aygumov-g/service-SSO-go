package app

import (
	"html/template"
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/config"
	"github.com/aygumov-g/service-SSO-go/internal/domain/user"
	"github.com/aygumov-g/service-SSO-go/internal/http/handler"
	"github.com/aygumov-g/service-SSO-go/internal/http/middleware"
	"github.com/aygumov-g/service-SSO-go/internal/http/router"
	"github.com/aygumov-g/service-SSO-go/internal/http/server"
	"github.com/aygumov-g/service-SSO-go/internal/logger"
	"github.com/aygumov-g/service-SSO-go/internal/storage/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func buildHTTP(cfg *config.Config, pool *pgxpool.Pool, log logger.Logger) *server.Server {
	userRepo := postgres.NewUserRepository(pool)
	userService := user.NewService(userRepo)

	tmpl := template.Must(
		template.ParseFiles("web/templates/index.html"),
	)

	indexHandler := handler.NewIndexHandler(tmpl)
	testHandler := handler.NewTestHandler(userService)

	r := router.New()

	r.Handle("/", indexHandler)
	r.Handle("/test", testHandler)

	staticsFS := http.FileServer(http.Dir("web/statics"))
	r.Handle("/stt/", http.StripPrefix("/stt/", staticsFS))

	r.Use(middleware.Logging(log))

	return server.New(":"+cfg.AppPort, r.Handler())
}
