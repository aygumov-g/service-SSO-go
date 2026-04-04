package token

import (
	"context"

	authorization_code_d "github.com/aygumov-g/service-SSO-go/internal/domain/authorization_code"
	oauth_client_d "github.com/aygumov-g/service-SSO-go/internal/domain/oauth_client"
)

type Usecase struct {
	authorization_codesRepo AuthorizationCodeRepository
	oauth_clientsRepo       OAuthClientRepository
	accountsRepo            AccountRepository
	sessionsSrv             SessionService
	tokens                  TokenProvider
	clk                     Clock
}

func NewUsecase(
	authorization_codesRepo AuthorizationCodeRepository,
	oauth_clientsRepo OAuthClientRepository,
	accountsRepo AccountRepository,
	sessionsSrv SessionService,
	tokens TokenProvider,
	clk Clock,
) *Usecase {
	return &Usecase{
		authorization_codesRepo: authorization_codesRepo,
		oauth_clientsRepo:       oauth_clientsRepo,
		accountsRepo:            accountsRepo,
		sessionsSrv:             sessionsSrv,
		tokens:                  tokens,
		clk:                     clk,
	}
}

func (uc *Usecase) Execute(ctx context.Context, in Input) (*Output, error) {
	client, err := uc.oauth_clientsRepo.GetByClientID(ctx, in.ClientID)
	if err != nil {
		return nil, oauth_client_d.ErrInvalidClient
	}

	if client.Secret != in.ClientSecret {
		return nil, oauth_client_d.ErrInvalidSecret
	}

	if client.RedirectURI != in.RedirectURI {
		return nil, oauth_client_d.ErrInvalidRedirectURI
	}

	codeEntity, err := uc.authorization_codesRepo.UseCode(ctx, in.Code)
	if err != nil {
		return nil, authorization_code_d.ErrInvalidAuthorizationCode
	}

	if uc.clk.Now().After(codeEntity.ExpiresAt) {
		return nil, authorization_code_d.ErrAuthorizationCodeExpired
	}

	account, err := uc.accountsRepo.GetByID(ctx, codeEntity.AccountID)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := uc.sessionsSrv.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return &Output{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
