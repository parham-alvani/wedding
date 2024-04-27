package generator

type Generator interface {
	ID() string
}

func Provide(cfg Config) Generator {
	// nolint: gocritic
	switch cfg.Type {
	case "simple":
		return new(Simple)
	}

	return new(Simple)
}
