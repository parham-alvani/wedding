package config

import (
	"time"

	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/generator"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logger"
	"github.com/parham-alvani/wedding/wedback/internal/infra/wedding"
)

// Default returns the default configuration.
//
// The wedding section holds deliberately generic placeholders: forks that
// never write a config.toml should see obvious dummy values rather than
// somebody else's names silently baked into their invitations.
func Default() Config {
	// nolint: exhaustruct, mnd
	return Config{
		Logger: logger.Config{
			Level: "debug",
		},
		Database: db.Config{
			DSN:             "wedding.db",
			Debug:           true,
			MaxIdelConns:    10,
			MaxOpenConns:    10,
			ConnMaxIdleTime: 10 * time.Second,
			ConnMaxLifetime: 10 * time.Second,
		},
		Generator: generator.Config{
			Type: "simple",
		},
		Wedding: wedding.Config{
			HusbandName: "Partner One",
			WifeName:    "Partner Two",
			BaseURL:     "http://127.0.0.1:4321",
		},
	}
}
