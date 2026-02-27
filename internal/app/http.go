package app

import (
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/config"
	"github.com/aygumov-g/service-SSO-go/internal/infrastructure/security/password"
	"github.com/aygumov-g/service-SSO-go/internal/infrastructure/security/token"
	account_db "github.com/aygumov-g/service-SSO-go/internal/repository/postgres/account"
	session_db "github.com/aygumov-g/service-SSO-go/internal/repository/postgres/session"
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
	change_password_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/change_password"
	get_me_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/get_me"
	login_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/login"
	refresh_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/refresh"
	register_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/register"
	"github.com/aygumov-g/service-SSO-go/pkg/clock"
	"github.com/aygumov-g/service-SSO-go/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func buildHTTP(cfg *config.Config, db *pgxpool.Pool, log logger.Logger) *server.Server {
	clk := clock.NewSystemClock()

	passwordHasher := password.NewBcryptHasher(12)
	tokenProvider := token.NewJWTProvider(cfg.JWT.Secret, cfg.JWT.TTL, clk)

	sessionRepo := session_db.NewRepository(db)
	accountRepo := account_db.NewRepository(db)

	sessionService := session_srv.NewService(sessionRepo, tokenProvider, clk, cfg.Refresh.TTL)

	chahge_passwordUsecase := change_password_uc.NewChangePassword(accountRepo, sessionService, passwordHasher, clk)
	registerUsecase := register_uc.NewRegister(accountRepo, passwordHasher, clk)
	refreshUsecase := refresh_uc.NewRefresh(accountRepo, sessionService)
	loginUsecase := login_uc.NewLogin(accountRepo, sessionService, passwordHasher)
	getMeUsecase := get_me_uc.NewGetMe(accountRepo)

	authIdentity := auth_id.NewIdentity()

	authMW := auth_mw.NewMiddleware(getMeUsecase, authIdentity, tokenProvider)
	methodsMW := methods_mw.NewMiddleware()

	meHandler := me_handler.NewHandler(getMeUsecase, authIdentity)
	change_passwordHandler := change_password_handler.NewHandler(chahge_passwordUsecase, authIdentity)
	refreshHandler := refresh_handler.NewHandler(refreshUsecase)
	registerHandler := register_handler.NewHandler(registerUsecase)
	loginHandler := login_handler.NewHandler(loginUsecase)

	r := router.NewRouter()
	r.Handle("/accounts/me", methodsMW.Handle([]string{http.MethodGet}, authMW.Handle(meHandler)))
	r.Handle("/accounts/change_password", methodsMW.Handle([]string{http.MethodPost}, authMW.Handle(change_passwordHandler)))
	r.Handle("/accounts/register", methodsMW.Handle([]string{http.MethodPost}, registerHandler))
	r.Handle("/accounts/refresh", methodsMW.Handle([]string{http.MethodPost}, refreshHandler))
	r.Handle("/accounts/login", methodsMW.Handle([]string{http.MethodPost}, loginHandler))

	return server.NewServer(cfg.App.Port, r.Handler())
}
