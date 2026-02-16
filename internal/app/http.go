package app

import (
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/config"
	account_db "github.com/aygumov-g/service-SSO-go/internal/repository/postgres/account"
	account_srv "github.com/aygumov-g/service-SSO-go/internal/service/account"
	"github.com/aygumov-g/service-SSO-go/internal/service/jwt"
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

	accountRepo := account_db.NewRepository(db)
	accountService := account_srv.NewService(accountRepo, jwtService, clk)

	authIdentity := auth_id.NewIdentity("identity")

	authMW := auth_mw.NewMiddleware(jwtService, authIdentity)
	methodsMW := methods_mw.NewMiddleware()

	meHandler := me_handler.NewHandler(accountService, authIdentity)
	registerHandler := register_handler.NewHandler(accountService)
	loginHandler := login_handler.NewHandler(accountService)

	r := router.NewRouter()
	r.Handle("/accounts/me", methodsMW.Handle([]string{http.MethodGet}, authMW.Handle(meHandler)))
	r.Handle("/accounts/register", methodsMW.Handle([]string{http.MethodPost}, registerHandler))
	r.Handle("/accounts/login", methodsMW.Handle([]string{http.MethodPost}, loginHandler))

	return server.NewServer(cfg.AppPort, r.Handler())
}
