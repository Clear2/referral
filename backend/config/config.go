package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v2"
)

const DefaultPath = "config.yaml"

type (
	// Config -.
	Config interface {
		App() App
		HTTP() HTTP
		Log() Log
		PG() PG
		Redis() Redis
		Sentry() Sentry
		Tracing() Tracing
		Metrics() Metrics
		Swagger() Swagger
		Google() Google
		JWT() JWT
		Admin() Admin
		Registration() Registration
	}

	// configImpl -.
	configImpl struct {
		AppData          App          `yaml:"app"`
		HTTPData         HTTP         `yaml:"http"`
		LogData          Log          `yaml:"log"`
		PGData           PG           `yaml:"postgres"`
		SentryData       Sentry       `yaml:"sentry"`
		TracingData      Tracing      `yaml:"tracing"`
		MetricsData      Metrics      `yaml:"metrics"`
		SwaggerData      Swagger      `yaml:"swagger"`
		GoogleData       Google       `yaml:"google"`
		GitHubData       Google       `yaml:"github"`
		RedisData        Redis        `yaml:"redis"`
		GRPCData         GRPC         `yaml:"grpc"`
		RMQData          RMQ          `yaml:"rmq"`
		JWTData          JWT          `yaml:"jwt"`
		AdminData        Admin        `yaml:"admin"`
		RegistrationData Registration `yaml:"registration"`
	}

	// App -.
	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Env     string `yaml:"env"`
	}

	// HTTP -.
	HTTP struct {
		Port              string   `yaml:"port"`
		CookieName        string   `yaml:"cookie_name"`
		RefreshCookieName string   `yaml:"refresh_cookie_name"`
		CookieSecure      bool     `yaml:"cookie_secure"`
		AllowedOrigins    []string `yaml:"allowed_origins"`
	}
	JWT struct {
		Secret     string `yaml:"secret"`
		AccessTTL  string `yaml:"access_ttl"`
		RefreshTTL string `yaml:"refresh_ttl"`
	}
	Admin struct {
		UserIDs   []int          `yaml:"user_ids"`
		Emails    []string       `yaml:"emails"`
		Bootstrap AdminBootstrap `yaml:"bootstrap"`
	}
	AdminBootstrap struct {
		Enabled  bool   `yaml:"enabled"`
		Name     string `yaml:"name"`
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
	}
	Registration struct {
		DemoEnabled bool   `yaml:"demo_enabled"`
		DemoCode    string `yaml:"demo_code"`
	}

	// Log -.
	Log struct {
		Dir   string `yaml:"dir"`
		Level string `yaml:"level"`
	}

	// PG -.
	PG struct {
		PoolMax int    `yaml:"pool_max"`
		URL     string `yaml:"url"`
	}

	// Sentry -.
	Sentry struct {
		DSN string `yaml:"dsn"`
	}

	// Metrics -.
	Tracing struct {
		Enabled bool `yaml:"enabled"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `yaml:"enabled"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `yaml:"enabled"`
	}

	Google struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		RedirectURL  string `yaml:"redirect_url"`
	}
	Redis struct {
		PoolMax  int    `yaml:"pool_max"`
		URL      string `yaml:"url"`
		Password string `yaml:"password"`
		Prefix   string `yaml:"prefix"`
	}
	GRPC struct {
		Port string `yaml:"port"`
	}
	RMQ struct {
		ServerExchange string `yaml:"server_exchange"`
		ClientExchange string `yaml:"client_exchange"`
		URL            string `yaml:"url"`
	}
)

func (a App) IsDev() bool {
	return a.Env == "development" || a.Env == "debug" || a.Env == "test"
}

// NewConfig returns app config.
func NewConfig() (Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = DefaultPath
	}
	return NewConfigFromFile(path)
}

func NewConfigFromFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}
	cfg := &configImpl{}
	if err = yaml.UnmarshalStrict(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}
	if err = cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config %q: %w", path, err)
	}
	return cfg, nil
}

func (c *configImpl) validate() error {
	required := map[string]string{
		"app.name": c.AppData.Name, "app.version": c.AppData.Version, "app.env": c.AppData.Env,
		"http.port": c.HTTPData.Port,
		"log.dir":   c.LogData.Dir, "log.level": c.LogData.Level,
		"postgres.url":     c.PGData.URL,
		"http.cookie_name": c.HTTPData.CookieName,
		"jwt.secret":       c.JWTData.Secret,
	}
	for name, value := range required {
		if value == "" {
			return fmt.Errorf("%s is required", name)
		}
	}
	if c.PGData.PoolMax < 1 {
		return fmt.Errorf("postgres.pool_max must be greater than zero")
	}
	if c.RedisData.URL == "" {
		return fmt.Errorf("redis.url is required")
	}
	if c.RedisData.PoolMax < 1 {
		return fmt.Errorf("redis.pool_max must be greater than zero")
	}
	return nil
}

// Implementation of ConfigInter methods
func (c *configImpl) App() App                   { return c.AppData }
func (c *configImpl) HTTP() HTTP                 { return c.HTTPData }
func (c *configImpl) Log() Log                   { return c.LogData }
func (c *configImpl) PG() PG                     { return c.PGData }
func (c *configImpl) Redis() Redis               { return c.RedisData }
func (c *configImpl) Sentry() Sentry             { return c.SentryData }
func (c *configImpl) Tracing() Tracing           { return c.TracingData }
func (c *configImpl) Metrics() Metrics           { return c.MetricsData }
func (c *configImpl) Swagger() Swagger           { return c.SwaggerData }
func (c *configImpl) Google() Google             { return c.GoogleData }
func (c *configImpl) JWT() JWT                   { return c.JWTData }
func (c *configImpl) Admin() Admin               { return c.AdminData }
func (c *configImpl) Registration() Registration { return c.RegistrationData }
