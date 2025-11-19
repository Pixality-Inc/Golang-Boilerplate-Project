package migrations

import migrate "github.com/rubenv/sql-migrate"

var Migrations = &migrate.MemoryMigrationSource{
	Migrations: []*migrate.Migration{
		{
			Id: "2025_11_19_09_00_00_Init",
			Up: []string{
				// Extensions
				`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,
				`CREATE EXTENSION IF NOT EXISTS "postgis"`,
				`CREATE EXTENSION IF NOT EXISTS "pg_stat_statements"`,
				// Search path
				`SET search_path TO public`,
				// Books
				`CREATE TABLE IF NOT EXISTS "books" (
	"id"         UUID         NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
	"title"      VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMPTZ  NOT NULL DEFAULT NOW()
)`,
				`CREATE UNIQUE INDEX IF NOT EXISTS "books_id_idx" ON "books" ("id")`,
			},
			Down: []string{
				`DROP TABLE IF EXISTS "books"`,
			},
		},
	},
}
