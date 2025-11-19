package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/pixality-inc/golang-boilerplate-project/internal/migrations"
	"github.com/pixality-inc/golang-boilerplate-project/internal/wiring"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	wire := wiring.New()
	defer wire.Shutdown()

	log := wire.Log

	if len(os.Args) < 2 {
		log.Fatal("No direction is provided.\nUsage: migrate <up|down>") //nolint:gocritic
	}

	migrateDirectionArg := os.Args[1]

	var migrateDirection migrate.MigrationDirection

	switch migrateDirectionArg {
	case "up":
		migrateDirection = migrate.Up
	case "down":
		migrateDirection = migrate.Down
	default:
		log.Fatalf("Unknown direction '%s'.\nUsage: migrate <up|down>", migrateDirectionArg)
	}

	dbConfig := wire.Config.Database

	migrationsUrl := dbConfig.ParamsUrl()

	log.
		WithFields(map[string]any{
			"host":     dbConfig.HostValue,
			"port":     dbConfig.PortValue,
			"user":     dbConfig.UserValue,
			"database": dbConfig.DatabaseValue,
			"schema":   dbConfig.SchemaValue,
		}).
		Infof("Applying migrations")

	db, err := sql.Open("postgres", migrationsUrl)
	if err != nil {
		log.WithError(err).Fatal("Error connecting to database")
	}

	n, err := migrate.Exec(db, "postgres", migrations.Migrations, migrateDirection)
	if err != nil {
		log.WithError(err).Fatal("Error running migrations")
	}

	log.Infof("Applied %d migrations", n)
}
