package token

import (
	"context"

	authorization_code_d "github.com/aygumov-g/service-SSO-go/internal/domain/authorization_code"
	oauth_client_d "github.com/aygumov-g/service-SSO-go/internal/domain/oauth_client"
	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"
)

type Usecase struct {
	tx                      TxManager
	authorization_codesRepo AuthorizationCodeRepository
	oauth_clientsRepo       OAuthClientRepository
	accountsRepo            AccountRepository
	sessionsSrv             SessionService
	tokens                  TokenProvider
	clk                     Clock
}

func NewUsecase(
	tx TxManager,
	authorization_codesRepo AuthorizationCodeRepository,
	oauth_clientsRepo OAuthClientRepository,
	accountsRepo AccountRepository,
	sessionsSrv SessionService,
	tokens TokenProvider,
	clk Clock,
) *Usecase {
	return &Usecase{
		tx:                      tx,
		authorization_codesRepo: authorization_codesRepo,
		oauth_clientsRepo:       oauth_clientsRepo,
		accountsRepo:            accountsRepo,
		sessionsSrv:             sessionsSrv,
		tokens:                  tokens,
		clk:                     clk,
	}
}

func (uc *Usecase) Execute(ctx context.Context, in Input) (*Output, error) {
	var tokens *session_srv.Output
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		client, err := uc.oauth_clientsRepo.GetByClientID(txCtx, in.ClientID)
		if err != nil {
			return oauth_client_d.ErrInvalidClient
		}

		if client.Secret != in.ClientSecret {
			return oauth_client_d.ErrInvalidSecret
		}

		if client.RedirectURI != in.RedirectURI {
			return oauth_client_d.ErrInvalidRedirectURI
		}

		codeEntity, err := uc.authorization_codesRepo.UseCode(txCtx, in.Code)
		if err != nil {
			return err
		}

		if uc.clk.Now().After(codeEntity.ExpiresAt) {
			return authorization_code_d.ErrAuthorizationCodeExpired
		}

		account, err := uc.accountsRepo.GetByID(txCtx, codeEntity.AccountID)
		if err != nil {
			return err
		}

		tokens, err = uc.sessionsSrv.Create(txCtx, account)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &Output{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}
