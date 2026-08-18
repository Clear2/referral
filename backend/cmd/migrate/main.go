package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/keep/sunny/config"
)

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "up" && os.Args[1] != "down") {
		log.Fatal("usage: migrate <up|down>")
	}
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	runner, err := migrate.New("file://ent/migrate/migrations", cfg.PG().URL)
	if err != nil {
		log.Fatalf("create migrator: %v", err)
	}
	defer func() {
		sourceErr, databaseErr := runner.Close()
		if sourceErr != nil || databaseErr != nil {
			log.Printf("close migrator: source=%v database=%v", sourceErr, databaseErr)
		}
	}()
	if os.Args[1] == "up" {
		err = runner.Up()
	} else {
		err = runner.Steps(-1)
	}
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(fmt.Errorf("migrate %s: %w", os.Args[1], err))
	}
}
