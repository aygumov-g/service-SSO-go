package app

import (
	"context"
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/config"
	"github.com/aygumov-g/service-SSO-go/internal/infrastructure/health"
	postgres_db "github.com/aygumov-g/service-SSO-go/internal/infrastructure/postgres"
	"github.com/aygumov-g/service-SSO-go/internal/infrastructure/security/password"
	"github.com/aygumov-g/service-SSO-go/internal/infrastructure/security/token"
	account_db "github.com/aygumov-g/service-SSO-go/internal/repository/postgres/account"
	authorization_code_db "github.com/aygumov-g/service-SSO-go/internal/repository/postgres/authorization_code"
	oauth_client_db "github.com/aygumov-g/service-SSO-go/internal/repository/postgres/oauth_client"
	session_db "github.com/aygumov-g/service-SSO-go/internal/repository/postgres/session"
	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"
	authorize_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/authorize"
	change_password_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/change_password"
	login_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/login"
	me_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/me"
	ready_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/ready"
	refresh_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/refresh"
	register_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/register"
	token_handler "github.com/aygumov-g/service-SSO-go/internal/transport/http/handler/token"
	auth_id "github.com/aygumov-g/service-SSO-go/internal/transport/http/identity/auth"
	auth_mw "github.com/aygumov-g/service-SSO-go/internal/transport/http/middleware/auth"
	methods_mw "github.com/aygumov-g/service-SSO-go/internal/transport/http/middleware/method"
	"github.com/aygumov-g/service-SSO-go/internal/transport/http/router"
	"github.com/aygumov-g/service-SSO-go/internal/transport/http/server"
	authorize_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/authorize"
	change_password_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/change_password"
	get_me_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/get_me"
	login_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/login"
	ready_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/ready"
	refresh_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/refresh"
	register_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/register"
	token_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/token"
	"github.com/aygumov-g/service-SSO-go/pkg/clock"
	"github.com/aygumov-g/service-SSO-go/pkg/logger"
)

type App struct {
	logger     logger.Logger
	httpServer *server.Server
	db         *postgres_db.DB
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.Load()
	log := logger.New()

	db, err := postgres_db.New(ctx, cfg.DB.DSN())
	if err != nil {
		return nil, err
	}

	healthchecker := health.New(db)

	clk := clock.NewSystemClock()

	passwordHasher := password.NewBcryptHasher(12)
	tokenProvider := token.NewJWTProvider(cfg.JWT.Secret, cfg.JWT.TTL, clk)

	accountRepo := account_db.NewRepository(db)
	sessionRepo := session_db.NewRepository(db)
	authorization_codeRepo := authorization_code_db.NewRepository(db)
	oauth_clientRepo := oauth_client_db.NewRepository(db)

	sessionService := session_srv.NewService(db, sessionRepo, accountRepo, tokenProvider, clk, cfg.Refresh.TTL)

	tokenUsecase := token_uc.NewUsecase(db, authorization_codeRepo, oauth_clientRepo, accountRepo, sessionService, tokenProvider, clk)
	authorizeUsecase := authorize_uc.NewUsecase(authorization_codeRepo, oauth_clientRepo, tokenProvider, clk)
	change_passwordUsecase := change_password_uc.NewUsecase(accountRepo, sessionService, passwordHasher, clk)
	registerUsecase := register_uc.NewUsecase(accountRepo, passwordHasher, clk)
	refreshUsecase := refresh_uc.NewUsecase(accountRepo, sessionService)
	loginUsecase := login_uc.NewUsecase(db, accountRepo, sessionService, passwordHasher)
	getMeUsecase := get_me_uc.NewUsecase(accountRepo)
	readyUsecase := ready_uc.NewUsecase(healthchecker)

	authIdentity := auth_id.NewIdentity()

	authMW := auth_mw.NewMiddleware(getMeUsecase, authIdentity, tokenProvider)
	methodsMW := methods_mw.NewMiddleware()

	meHandler := me_handler.NewHandler(getMeUsecase, authIdentity)
	authorizeHandler := authorize_handler.NewHandler(authorizeUsecase, authIdentity)
	change_passwordHandler := change_password_handler.NewHandler(change_passwordUsecase, authIdentity)
	registerHandler := register_handler.NewHandler(registerUsecase)
	refreshHandler := refresh_handler.NewHandler(refreshUsecase)
	tokenHandler := token_handler.NewHandler(tokenUsecase)
	loginHandler := login_handler.NewHandler(loginUsecase)
	readyHandler := ready_handler.NewHandler(readyUsecase)

	r := router.NewRouter()
	r.Handle("/auth/me", methodsMW.Handle([]string{http.MethodGet}, authMW.Handle(meHandler)))
	r.Handle("/auth/change_password", methodsMW.Handle([]string{http.MethodPost}, authMW.Handle(change_passwordHandler)))
	r.Handle("/auth/register", methodsMW.Handle([]string{http.MethodPost}, registerHandler))
	r.Handle("/auth/refresh", methodsMW.Handle([]string{http.MethodPost}, refreshHandler))
	r.Handle("/auth/login", methodsMW.Handle([]string{http.MethodPost}, loginHandler))
	r.Handle("/oauth/authorize", methodsMW.Handle([]string{http.MethodGet}, authMW.Handle(authorizeHandler)))
	r.Handle("/oauth/token", methodsMW.Handle([]string{http.MethodPost}, tokenHandler))
	r.Handle("/health/ready", methodsMW.Handle([]string{http.MethodGet}, readyHandler))

	s := server.NewServer(cfg.App.Port, r.Handler())

	return &App{logger: log, httpServer: s, db: db}, nil
}

func (a *App) Run() {
	a.logger.Info("http server started", "addr", a.httpServer.Addr())

	if err := a.httpServer.Start(); err != nil {
		a.logger.Error("http server failed started", "error", err)
	}
}

func (a *App) Shutdown(ctx context.Context) {
	a.logger.Info("shutdown started")

	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		a.logger.Error("http server failed shutdown", "error", err)
	}

	a.db.Close()

	a.logger.Info("shutdown completed")
}
