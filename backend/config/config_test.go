package config

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validYAML = `
app: {name: referral, version: 1.0.0, env: test}
http: {port: "8999", cookie_name: session, refresh_cookie_name: refresh}
jwt: {secret: "test-secret-at-least-32-characters", access_ttl: 15m, refresh_ttl: 168h}
admin: {bootstrap: {enabled: false}}
log: {dir: /tmp/logs, level: debug}
postgres: {pool_max: 2, url: "postgres://localhost/referral"}
redis: {pool_max: 2, url: "redis://localhost:6379/0", prefix: referral}
sentry: {dsn: "https://public@example.com/1"}
tracing: {enabled: false}
metrics: {enabled: true}
swagger: {enabled: true}
`

func TestNewConfigFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(validYAML), 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", path, err)
	}

	cfg, err := NewConfigFromFile(path)
	if err != nil {
		t.Fatalf("NewConfigFromFile(%q) error = %v", path, err)
	}
	if got, want := cfg.HTTP().Port, "8999"; got != want {
		t.Errorf("NewConfigFromFile(%q).HTTP().Port = %q, want %q", path, got, want)
	}
	if got, want := cfg.Redis().Prefix, "referral"; got != want {
		t.Errorf("NewConfigFromFile(%q).Redis().Prefix = %q, want %q", path, got, want)
	}
}

func TestReferralConfigFilesStayAligned(t *testing.T) {
	t.Parallel()
	paths := []string{
		filepath.Join("..", "config.yaml"),
		filepath.Join("..", "config.example.yaml"),
		filepath.Join("..", "..", "deploy", "config.production.yaml"),
		filepath.Join("..", "..", "deploy", "config.production.example.yaml"),
	}
	for _, path := range paths {
		path := path
		t.Run(path, func(t *testing.T) {
			cfg, err := NewConfigFromFile(path)
			if err != nil {
				t.Fatalf("NewConfigFromFile(%q): %v", path, err)
			}
			postgresURL, err := url.Parse(cfg.PG().URL)
			if err != nil {
				t.Fatalf("parse postgres URL in %q: %v", path, err)
			}
			if cfg.App().Name != "referral" || postgresURL.Path != "/referral" || cfg.Redis().Prefix != "referral" {
				t.Fatalf("%q is not Referral-aligned", path)
			}
			if cfg.HTTP().CookieName != "referral_session" || cfg.HTTP().RefreshCookieName != "referral_refresh" {
				t.Fatalf("%q uses non-Referral cookie names", path)
			}
		})
	}
}

func TestNewConfigFromFileRejectsUnknownField(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	data := validYAML + "unknown_section: true\n"
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", path, err)
	}

	_, err := NewConfigFromFile(path)
	if err == nil || !strings.Contains(err.Error(), "unknown_section") {
		t.Errorf("NewConfigFromFile(%q) error = %v, want unknown-field error", path, err)
	}
}
