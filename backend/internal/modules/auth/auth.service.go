package auth

import (
	"context"
	"net/http"
	"strings"

	appErrors "github.com/keep/sunny/pkg/errors"
)

type Service interface {
	AuthorizationURL(state string) (string, error)
	Login(ctx context.Context, code string) (*LoginView, error)
}

type service struct {
	provider   GoogleProvider
	repository Repository
	tokens     TokenIssuer
}

func NewService(provider GoogleProvider, repository Repository, tokens TokenIssuer) Service {
	return &service{provider: provider, repository: repository, tokens: tokens}
}

func (service *service) AuthorizationURL(state string) (string, error) {
	if !service.provider.Enabled() {
		return "", appErrors.NewAPIError(http.StatusServiceUnavailable, "Google 登录尚未配置")
	}
	return service.provider.AuthorizationURL(state), nil
}

func (service *service) Login(ctx context.Context, code string) (*LoginView, error) {
	if !service.provider.Enabled() {
		return nil, appErrors.NewAPIError(http.StatusServiceUnavailable, "Google 登录尚未配置")
	}
	profile, err := service.provider.Exchange(ctx, code)
	if err != nil {
		return nil, appErrors.WrapAPIError(appErrors.NewAPIError(http.StatusBadGateway, "Google 登录失败"), err)
	}
	profile.Email = strings.ToLower(strings.TrimSpace(profile.Email))
	profile.Name = strings.TrimSpace(profile.Name)
	if profile.Email == "" || !profile.VerifiedEmail {
		return nil, appErrors.NewAPIError(http.StatusUnauthorized, "Google 邮箱未验证")
	}
	if profile.Name == "" {
		profile.Name = strings.Split(profile.Email, "@")[0]
	}
	entity, isNewUser, err := service.repository.FindOrCreate(ctx, profile.Name, profile.Email)
	if err != nil {
		return nil, appErrors.WrapAPIError(appErrors.ErrInternalServerError, err)
	}
	tokens, err := service.tokens.Issue(entity)
	if err != nil {
		return nil, appErrors.WrapAPIError(appErrors.ErrInternalServerError, err)
	}
	return &LoginView{
		AccessToken: tokens.AccessToken, TokenType: "Bearer", ExpiresAt: tokens.AccessExpiresAt, IsNewUser: isNewUser, session: tokens,
		User: UserView{ID: entity.ID, Name: entity.Name, Email: entity.Email, ReferralCode: entity.ReferralCode, CreditBalance: entity.CreditBalance},
	}, nil
}
