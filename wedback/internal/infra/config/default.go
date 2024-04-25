package config

import "github.com/parham-alvani/wedding/wedback/internal/infra/logger"

// Default return default configuration.
func Default() Config {
	// nolint: exhaustruct
	return Config{
		Logger: logger.Config{
			Level: "debug",
		},
	}
}
