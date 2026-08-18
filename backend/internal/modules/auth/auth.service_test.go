package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/keep/sunny/ent"
	appErrors "github.com/keep/sunny/pkg/errors"
)

type fakeGoogleProvider struct {
	profile GoogleProfile
	err     error
	enabled bool
}

func (provider *fakeGoogleProvider) Enabled() bool { return provider.enabled }
func (provider *fakeGoogleProvider) AuthorizationURL(state string) string {
	return "https://accounts.example/authorize?state=" + state
}

func (provider *fakeGoogleProvider) Exchange(context.Context, string) (GoogleProfile, error) {
	return provider.profile, provider.err
}

type fakeRepository struct {
	user  *ent.User
	isNew bool
	name  string
	email string
	err   error
}

func (repository *fakeRepository) FindOrCreate(_ context.Context, name, email string) (*ent.User, bool, error) {
	repository.name, repository.email = name, email
	return repository.user, repository.isNew, repository.err
}

type fakeTokenIssuer struct {
	tokens SessionTokens
	err    error
}

func (issuer *fakeTokenIssuer) Issue(*ent.User) (SessionTokens, error) {
	return issuer.tokens, issuer.err
}

func TestGoogleLoginCreatesUserAndIssuesToken(t *testing.T) {
	t.Parallel()
	expiresAt := time.Date(2026, time.August, 20, 0, 0, 0, 0, time.UTC)
	provider := &fakeGoogleProvider{enabled: true, profile: GoogleProfile{Name: " Alice ", Email: " ALICE@Example.com ", VerifiedEmail: true}}
	repository := &fakeRepository{user: &ent.User{ID: 1, Name: "Alice", Email: "alice@example.com"}, isNew: true}

	result, err := NewService(provider, repository, &fakeTokenIssuer{tokens: SessionTokens{AccessToken: "token", AccessExpiresAt: expiresAt, RefreshToken: "refresh", RefreshExpiresAt: expiresAt.Add(time.Hour)}}).Login(context.Background(), "code")
	if err != nil {
		t.Fatalf("Login(%q) error = %v, want nil", "code", err)
	}
	if repository.name != "Alice" || repository.email != "alice@example.com" {
		t.Errorf("Login(%q) normalized identity = (%q, %q), want (%q, %q)", "code", repository.name, repository.email, "Alice", "alice@example.com")
	}
	if result.AccessToken != "token" || !result.IsNewUser || !result.ExpiresAt.Equal(expiresAt) {
		t.Errorf("Login(%q) result = %+v, want token, new user and expiry %v", "code", result, expiresAt)
	}
}

func TestGoogleLoginRejectsUnverifiedEmail(t *testing.T) {
	t.Parallel()
	provider := &fakeGoogleProvider{enabled: true, profile: GoogleProfile{Email: "alice@example.com", VerifiedEmail: false}}

	_, err := NewService(provider, &fakeRepository{}, &fakeTokenIssuer{}).Login(context.Background(), "code")
	var apiErr *appErrors.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Login(%q) error = %v, want *errors.APIError", "code", err)
	}
	if apiErr.Code != 401 {
		t.Errorf("Login(%q) status = %d, want 401", "code", apiErr.Code)
	}
}

func TestAuthorizationURLRejectsMissingConfiguration(t *testing.T) {
	t.Parallel()
	_, err := NewService(&fakeGoogleProvider{}, &fakeRepository{}, &fakeTokenIssuer{}).AuthorizationURL("state")
	var apiErr *appErrors.APIError
	if !errors.As(err, &apiErr) || apiErr.Code != 503 {
		t.Errorf("AuthorizationURL(%q) error = %v, want status 503", "state", err)
	}
}
