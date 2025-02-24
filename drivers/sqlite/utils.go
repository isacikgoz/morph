package sqlite

import "github.com/isacikgoz/morph/drivers"

func getDefaultConfig() *Config {
	return &Config{
		Config: drivers.Config{
			MigrationsTable:        "db_migrations",
			StatementTimeoutInSecs: 60,
			MigrationMaxSize:       defaultMigrationMaxSize,
		},
	}
}
