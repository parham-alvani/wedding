package config

import (
	"time"

	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logger"
)

// Default return default configuration.
func Default() Config {
	// nolint: exhaustruct
	return Config{
		Logger: logger.Config{
			Level: "debug",
		},
		Database: db.Config{
			DSN:             "sqlite://wedding.db",
			Debug:           true,
			MaxIdelConns:    10,
			MaxOpenConns:    10,
			ConnMaxIdleTime: 10 * time.Second,
			ConnMaxLifetime: 10 * time.Second,
		},
	}
}
