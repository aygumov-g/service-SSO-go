package app

import (
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/config"
	account_db "github.com/aygumov-g/service-SSO-go/internal/repository/postgres/account"
	session_db "github.com/aygumov-g/service-SSO-go/internal/repository/postgres/session"
	account_srv "github.com/aygumov-g/service-SSO-go/internal/service/account"
	jwt_srv "github.com/aygumov-g/service-SSO-go/internal/service/jwt"
	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"
	change_password_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/change_password"
	login_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/login"
	me_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/me"
	refresh_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/refresh"
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

	jwtService := jwt_srv.NewJWTService([]byte(cfg.JWT.Secret), cfg.JWT.TTL, clk)

	sessionRepo := session_db.NewRepository(db)
	accountRepo := account_db.NewRepository(db)

	sessionService := session_srv.NewService(sessionRepo, accountRepo, jwtService, cfg.Refresh.TTL, clk)
	accountService := account_srv.NewService(accountRepo, sessionService, clk)

	authIdentity := auth_id.NewIdentity()

	authMW := auth_mw.NewMiddleware(jwtService, authIdentity, accountService)
	methodsMW := methods_mw.NewMiddleware()

	meHandler := me_handler.NewHandler(accountService, authIdentity)
	change_passwordHandler := change_password_handler.NewHandler(accountService, authIdentity)
	refreshHandler := refresh_handler.NewHandler(sessionService)
	registerHandler := register_handler.NewHandler(accountService)
	loginHandler := login_handler.NewHandler(accountService)

	r := router.NewRouter()
	r.Handle("/accounts/me", methodsMW.Handle([]string{http.MethodGet}, authMW.Handle(meHandler)))
	r.Handle("/accounts/change_password", methodsMW.Handle([]string{http.MethodPost}, authMW.Handle(change_passwordHandler)))
	r.Handle("/accounts/refresh", methodsMW.Handle([]string{http.MethodPost}, refreshHandler))
	r.Handle("/accounts/register", methodsMW.Handle([]string{http.MethodPost}, registerHandler))
	r.Handle("/accounts/login", methodsMW.Handle([]string{http.MethodPost}, loginHandler))

	return server.NewServer(cfg.App.Port, r.Handler())
}
