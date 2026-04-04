package authorize

import (
	"context"
	"time"

	authorization_code_d "github.com/aygumov-g/service-SSO-go/internal/domain/authorization_code"
	oauth_client_d "github.com/aygumov-g/service-SSO-go/internal/domain/oauth_client"
)

type Usecase struct {
	authorization_codesRepo AuthorizationCodeRepository
	oauth_clientsRepo       OAuthClientRepository
	tokens                  TokenProvider
	clk                     Clock
}

func NewUsecase(
	authorization_codesRepo AuthorizationCodeRepository,
	oauth_clientsRepo OAuthClientRepository,
	tokens TokenProvider,
	clk Clock,
) *Usecase {
	return &Usecase{
		authorization_codesRepo: authorization_codesRepo,
		oauth_clientsRepo:       oauth_clientsRepo,
		tokens:                  tokens,
		clk:                     clk,
	}
}

func (uc *Usecase) Execute(ctx context.Context, in Input) (string, error) {
	client, err := uc.oauth_clientsRepo.GetByClientID(ctx, in.ClientID)
	if err != nil {
		return "", err
	}

	if client.RedirectURI != in.RedirectURI {
		return "", oauth_client_d.ErrInvalidRedirectURI
	}

	code, err := uc.tokens.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	now := uc.clk.Now()

	entity := &authorization_code_d.AuthorizationCode{
		Code:        code,
		AccountID:   in.AccountID,
		ClientID:    in.ClientID,
		RedirectURI: in.RedirectURI,
		ExpiresAt:   now.Add(time.Minute),
		CreatedAt:   now,
		Used:        false,
	}

	if err := uc.authorization_codesRepo.Create(ctx, entity); err != nil {
		return "", err
	}

	return code, nil
}
