package generator

import "github.com/parham-alvani/wedding/wedback/internal/domain/generator"

func Provide(cfg Config) generator.Generator {
	// nolint: gocritic
	switch cfg.Type {
	case "simple":
		return new(Simple)
	}

	return new(Simple)
}
