package main

import (
	"errors"
	"flag"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"

	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/cmd"
)

func main() {
	flag.Parse()

	config := cmd.GetConfig(cmd.ConfigFile)

	filesPath := "file://migrations"
	dbURL := config.Database.String()
	migration, err := migrate.New(filesPath, dbURL)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %s", err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed to apply migrations: %s", err)
	}

	log.Println("Migrations applied successfully!")
}
