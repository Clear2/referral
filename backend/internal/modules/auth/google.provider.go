package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/keep/sunny/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

type GoogleProvider interface {
	Enabled() bool
	AuthorizationURL(state string) string
	Exchange(ctx context.Context, code string) (GoogleProfile, error)
}

type googleProvider struct {
	config  oauth2.Config
	enabled bool
}

func NewGoogleProvider(cfg config.Config) GoogleProvider {
	googleConfig := cfg.Google()
	return &googleProvider{
		enabled: googleConfig.ClientID != "" && googleConfig.ClientSecret != "" && googleConfig.RedirectURL != "",
		config: oauth2.Config{
			ClientID: googleConfig.ClientID, ClientSecret: googleConfig.ClientSecret, RedirectURL: googleConfig.RedirectURL,
			Endpoint: google.Endpoint, Scopes: []string{"openid", "email", "profile"},
		},
	}
}

func (provider *googleProvider) Enabled() bool { return provider.enabled }

func (provider *googleProvider) AuthorizationURL(state string) string {
	return provider.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (provider *googleProvider) Exchange(ctx context.Context, code string) (GoogleProfile, error) {
	token, err := provider.config.Exchange(ctx, code)
	if err != nil {
		return GoogleProfile{}, fmt.Errorf("exchange Google authorization code: %w", err)
	}
	response, err := provider.config.Client(ctx, token).Get(googleUserInfoURL)
	if err != nil {
		return GoogleProfile{}, fmt.Errorf("get Google user info: %w", err)
	}
	defer func() {
		// The response body is fully consumed below; a close error cannot affect the result.
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		return GoogleProfile{}, fmt.Errorf("get Google user info: unexpected status %d", response.StatusCode)
	}
	var profile GoogleProfile
	if err := json.NewDecoder(response.Body).Decode(&profile); err != nil {
		return GoogleProfile{}, fmt.Errorf("decode Google user info: %w", err)
	}
	return profile, nil
}
