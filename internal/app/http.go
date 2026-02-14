package app

import (
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/config"
	"github.com/aygumov-g/service-SSO-go/internal/repository/postgres/user"
	"github.com/aygumov-g/service-SSO-go/internal/service/jwt"
	user_srv "github.com/aygumov-g/service-SSO-go/internal/service/user"
	login_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/login"
	me_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/me"
	register_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/register"
	auth_id "github.com/aygumov-g/service-SSO-go/internal/transport/http/identity/auth"
	auth_mw "github.com/aygumov-g/service-SSO-go/internal/transport/http/middleware/auth"
	methods_mw "github.com/aygumov-g/service-SSO-go/internal/transport/http/middleware/method"
	"github.com/aygumov-g/service-SSO-go/internal/transport/http/router"
	"github.com/aygumov-g/service-SSO-go/internal/transport/http/server"
	"github.com/aygumov-g/service-SSO-go/pkg/clock"
	"github.com/aygumov-g/service-SSO-go/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func buildHTTP(cfg *config.Config, db *pgxpool.Pool, log logger.Logger) *server.Server {
	clk := clock.NewSystemClock()

	jwtService := jwt.NewJWTService([]byte(cfg.JWT.Secret), cfg.JWT.TTL, clk)

	userRepo := user.NewRepository(db)
	userService := user_srv.NewService(userRepo, jwtService)

	authIdentity := auth_id.NewIdentity("identity")

	authMW := auth_mw.NewMiddleware(jwtService, authIdentity)
	methodsMW := methods_mw.NewMiddleware()

	meHandler := me_handler.NewHandler(userService, authIdentity)
	registerHandler := register_handler.NewHandler(userService)
	loginHandler := login_handler.NewHandler(userService)

	r := router.NewRouter()
	r.Handle("/auth/me", methodsMW.Handle([]string{http.MethodGet}, authMW.Handle(meHandler)))
	r.Handle("/auth/register", methodsMW.Handle([]string{http.MethodPost}, registerHandler))
	r.Handle("/auth/login", methodsMW.Handle([]string{http.MethodPost}, loginHandler))

	return server.NewServer(cfg.AppPort, r.Handler())
}
