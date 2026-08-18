//go:build ignore

package main

import (
	"context"
	"database/sql"
	"log"
	"net/url"
	"os"

	"ariga.io/atlas/sql/sqltool"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	pgx "github.com/jackc/pgx/v5/stdlib"
	"github.com/keep/sunny/config"
	"github.com/keep/sunny/ent/migrate"
)

func init() {
	pgxDriver := pgx.GetDefaultDriver()
	sql.Register("postgres", pgxDriver)
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("load migration config: %v", err)
	}
	databaseURL := cfg.PG().URL

	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		log.Fatalf("invalid postgres.url: %v", err)
	}
	query := parsedURL.Query()
	if query.Get("sslmode") == "" {
		query.Set("sslmode", "disable")
		parsedURL.RawQuery = query.Encode()
	}
	databaseURL = parsedURL.String()

	ctx := context.Background()
	// Create a local migration directory able to understand golang-migrate migration file format for replay.
	dir, err := sqltool.NewGolangMigrateDir("ent/migrate/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}
	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
	}
	if len(os.Args) != 2 {
		log.Fatalln("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
	}
	// Generate migrations using Atlas support for PostgreSQL (note the Ent dialect option passed above).
	err = migrate.NamedDiff(ctx, databaseURL, os.Args[1], opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}
