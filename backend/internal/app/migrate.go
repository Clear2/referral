//go:build migrate

package app

import (
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/keep/sunny/config"
)

const (
	_defaultAttempts = 10
	_defaultTimeout  = time.Second
)

func init() {
	cfg, configErr := config.NewConfig()
	if configErr != nil {
		log.Fatalf("Migrate: load config: %s", configErr)
	}
	databaseURL := cfg.PG().URL

	parsedURL, parseErr := url.Parse(databaseURL)
	if parseErr != nil {
		log.Fatalf("Migrate: invalid postgres.url: %s", parseErr)
	}
	query := parsedURL.Query()
	if query.Get("sslmode") == "" {
		query.Set("sslmode", "disable")
		parsedURL.RawQuery = query.Encode()
	}
	databaseURL = parsedURL.String()

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://ent/migrate/migrations", databaseURL)
		if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}
