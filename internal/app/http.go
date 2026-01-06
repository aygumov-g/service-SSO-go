package app

import (
	"github.com/aygumov-g/service-SSO-go/internal/config"
	"github.com/aygumov-g/service-SSO-go/internal/domain/auth"
	"github.com/aygumov-g/service-SSO-go/internal/http/handler"
	"github.com/aygumov-g/service-SSO-go/internal/http/middleware"
	"github.com/aygumov-g/service-SSO-go/internal/http/router"
	"github.com/aygumov-g/service-SSO-go/internal/http/server"
	"github.com/aygumov-g/service-SSO-go/internal/logger"
	"github.com/aygumov-g/service-SSO-go/internal/storage/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func buildHTTP(cfg *config.Config, pool *pgxpool.Pool, log logger.Logger) *server.Server {
	jwtManager := auth.NewJWTManager(auth.Config{
		Secret: []byte(cfg.JWT.Secret),
		TTL:    cfg.JWT.TTL,
	})

	authUserRepo := postgres.NewAuthUserRepository(pool)
	authService := auth.NewService(authUserRepo, jwtManager)

	registerHandler := handler.NewRegisterHandler(authService)
	loginHandler := handler.NewLoginHandler(authService)
	meHandler := handler.NewMeHandler(authUserRepo)

	authMW := middleware.Auth(jwtManager)

	r := router.New()

	r.Handle("/auth/register", registerHandler)
	r.Handle("/auth/login", loginHandler)
	r.Handle("/auth/me", authMW(meHandler))

	r.Use(middleware.Logging(log))

	return server.New(":"+cfg.AppPort, r.Handler())
}
